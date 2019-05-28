/**
scheduler 是一个FIFO队列，实现了广度优先安排待处理的请求
*/

package scheduler

import (
	"sync"
)

// DefaultScheduler 默认的全局调度器
var DefaultScheduler = NewScheduler()

// Scheduler URl调度器
type Scheduler struct {
	store    []string
	history  []string
	storeMux sync.RWMutex
}

// NewScheduler 返回一个Scheduler实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		store:   make([]string, 100)[0:0],
		history: make([]string, 100)[0:0],
	}
}

// Push 推入一个待处理的请求
func (s *Scheduler) Push(url string) {
	s.storeMux.Lock()
	defer s.storeMux.Unlock()
	if !s.hasDone(url) && !s.hasExist(url) {
		s.store = append(s.store, url)
	}
}

func (s *Scheduler) hasExist(url string) bool {
	var isExist = false
	for _, val := range s.store {
		if val == url {
			isExist = true
			break
		}
	}
	return isExist
}

// Pop 返回一个待处理的请求
func (s *Scheduler) Pop() (string, bool) {
	s.storeMux.Lock()
	defer s.storeMux.Unlock()
	if len(s.store) > 0 {
		url := s.store[0]
		s.store = s.store[1:]
		return url, true
	}

	return "", false
}

// Done 用来标记一个已完成的Request
// 这样这个请求再被push时将被忽略
func (s *Scheduler) Done(url string) {
	s.storeMux.Lock()
	defer s.storeMux.Unlock()
	if !s.hasDone(url) {
		s.history = append(s.history, url)
	}
}

// HasDone url是否已经抓取过
func (s *Scheduler) HasDone(url string) bool {
	s.storeMux.RLock()
	defer s.storeMux.RUnlock()
	return s.hasDone(url)
}

func (s *Scheduler) hasDone(url string) bool {
	var isExist = false
	for _, val := range s.history {
		if val == url {
			isExist = true
			break
		}
	}
	return isExist
}
