package pipeline

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

func Test_Item(t *testing.T) {
	kayVals := map[string]interface{}{
		"test": "hahaha",
	}
	it := NewItem("test", make([]string, 0), kayVals, make([]string, 0))
	expect(t, it.TableName(), "test")
	expect(t, len(it.SortedKeys()), 0)
	expect(t, len(it.KeyValues()), 1)
	expect(t, len(it.OutputTypes()), 0)
}

func Test_HandlerFunc(t *testing.T) {
	var pp = HandlerFunc(func(first Item, rest ...Item) error {
		expect(t, "我被执行了", "我被执行了")
		return nil
	})
	kayVals := map[string]interface{}{
		"test": "hahaha",
	}
	it := NewItem("test", make([]string, 0), kayVals, make([]string, 0))
	pp.Write(it)
	pp.Close()
}
