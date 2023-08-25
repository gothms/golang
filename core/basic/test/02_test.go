package test

import (
	"golang/core/basic"
	"testing"
)

func TestFlag(t *testing.T) {
	basic.TestFlag()
	//fmt.Println("test")
}

// TestFlagUsage 需要在 main 中测试，在 test 中测试无效
//func TestFlagUsage(t *testing.T) {
//	basic.TestFlagUsage()
//}
