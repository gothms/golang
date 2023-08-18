package api

import (
	"strings"
	"testing"
	"unicode/utf8"
)

/*
EqualFold

	报告s和t，解释为UTF-8字符串，在简单的Unicode大小写折叠下是否相等，这是一种更普遍的大小写不敏感的形式。
	对两个字符串进行大小写非敏感的匹配

utf8.RuneStart

	判断传入字节是否为某个字符的开始

utf8.DecodeRuneInString

	从字符串中解码出第一个字符

utf8.EncodeRune:func EncodeRune(p []byte, r rune) int {

	writes into p (which must be large enough) the UTF-8 encoding of the rune
*/
func TestEqualFold(t *testing.T) {
	s := "abCdEFghiJ够浪XYZ"
	ts := "AbCdefgHIJ够浪Xyz"
	fold := strings.EqualFold(s, ts)
	t.Log(fold)
}
func TestRuneStartAndDecodeRuneInString(t *testing.T) {
	s := "够浪"
	b := []byte(s)
	start := utf8.RuneStart(b[0]) // true
	t.Log(start)
	if start {
		r, size := utf8.DecodeRuneInString(s)
		t.Logf("%c, %d\n", r, size) // 够, 3
	}
	start = utf8.RuneStart(b[1]) // false
	t.Log(start)
	r, size := utf8.DecodeRuneInString(s[1:])
	t.Logf("%c, %d\n", r, size) // �, 1
}
func TestEncodeRune(t *testing.T) {
	s := "够浪"
	rb := [4]byte{}
	encodeRune := utf8.EncodeRune(rb[:], rune(s[0]))
	t.Log(rb, encodeRune)
}
