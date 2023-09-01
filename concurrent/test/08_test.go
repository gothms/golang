package test

import (
	"golang/concurrent"
	"sync"
	"testing"
)

func TestDockerOncePanic(t *testing.T) {
	concurrent.DockerOncePanic()
}
func TestNewOnceDemo(t *testing.T) {
	for i := 0; i < 5; i++ {
		demo := concurrent.NewOnceDemo()
		t.Logf("%p\n", &demo)
		//concurrent.NewOnceDemo()
		//t.Logf("%p\n", &concurrent.OnceDemo)
	}
}
func TestOnce(t *testing.T) {
	var once sync.Once
	f1 := func() {
		t.Log("call f1")
	}
	once.Do(f1) // 正常打印
	f2 := func() {
		t.Log("call f2")
	}
	once.Do(f2) // 无输出
}
func TestLazyNewDemo(t *testing.T) {
	concurrent.LazyNewDemo()
}
