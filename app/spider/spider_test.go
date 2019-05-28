package spider

import (
	"reflect"
	"sort"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

func Test_muxEntrySlice_sortable(t *testing.T) {
	val1 := muxEntry{
		pattern: "1",
	}
	val2 := muxEntry{
		pattern: "123",
	}
	slice := []muxEntry{val2, val1}
	sort.Sort(muxEntrySlice(slice))
	expect(t, slice[0].pattern, "1")
	expect(t, slice[1].pattern, "123")
	//结构体判等
	expect(t, slice[0], val1)
	expect(t, slice[1], val2)
}

func Test_Spiderlair(t *testing.T) {
	var spl = NewSpiderlair()

	val1 := muxEntry{
		pattern: "http://ps4.tgbus.com",
		s:       SpiderFunc(func(rw ResultWriter, doc *goquery.Document) {}),
	}
	val2 := muxEntry{
		pattern: "http://switch.tgbus.com",
		s:       SpiderFunc(func(rw ResultWriter, doc *goquery.Document) {}),
	}

	val3 := muxEntry{
		pattern: "http://XBoxOne.tgbus.com",
		s:       SpiderFunc(func(rw ResultWriter, doc *goquery.Document) {}),
	}
	val4 := muxEntry{
		pattern: "http://XBoxX360.tgbus.com",
		s:       SpiderFunc(func(rw ResultWriter, doc *goquery.Document) {}),
	}

	//乱序加入
	spl.Join(val3.pattern, val3.s)
	spl.Join(val2.pattern, val2.s)
	spl.Join(val4.pattern, val4.s)
	spl.Join(val1.pattern, val1.s)

	//保存应当正常
	expect(t, len(spl.m), 4)
	expect(t, len(spl.es), 4)

	//获取的spider有值
	refute(t, spl.Spider(val3.pattern), nil)

	//es众的顺序应当正确
	expect(t, spl.es[0].pattern, val1.pattern)
	expect(t, spl.es[1].pattern, val2.pattern)
	expect(t, spl.es[2].pattern, val3.pattern)
	expect(t, spl.es[3].pattern, val4.pattern)

	//最长匹配原则应当生效
	refute(t, spl.Spider("http://ps4.tgbus.com/test"), nil)

	//删除操作测试
	spl.Clean("http://switch.tgbus.com")
	expect(t, spl.Spider("http://switch.tgbus.com"), nil)
	expect(t, len(spl.es), 3)

	//清除操作测试
	spl.CleanAll()
	expect(t, len(spl.m), 0)
	expect(t, len(spl.es), 0)
}
