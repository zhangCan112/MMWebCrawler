package scheduler

import (
	"reflect"
	"testing"
)

/* 测试辅助函数 */
func expect(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestSchedulerPushPop(t *testing.T) {
	url := "http://www.baidu.com"
	DefaultScheduler.Push(url)

	expect(t, len(DefaultScheduler.store), 1)
	expect(t, DefaultScheduler.store[0], url)

	//push重复url将被忽略
	DefaultScheduler.Push(url)
	expect(t, len(DefaultScheduler.store), 1)
	expect(t, DefaultScheduler.store[0], url)

	//添加一个新的后，pop出来的是先进入的
	DefaultScheduler.Push("http://www.163.com")
	popurl, _ := DefaultScheduler.Pop()
	expect(t, len(DefaultScheduler.store), 1)
	expect(t, popurl, url)

	//push 已经执行过的将被忽略
	lastLen := len(DefaultScheduler.store)
	DefaultScheduler.Done(popurl)
	DefaultScheduler.Push(popurl)
	expect(t, len(DefaultScheduler.store), lastLen)
}

func TestSchedulerPopNull(t *testing.T) {
	DefaultScheduler = NewScheduler()
	url, ok := DefaultScheduler.Pop()
	expect(t, ok, false)
	expect(t, url, "")
}

func TestSchedulerDone(t *testing.T) {
	url := "http://www.baidu.com"
	DefaultScheduler.Done(url)

	expect(t, len(DefaultScheduler.history), 1)
	expect(t, DefaultScheduler.history[0], url)

	//再次done 将被忽略
	DefaultScheduler.Done(url)
	expect(t, len(DefaultScheduler.history), 1)
	expect(t, DefaultScheduler.history[0], url)
}

func TestSchedulerHasDone(t *testing.T) {
	url := "http://www.baidu.com"
	DefaultScheduler.Done(url)

	expect(t, DefaultScheduler.HasDone(url), true)
}

func TestSchedulerUnhanldedCount(t *testing.T) {

	url := "http://www.hao123.com"
	DefaultScheduler.Push(url)

	expect(t, DefaultScheduler.UnhanldedCount(), 1)

	//再次done 将被忽略
	DefaultScheduler.Pop()
	expect(t, DefaultScheduler.UnhanldedCount(), 0)
}
