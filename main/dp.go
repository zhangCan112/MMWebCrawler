package main

import (
	"github.com/PuerkitoBio/goquery"
	webcrawler "github.com/zhangCan112/webcrawler/app"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/spider"
)

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
