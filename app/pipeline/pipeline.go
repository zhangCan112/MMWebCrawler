package pipeline

// Item 保存数据的一个最小单元
type Item interface {
	// TableName 所属数据表
	TableName() string
	// SortedKeys item的键名限制和顺序
	SortedKeys() []string
	// KeyValues 获取数据键值对
	KeyValues() map[string]interface{}
	// 指定输出的类型，可以指定多个类型
	OutputTypes() []string
}

// Pipeline 输出管道
type Pipeline interface {
	// Put 将数据放入输出管道
	Put(first Item, rest ...Item) error
}

// HandlerFunc 就是一个允许普通函数做为Pipeline的适配器，
type HandlerFunc func(first Item, rest ...Item) error

//Put Pipeline
func (p HandlerFunc) Put(first Item, rest ...Item) error {
	return p(first, rest...)
}

// Collector Pipeline接口的扩展实现，实现了缓存数据并批量保存数据的功能
type Collector struct {
	Pipeline
	cacheSize int
	cache     []Item
}

// Put Pipeline
func (c *Collector) Put(first Item, rest ...Item) error {
	if c.cache == nil {
		c.cache = make([]Item, c.cacheSize)[0:0]
	}
	c.cache = append(c.cache, first)
	c.cache = append(c.cache, rest...)
	if len(c.cache) > c.cacheSize {
		err := c.Pipeline.Put(c.cache[0], c.cache[1:]...)
		if err == nil {
			c.cache = c.cache[0:0]
		}
		return err
	}
	return nil
}
