package test

import (
	"golang/concurrent"
	"testing"
)

func TestSemaWorkerPool(t *testing.T) {
	concurrent.SemaWorkerPool()
}
