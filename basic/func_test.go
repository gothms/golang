package basic

import (
	"fmt"
	"testing"
	"time"
)

/*
1.可以有多个返回值
2.所有参数都是值传递：slice，map，channel 会有传引用的错觉
3.函数可以作为变量的值
4.函数可以作为参数和返回值

functional programming：计算机程序的构造和解释
*/
func TestFn(t *testing.T) {
	sf := timeSpent(slowFunc)
	t.Log(sf(10))
}

// timeSpent f 的定义比较长，不便于阅读（参考自定义类型 type IntConv func(op int) int）
func timeSpent(f func(op int) int) func(int) int {
	return func(i int) int {
		start := time.Now()
		ret := f(i)
		fmt.Println("time spent:", time.Since(start).Seconds())
		return ret
	}
}
func slowFunc(op int) int {
	time.Sleep(time.Second * 2)
	return op
}
