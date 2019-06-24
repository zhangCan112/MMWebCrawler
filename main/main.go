package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	webcrawler "github.com/zhangCan112/webcrawler/app"
	"github.com/zhangCan112/webcrawler/app/crawler"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/spider"
)

func main() {
	dpDemo()
}

func dpDemo() {
	cr := crawler.NewCrawler()
	cr.Init(DBSpider, pipeline.CSVWriter, 3)
	cr.Start("http://t.dianping.com/list/xian?q=%E4%B8%93%E4%B8%9A%E8%84%B1%E6%AF%9B")
}

func exampleScrape() {
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

// DBSpider 点评解析spider
var DBSpider = spider.SpiderFunc(func(rw spider.ResultWriter, doc *goquery.Document) {
	doc.Find(".tg-floor-item").Each(func(i int, s *goquery.Selection) {
		title := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-title h3").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		subTitle := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-title h4").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		price := webcrawler.WrapedString(s.Find(".tg-floor-item-wrap .tg-floor-price-new em").Text()).TrimSpace().FilterLineBreaks().Unwrap()
		it := pipeline.NewItem(
			"DianPing",
			[]string{"title", "subTitle", "price"},
			map[string]interface{}{"title": title, "subTitle": subTitle, "price": price},
			[]string{"csv"},
		)
		rw.AddItem(it)
	})
})
