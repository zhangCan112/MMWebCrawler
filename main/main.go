package main

import (
	"github.com/zhangCan112/webcrawler/app/crawler"
	"github.com/zhangCan112/webcrawler/app/pipeline"
)

func main() {
	cr := crawler.NewCrawler()
	cr.Init(listSpider, pipeline.CSVWriter, 3)
	cr.Start("http://www.dianping.com/xian/ch70/g2784")
}
