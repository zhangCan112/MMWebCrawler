package downloader

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var (
	userAgentList = [...]string{
		// "User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		// "User-Agent, Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv,2.0.1) Gecko/20100101 Firefox/4.0.1",
		// "User-Agent,Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"User-Agent, Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		// "User-Agent,Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4094.1 Safari/537.36",
	}
	defaultDownloader = NewDownloader()
)

// Downloader 下载器
type Downloader struct{}

// Download 使用默认下载器请求指定url上的数据
func Download(url string) (doc *goquery.Document, err error) {
	return defaultDownloader.Download(url)
}

// NewDownloader 返回一个新的Downloader实例
func NewDownloader() *Downloader {
	return &Downloader{}
}

// Download 对指定url发起请求
func (dl *Downloader) Download(rawurl string) (doc *goquery.Document, err error) {
	res, err := get(rawurl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, _ = goquery.NewDocumentFromReader(res.Body)
	url, _ := url.Parse(rawurl)
	doc.Url = url
	return doc, nil
}

// get  get请求的简单封装
func get(url string) (resp *http.Response, err error) {
	req, err := newRequest("GET", url)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// NewRequest 返回一个新的Request请求
func newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err == nil {
		req.Header.Set("User-Agent", randomUserAgent())
	}

	return req, err
}

func randomUserAgent() string {
	idx := rand.Intn(len(userAgentList))
	return userAgentList[idx]
}
