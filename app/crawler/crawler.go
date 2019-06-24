package crawler

import (
	"fmt"

	"github.com/zhangCan112/webcrawler/app/downloader"
	"github.com/zhangCan112/webcrawler/app/pipeline"
	"github.com/zhangCan112/webcrawler/app/scheduler"
	"github.com/zhangCan112/webcrawler/app/spider"
)

type (
	// Crawler 抓取器的接口定义
	Crawler interface {
		// Init 初始化采集器
		Init(spider spider.Spider, pipeline pipeline.Pipeline, maxProcesses int) Crawler
		// Start 用种子URL启动采集器，至少一个
		Start(seed string, rest ...string)
		// Stop 停止采集器
		Stop()
		// HasEnd 是否以停止
		HasEnd() bool
	}

	// crawler 抓取器接口的默认实现
	crawler struct {
		sp           spider.Spider
		dl           *downloader.Downloader
		sc           *scheduler.Scheduler
		pl           pipeline.Pipeline
		maxProcesses int
		ch           chan struct{}
	}
)

// NewCrawler 新建一个Crawler实例
func NewCrawler() Crawler {
	return &crawler{}
}

func (c *crawler) Init(spider spider.Spider, pipeline pipeline.Pipeline, maxProcesses int) Crawler {
	c.sp = spider
	c.pl = pipeline
	c.sc = scheduler.NewScheduler()
	c.maxProcesses = 3
	if maxProcesses > 0 {
		c.maxProcesses = maxProcesses
	}
	return c
}

func (c *crawler) Start(seed string, rest ...string) {
	c.sc.Push(seed)
	for _, surl := range rest {
		c.sc.Push(surl)
	}

	if c.ch != nil {
		return
	}
	c.ch = make(chan struct{}, c.maxProcesses)
	for c.ch != nil {
		c.ch <- struct{}{}
		go c.work()
	}
	for {
	}
}

func (c *crawler) work() {
	defer func() {
		<-c.ch
	}()

	url, ok := c.sc.Pop()
	if !ok {
		return
	}
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
	ch := c.ch
	c.ch = nil
	close(ch)
	c.pl.Close()
}

func (c *crawler) HasEnd() bool {
	return c.sc.UnhanldedCount() == 0
}
