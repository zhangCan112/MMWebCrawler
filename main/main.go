package main

import (
	"github.com/zhangCan112/webcrawler/app/crawler"
	"github.com/zhangCan112/webcrawler/app/pipeline"
)

func main() {
	cr := crawler.NewCrawler()
	cr.Init(detailSpider, pipeline.CSVWriter, 3)
	cr.Start("http://www.dianping.com/shop/69043117")
}
