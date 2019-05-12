/**
scheduler 是一个先入先出的队列，实现了广度优先安排待处理的请求
*/

package webcrawler

import (
	"sync"
)

var (
	store      = make([]*Request, 100)
	history    = make([]string, 100)
	storeMux   sync.Mutex
	historyMux sync.RWMutex
)

// Push 推入一个待处理的请求
func Push(req *Request) {
	storeMux.Lock()
	defer storeMux.Unlock()
	store = append(store, req)
}

// Pop 返回一个待处理的请求
func Pop() *Request {
	storeMux.Lock()
	defer storeMux.Unlock()
	if len(store) > 0 {
		req := store[0]
		store = store[1:]
		return req
	}

	return nil
}

// Done 用来标记一个已完成的Request
// 这样这个请求再被push时将被忽略
func Done(req *Request) {
	historyMux.Lock()
	defer historyMux.Unlock()
	if !hasDone(req) {
		history = append(history)
	}
}

func HasDone(req *Request) bool {
	return false
}

func hasDone(req *Request) bool {
	return true
}
