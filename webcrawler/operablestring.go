package webcrawler

import "strings"

// OperableString 对string的包装，用来方便做一些string的链式操作
type OperableString string

// WrapedString 包装string为OperableString
func WrapedString(str string) OperableString {
	return OperableString(str)
}

// Unwrap 打开包装返回string
func (opstr OperableString) Unwrap() string {
	return string(opstr)
}

// TrimSpace 返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串。
func (opstr OperableString) TrimSpace() OperableString {
	return WrapedString(strings.TrimSpace(opstr.Unwrap()))
}

// FilterLineBreaks 过滤掉换行符
func (opstr OperableString) FilterLineBreaks() OperableString {
	return WrapedString(strings.Replace(opstr.Unwrap(), "\n", "", -1))
}
