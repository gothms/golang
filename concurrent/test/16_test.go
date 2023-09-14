package test

import (
	"golang/concurrent"
	"testing"
)

func TestStructVal(t *testing.T) {
	concurrent.StructVal()
}

func TestSemaWorkerPool(t *testing.T) {
	concurrent.SemaWorkerPool()
}
