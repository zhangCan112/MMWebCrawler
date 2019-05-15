package webcrawler

import (
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

// Spiderlair 负责所有Spider统一管理和调度
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

// match 在给定路径字符串的Spider映射上查找处理Spider
// 最具体（最长）匹配优先
func (sl *Spiderlair) match(path string) (h Spider, pattern string) {
	return nil, ""
}
