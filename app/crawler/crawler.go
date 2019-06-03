package crawler

import (
	"fmt"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/scheduler"
	"github.com/zhangCan112/webcrawler/app/downloader"
	"github.com/zhangCan112/webcrawler/app/spider"
)

type (
	Crawler interface {
		// Init 初始化采集器
		Init(spider spider.Spider)Crawler
		// Start 用种子URL启动采集器，至少一个
		Start(seed string, rest ...string)
		// Stop 停止采集器
		Stop()
		// HasEnd 是否以停止
		HasEnd() bool
	}

	crawler struct {
		sp spider.Spider
		dl *downloader.Downloader
		sc *scheduler.Scheduler
		pl pipeline.Pipeline
	}
)

func (c *crawler) Init(spider spider.Spider) Crawler {
	c.sp = spider
	return c
}

func (c *crawler) Start(seed string, rest ...string) {
	c.sc.Push(seed)
	for _, surl := range rest {
		c.sc.Push(surl)		
	}
}

func (c *crawler) work()  {
	url, ok := c.sc.Pop()
	if !ok {return}	
	doc, err := c.dl.Download(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.sc.Done(url)
	rw := spider.NewResultWriter()
	c.sp.ExtractHTML(rw, doc)

	urls := rw.URLs()
	for _, surl := range urls {
		c.sc.Push(surl)		
	}	

	its := rw.Items()
	c.pl.Write(its[0], its[1:]...)
}

func (c *crawler) Stop() {
	
}

func (c *crawler) HasEnd() bool {
	return false
}

