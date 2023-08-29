package test

import (
	"golang/core/advance"
	"testing"
)

func TestPanicTest(t *testing.T) {
	advance.RecoverWrongCall()
}
func TestDeferStack(t *testing.T) {
	advance.DeferStack()
}
