package test

import (
	"golang/concurrent"
	"testing"
)

func TestContextTimeout(t *testing.T) {
	concurrent.ContextTimeout()
}
