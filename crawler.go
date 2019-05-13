package webcrawler

import "net/http"

// Item 解析好的数据单元接口
type Item interface {
}

// Spider 爬虫解析模块的接口
type Spider interface {
	Parse(resp *http.Response) (items []*Item, urls []string, err error)
}
