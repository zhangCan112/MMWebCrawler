package webcrawler

import "net/http"

// Item 解析好的数据单元接口
type Item interface {
}

// Request 请求数据
type Request struct {
	http.Request
}

// Response 请求响应数据
type Response struct {
	http.Response
}

// Spider 爬虫解析模块的接口
type Spider interface {
	Parse(resp *Response) (items []*Item, reqs []*Request, err error)
}
