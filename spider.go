package webcrawler

import (
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// ResultWriter 爬虫结果写入器接口
type ResultWriter interface {
	Items() []interface{}
	URLs() []string
	AddItem(item interface{})
	AddURL(url string)
}

// NewResultWriter 返回一个ResultWriter的默认实现实例
func NewResultWriter() ResultWriter {
	return &defaultResultWriter{
		items: make([]interface{}, 0),
		urls:  make([]string, 0),
	}
}

type defaultResultWriter struct {
	itemsMu sync.RWMutex
	items   []interface{}
	urlsMu  sync.RWMutex
	urls    []string
}

func (rw *defaultResultWriter) Items() []interface{} {
	rw.itemsMu.RLock()
	defer rw.itemsMu.RUnlock()

	cp := make([]interface{}, len(rw.items))
	copy(cp, rw.items)
	return cp
}

func (rw *defaultResultWriter) URLs() []string {
	rw.urlsMu.RLock()
	defer rw.urlsMu.RUnlock()

	cp := make([]string, len(rw.urls))
	copy(cp, rw.urls)
	return cp
}

func (rw *defaultResultWriter) AddItem(item interface{}) {
	rw.itemsMu.Lock()
	defer rw.itemsMu.Unlock()

	rw.items = append(rw.items, item)
}

func (rw *defaultResultWriter) AddURL(url string) {
	rw.urlsMu.Lock()
	defer rw.urlsMu.Unlock()

	rw.urls = append(rw.urls, url)
}

// Spider 爬虫解析模块的接口
type Spider interface {
	extractHTML(rw ResultWriter, doc *goquery.Document)
}

// SpiderFunc 就是一个允许普通函数做为Spider的适配器，
type SpiderFunc func(rw ResultWriter, doc *goquery.Document)

// extractHTML Spider接口的实现
func (sf SpiderFunc) extractHTML(rw ResultWriter, doc *goquery.Document) {
	sf(rw, doc)
}

// SpiderServer 用于实现Spider的监听功能
type SpiderServer interface {
	Do(url string, doc *goquery.Document)
	Handle(pattern string, sp Spider)
	HandleFunc(pattern string, spfunc func(rw ResultWriter, doc *goquery.Document))
	ItemsChan() <-chan []interface{}
	URLsChan() <-chan []string
}

// NewSpiderServer 创建一个默认的SpiderServer实例
func NewSpiderServer() SpiderServer {
	return &defaultSpiderServer{
		spl:       NewSpiderlair(),
		itemsChan: make(chan []interface{}),
	}
}

// defaultSpiderServer SpiderServer的默认内部实现
type defaultSpiderServer struct {
	spl       *Spiderlair
	itemsChan chan []interface{}
	urlsChan  chan []string
}

func (ss *defaultSpiderServer) Do(url string, doc *goquery.Document) {
	go ss.do(url, doc)
}

func (ss *defaultSpiderServer) do(url string, doc *goquery.Document) {
	sp := ss.spl.Spider(url)
	rw := NewResultWriter()
	sp.extractHTML(rw, doc)
	ss.itemsChan <- rw.Items()
	ss.urlsChan <- rw.URLs()
}

// Handle 指定路径注册一个爬虫
func (ss *defaultSpiderServer) Handle(pattern string, sp Spider) {
	ss.spl.Join(pattern, sp)
}

// HandleFunc 指定路径注册一个爬虫函数
func (ss *defaultSpiderServer) HandleFunc(pattern string, spfunc func(rw ResultWriter, doc *goquery.Document)) {
	ss.spl.Join(pattern, SpiderFunc(spfunc))
}

// ItemsChan 获取接收Items数据的通道
func (ss *defaultSpiderServer) ItemsChan() <-chan []interface{} {
	return ss.itemsChan
}

// URLsChan 获取接收Urls数据的通道
func (ss *defaultSpiderServer) URLsChan() <-chan []string {
	return ss.urlsChan
}

// DefaultSpiderlair 默认的Spiderlair实例
var DefaultSpiderlair = &defaultSpiderlair
var defaultSpiderlair Spiderlair

// Spiderlair 负责所有Spider统一注册管理和调度
// 它根据注册模式列表匹配每个传入的URL路径，并调用与该URL最匹配的Spider进行处理。
// 模式名固定，根路径，如“/favicon.ico”，或根子树，如“/images/”（注意尾随斜杠）。
// 较长的模式优先于较短的模式，因此，如果同时为“/images/”和“/images/thumbnails/”注册了Spider，
// 则会对以“/images/thumbnails/”开头的路径调用后者的Spider，
// 前者将接收对“/images/”子树中任何其他路径的请求。
// 此外如果注册了相同路径的多个Spider, 他们不会彼此覆盖，而是会依次被调用
type Spiderlair struct {
	mu sync.RWMutex
	m  map[string]muxEntry
	es []muxEntry
}

type muxEntry struct {
	s       Spider
	pattern string
}

type muxEntrySlice []muxEntry

// Len is the number of elements in the collection.
func (es muxEntrySlice) Len() int {
	return len(es)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (es muxEntrySlice) Less(i, j int) bool {
	return len(es[i].pattern) < len(es[j].pattern)
}

// Swap swaps the elements with indexes i and j.
func (es muxEntrySlice) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

// NewSpiderlair 返回一个新的Spiderlair实例
func NewSpiderlair() *Spiderlair {
	return new(Spiderlair)
}

// Join 指定pattern注册一个Spider到Spiderlair
func (sl *Spiderlair) Join(pattern string, sp Spider) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	me := muxEntry{
		s:       sp,
		pattern: pattern,
	}

	if sl.m == nil {
		sl.m = make(map[string]muxEntry)
	}

	sl.m[pattern] = me
	//first clean last append
	sl.cleanEs(pattern)
	sl.appendSorted(me)
}

// Spider 根据指定path查找合适的Spider
func (sl *Spiderlair) Spider(url string) Spider {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	s, _ := sl.match(url)
	return s
}

// Clean 根据指定pattern清除Spider
func (sl *Spiderlair) Clean(pattern string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	delete(sl.m, pattern)
	sl.cleanEs(pattern)
}

// CleanAll 清理所有的spider
func (sl *Spiderlair) CleanAll() {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.m = map[string]muxEntry{}
	sl.es = nil
}

// extractHTML Spider接口的实现
func (sl *Spiderlair) extractHTML(rw ResultWriter, doc *goquery.Document) {
	sp := sl.Spider(doc.Url.String())
	if sp != nil {
		sp.extractHTML(rw, doc)
	}
}

// appendSorted 排序插入muxEntry，依照pattern长度从小到大
func (sl *Spiderlair) appendSorted(me muxEntry) {
	if sl.es == nil {
		sl.es = make([]muxEntry, 100)[0:0]
	}
	sl.es = append(sl.es, me)
	sort.Sort(muxEntrySlice(sl.es))
}

// cleanEs 清除es中指定pattern的值
func (sl *Spiderlair) cleanEs(pattern string) {
	for idx, val := range sl.es {
		if val.pattern == pattern {
			sl.es = append(sl.es[:idx], sl.es[idx+1:]...)
			break
		}
	}
}

// match 在给定路径字符串的Spider映射上查找处理Spider
// 最具体（最长）匹配优先
func (sl *Spiderlair) match(url string) (s Spider, pattern string) {
	// Check for exact match first.
	v, ok := sl.m[url]
	if ok {
		return v.s, v.pattern
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest.
	for _, e := range sl.es {
		if strings.HasPrefix(url, e.pattern) {
			return e.s, e.pattern
		}
	}
	return nil, ""
}
