package pipeline

import "sync"

// Item 保存数据的一个最小单元
type Item interface {
	// TableName 所属数据表
	TableName() string
	// SortedKeys item的键名限制和顺序
	SortedKeys() []string
	// KeyValues 获取数据键值对
	KeyValues() map[string]interface{}
	// OutputTypes 指定输出的类型，可以指定多个类型
	OutputTypes() []string
}

// Writer 写入器
type Writer interface {
	// Put 将数据放入输出管道
	Write(first Item, rest ...Item) error
}

// Closer 关闭
type Closer interface {
	Close()
}

// Pipeline 输出管道
type Pipeline interface {
	Writer
	Closer
}

// HandlerFunc 就是一个允许普通函数做为Writer的适配器，
type HandlerFunc func(first Item, rest ...Item) error

//Write Writer
func (p HandlerFunc) Write(first Item, rest ...Item) error {
	return p(first, rest...)
}

// Collector Pipeline接口的扩展实现，实现了数据单元的分类，缓存和批量保存数据的功能
type Collector struct {
	cacheSize int
	classed   sync.Map
}

// NewCollector 将一个Pipeline包装为Collector
func NewCollector(w Writer, cacheSize int) *Collector {
	return &Collector{
		cacheSize: cacheSize,
	}
}

// Put Pipeline
func (c *Collector) Put(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)

	for _, it := range total {
		c.classed.Store(it.TableName(), it)
	}

	if len(c.cache) > c.cacheSize {
		err := c.Pipeline.Put(c.cache[0], c.cache[1:]...)
		if err == nil {
			c.cache = c.cache[0:0]
		}
		return err
	}
	return nil
}

// Mux Pipeline接口的扩展实现，实现了将Pipeline按照outputType进行映射聚合的功能
type Mux struct {
	Pipeline
	store sync.Map
}

// Put Pipeline
func (m *Mux) Put(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)
	for _, it := range total {
		for _, output := range it.OutputTypes() {
			if val, ok := m.store.Load(output); ok {
				p := val.(Pipeline)
				p.Put(it)
			}
		}
	}
	return nil
}

// Handle 注册一个pipeline到Mux
func (m *Mux) Handle(outputType string, p Pipeline) {
	m.store.Store(outputType, p)
}
