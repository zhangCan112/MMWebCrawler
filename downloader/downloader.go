package webcrawler

import "net/http"

// Get get请求的简单封装
func Get(url string) (resp *http.Response, err error) {
	req, err := NewRequest("GET", url)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// NewRequest 返回一个新的Request请求
func NewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err == nil {
		// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		// req.Header.Set("Accept-Encoding", "gzip, deflate")
		// req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		// req.Header.Set("Cache-Control", "max-age=0")
		// req.Header.Set("Connection", "keep-alive")
		// req.Header.Set("Host", "t.dianping.com")
		// req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.3 Safari/605.1.15")
	}

	return req, err
}
