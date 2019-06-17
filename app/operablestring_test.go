package webcrawler

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

func TestOPStringWraped(t *testing.T) {
	str := "test string"
	wraped := WrapedString(str)
	expect(t, reflect.TypeOf(wraped.Unwrap()), reflect.TypeOf(str))
}

func TestOPStringTrimSpace(t *testing.T) {
	str := "  asdasdads  "
	trimed := WrapedString(str).TrimSpace().Unwrap()
	expect(t, trimed, "asdasdads")
}

func TestOPStringFilterLineBreaks(t *testing.T) {
	str := "asda\nsd\nads"
	filted := WrapedString(str).FilterLineBreaks().Unwrap()
	expect(t, filted, "asdasdads")
}

func TestOPStringTrimSpaceAndFilterLineBreaks(t *testing.T) {
	str := "  asda\nsd\nads  "
	chaininged := WrapedString(str).FilterLineBreaks().TrimSpace().Unwrap()
	expect(t, chaininged, "asdasdads")
}
