package webcrawler

import (
	"math/rand"
	"net/http"
)

var (
	userAgentList = [...]string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1",
		"Mozilla/5.0 (X11; CrOS i686 2268.111.0) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1092.0 Safari/536.6",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1090.0 Safari/536.6",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/19.77.34.5 Safari/537.1",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.9 Safari/536.5",
		"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.36 Safari/536.5",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_0) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.0 Safari/536.3",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
	}
	DefaultDownloader = NewDownloader()
)

// Downloader 下载器
type Downloader struct {
	Sender chan (*http.Response, error)
}

// RunDownloader启动默认的Downloader模块
func RunDownloader(receiver <-chan string, sender chan<- *http.Response) {
	DefaultDownloader.Run(receiver, sender)
}

func NewDownloader() *Downloader {
	return &Downloader{
		Sender: make(chan)
	}
}

// Run 下载器开启工作模式
func (dl *Downloader) Run(receiver <-chan string, sender chan<- *http.Response) {
	dl.Receiver = receiver
	dl.Sender = sender
	go dl.work()
}

func (dl *Downloader) work() {
	for {
		select {
		case <-dl.Receiver:
		case dl.Sender <- nil:
		}
	}
}

// Get get请求的简单封装
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

// download 开启请求指定url的数据
func download(url string) (res *http.Response, err error) {
	return get(url)
}
