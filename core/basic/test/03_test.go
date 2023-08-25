package test

import (
	"golang/core/basic"
	"golang/core/basic/internal"
	"os"
	"testing"
)

func TestHelloInternal(t *testing.T) {
	basic.Hello("lee")
}

// TestInternal 能调用，但怎么输入参数呢？
func TestInternal(t *testing.T) {
	internal.Hello(os.Stdout, os.Args[1])
}
