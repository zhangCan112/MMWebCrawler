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

// Spider 爬虫解析模块的接口
type Spider interface {
	extractHTML(rw *ResultWriter, doc *goquery.Document)
}

// SpiderFunc 就是一个允许普通函数做为Spider的适配器，
type SpiderFunc func(rw *ResultWriter, doc *goquery.Document)

// extractHTML Spider接口的实现
func (sf SpiderFunc) extractHTML(rw *ResultWriter, doc *goquery.Document) {
	sf(rw, doc)
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
	mux sync.RWMutex
	m   map[string]muxEntry
	es  []muxEntry
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
	sl.mux.Lock()
	defer sl.mux.Unlock()

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
	sl.mux.RLock()
	defer sl.mux.RUnlock()

	s, _ := sl.match(url)
	return s
}

// Clean 根据指定pattern清除Spider
func (sl *Spiderlair) Clean(pattern string) {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	delete(sl.m, pattern)
	sl.cleanEs(pattern)
}

// CleanAll 清理所有的spider
func (sl *Spiderlair) CleanAll() {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	sl.m = map[string]muxEntry{}
	sl.es = nil
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
