package test

import (
	"golang/core/practice"
	"testing"
)

func TestBufferContentLeak(t *testing.T) {
	practice.BufferContentLeak()
}
