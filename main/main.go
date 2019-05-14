package main

import (
	"fmt"
	"log"
	"net/http"
	"webcrawler"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	test()
}

func test() {
	docReceiver := webcrawler.RunDownloader()
	webcrawler.Download("http://www.baidu.com")
	resultFunc := <-docReceiver
	doc, url, err := resultFunc()
	fmt.Println(doc)
	fmt.Println(url)
	fmt.Println(err)
	fmt.Println(doc.Text())
}

func ExampleScrape() {
	// Request the HTML page.
	res, err := Get("http://t.dianping.com/list/xian?q=%E4%B8%93%E4%B8%9A%E8%84%B1%E6%AF%9B")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".tg-floor-item").Each(func(i int, s *goquery.Selection) {
		title := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-title h3").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		subTitle := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-title h4").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		price := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-price-new em").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		// saleCount, _ := s.Attr("data-static-join")

		fmt.Printf("%s %s 价格：%s \n", title, subTitle, price)
	})
	// Find the review items
	// doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	band := s.Find("a").Text()
	// 	title := s.Find("i").Text()
	// 	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// })
}

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
