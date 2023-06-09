package main

import (
	"fmt"
	"unicode/utf8"
)

/*
1.utf8.RuneCountInString
	使用 len() 求字符串长度时，一个汉字一般会返回 3 的长度
*/
func main() {
	s := "Go浪"
	l := len(s)
	fmt.Println(l) // 5
	l = LenOfString(s)
	fmt.Println(l) // 3
}
func LenOfString(s string) int {
	return utf8.RuneCountInString(s)
}
