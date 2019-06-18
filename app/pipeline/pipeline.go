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

// NewItem 新建一个Item接口的默认实例
func NewItem(tableName string, sortedKeys []string, keyValues map[string]interface{}, outputTypes []string) Item {
	return &item{
		tableName:   tableName,
		sortedKeys:  sortedKeys,
		keyValues:   keyValues,
		outputTypes: outputTypes,
	}
}

type item struct {
	tableName   string
	sortedKeys  []string
	keyValues   map[string]interface{}
	outputTypes []string
}

func (it *item) TableName() string {
	return it.tableName
}

func (it *item) SortedKeys() []string {
	return it.sortedKeys[0:]
}

func (it *item) KeyValues() map[string]interface{} {
	var kvs = make(map[string]interface{})
	for key, val := range it.keyValues {
		kvs[key] = val
	}
	return kvs
}

func (it *item) OutputTypes() []string {
	return it.outputTypes[0:]
}

// Writer 写入器
type Writer interface {
	// Write 将数据放入输出管道
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

// HandlerFunc 就是一个允许普通函数做为Pipeline的适配器，
type HandlerFunc func(first Item, rest ...Item) error

//Write Writer
func (p HandlerFunc) Write(first Item, rest ...Item) error {
	return p(first, rest...)
}

//Close Closer
func (p HandlerFunc) Close() {
}

// Collector Pipeline接口的扩展实现，实现了数据单元的分类，缓存和批量保存数据的功能
type Collector struct {
	Writer
	cacheSize int
	classed   map[string][]Item
}

// NewCollector 将一个Pipeline包装为Collector
func NewCollector(w Writer, cacheSize int) *Collector {
	return &Collector{
		cacheSize: cacheSize,
		classed:   make(map[string][]Item),
	}
}

// Write Pipeline
func (c *Collector) Write(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)

	for _, it := range total {
		c.classed[it.TableName()] = append(c.classed[it.TableName()], it)
	}

	for key, its := range c.classed {
		if len(its) >= c.cacheSize {
			c.Writer.Write(its[0], its[1:]...)
			delete(c.classed, key)
		}
	}

	return nil
}

// Close Pipeline
func (c *Collector) Close() {
	for _, its := range c.classed {
		if len(its) > 0 {
			c.Writer.Write(its[0], its[1:]...)
		}
	}
	c.classed = make(map[string][]Item)
}

// Mux Pipeline接口的扩展实现，实现了将Pipeline按照outputType进行映射聚合的功能
type Mux struct {
	Pipeline
	store sync.Map
}

// Write Writer
func (m *Mux) Write(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)
	for _, it := range total {
		for _, output := range it.OutputTypes() {
			if val, ok := m.store.Load(output); ok {
				p := val.(Pipeline)
				p.Write(it)
			}
		}
	}
	return nil
}

// Handle 注册一个pipeline到Mux
func (m *Mux) Handle(outputType string, p Pipeline) {
	m.store.Store(outputType, p)
}
