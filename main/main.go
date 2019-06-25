package main

import (
	"github.com/zhangCan112/webcrawler/app/crawler"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/spider"
)

func main() {
	//http://www.dianping.com/shop/77005354
	//http://www.dianping.com/xian/ch70/g2784p2?aid=01b2e5cd09781cb45d747fbce1d216d2
	sl := spider.NewSpiderlair()
	sl.Join("http://www.dianping.com/shop/", detailSpider)
	sl.Join("http://www.dianping.com/xian/ch70/", listSpider)
	cr := crawler.NewCrawler()
	cr.Init(sl, pipeline.CSVWriter, 3)
	cr.Start("http://www.dianping.com/xian/ch70/g2784")
}
