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
	doc.Find("#J_boxList ul li[data-shopid]").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Find("a[href]").Attr("href")
		if ok {
			rw.AddURL("http:" + url)
		}
	})

	doc.Find("#J_boxList .Pages .PageLink").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			rw.AddURL("http://www.dianping.com" + url)
		}
	})
})

//  detailSpider 点评月子中心详情页Spider
var detailSpider = spider.SpiderFunc(func(rw spider.ResultWriter, doc *goquery.Document) {
	store := webcrawler.WrapedString(doc.Find("#J_boxDetail .shop-info .shop-name .shop-title").Text()).TrimSpace().FilterLineBreaks().Unwrap()
	address, _ := doc.Find("#J_boxDetail .shop-info .shop-addr span[title]").Attr("title")
	price := webcrawler.WrapedString(doc.Find("#J_boxDetail div div[class=comment-rst] div span strong").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("￥").Unwrap()
	commentTotal := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(1) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	comment5Star := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(2) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	comment4Star := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(3) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	comment3Star := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(4) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	comment2Star := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(5) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	comment1Star := webcrawler.WrapedString(doc.Find("#J_boxReview div.comment-mode.shop-comment div.J_wrapFilter div.comment-star dl dd:nth-child(6) em").Text()).TrimSpace().FilterLineBreaks().TrimPrefix("(").TrimSuffix(")").Unwrap()
	it := pipeline.NewItem(
		"月子中心",
		[]string{
			"store",         //店名
			"address",       //地址
			"price",         //价格
			"comment_total", // 评价数
			"5Star",         //5星评价数
			"4Star",         //4星评价数
			"3Star",         //3星评价数
			"2Star",         //2星评价数
			"1Star"},        //1星评价数
		map[string]interface{}{
			"store":         store,
			"address":       address,
			"price":         price,
			"comment_total": commentTotal,
			"5Star":         comment5Star,
			"4Star":         comment4Star,
			"3Star":         comment3Star,
			"2Star":         comment2Star,
			"1Star":         comment1Star,
		},
		[]string{"csv"},
	)

	rw.AddItem(it)
})
