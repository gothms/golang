package test

import (
	"os"
	"testing"
)

/*
1.Test：单元测试
2.Benchmark：性能测试
3.Example：示例测试
4.子测试
5.Main测试
	TestMain 用于主动执行各种测试，可以测试前后做setup和tear-down操作
*/
func TestTest(t *testing.T) {
	a, b := 1, 2
	a, b = b, a
	t.Log(a, b)
}
func TestMain(t *testing.M) {
	println("Test setup.")
	retCode := t.Run() // 执行测试，包括单元测试、性能测试和示例测试
	println("Test tear-down.")
	os.Exit(retCode)
}
