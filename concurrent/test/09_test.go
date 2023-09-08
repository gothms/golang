package test

import (
	"golang/concurrent"
	"sync"
	"testing"
)

func TestSyncMap(t *testing.T) {
	var m sync.Map
	m.Store(1, 2)
	d, ok := m.LoadAndDelete(1)
	t.Log(d, ok)
	m.Store(1, 3)
	value, ok := m.Load(1)
	t.Log(value, ok)
}

func TestMap(t *testing.T) {
	concurrent.MapNilValue()
}
func TestMapSyncPanic(t *testing.T) {
	var m = make(map[int]int, 10) // 初始化一个map
	go func() {
		for {
			m[1] = 1 //设置key
		}
	}()
	go func() {
		for {
			_ = m[2] //访问这个map
		}
	}()
	select {}
}
