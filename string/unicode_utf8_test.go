package string

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

/*
string

	1.string是数据类型，不是引用/指针类型，默认值是空字符串，而不是 nil
	2.string是只读的byte slice，len函数返回它所包含的byte数
	3.string的byte数组可以存放任何数据

Unicode UTF8

	1.Unicode是一种字符集（code point）
	2.UTF8是Unicode的存储实现（转换为字节序列的规则）
*/
func TestString(t *testing.T) {
	s := "Go浪"
	l := len(s)
	fmt.Println(l) // 5
	l = LenOfString(s)
	fmt.Println(l) // 3

	s = "中"
	c := []rune(s)
	t.Logf("中 unicode %x", c) // 中 unicode [4e2d]
	t.Logf("中 UTF8 %x", s)    // 中 UTF8 e4b8ad
}

// TestStingToRune rune 测试
func TestStingToRune(t *testing.T) {
	s := "Go语言&云原生"
	for _, c := range s {
		t.Logf("%[2]c %[2]x", 'a', c) // [2]：第2个参数
	}
}

// LenOfString utf8.RuneCountInString：使用 len() 求字符串长度时，一个汉字一般会返回 3 的长度
func LenOfString(s string) int {
	return utf8.RuneCountInString(s)
}
