package main

import (
	"github.com/PuerkitoBio/goquery"
	webcrawler "github.com/zhangCan112/webcrawler/app"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/spider"
)

// DBSpider 点评解析spider
var dpSpider = spider.SpiderFunc(func(rw spider.ResultWriter, doc *goquery.Document) {
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
		rw.AddURL("http://t.dianping.com/list/xian?q=月子中心")
	})
})

//  listSpider 点评月子中心列表页Spider
var listSpider = spider.SpiderFunc(func(rw spider.ResultWriter, doc *goquery.Document) {
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
		rw.AddURL("http://t.dianping.com/list/xian?q=月子中心")
	})
})

//  detailSpider 点评月子中心详情页Spider
var detailSpider = spider.SpiderFunc(func(rw spider.ResultWriter, doc *goquery.Document) {
	store := webcrawler.WrapedString(doc.Find("#J_boxDetail .shop-info .shop-name .shop-title").Text()).TrimSpace().FilterLineBreaks().Unwrap()
	address, _ := doc.Find("#J_boxDetail .shop-info .shop-addr span[title]").Attr("title")
	price := webcrawler.WrapedString(doc.Find("#J_boxDetail div div[class=comment-rst] div span strong").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("￥").Unwrap()
	commentTotal := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(1) em").Text()).TrimSpace().FilterLineBreaks().Unwrap()

	it := pipeline.NewItem(
		"月子中心",
		[]string{"store", "address", "price", "comment_total"},
		map[string]interface{}{"store": store, "address": address, "price": price, "comment_total": commentTotal},
		[]string{"csv"},
	)

	rw.AddItem(it)
})
