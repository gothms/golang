package main

import (
	"fmt"
	"time"
)

func main() {
	m := map[int]string{
		1: "haha",
	}

	//go read(m)
	//time.Sleep(time.Second)
	//go write(m)
	//time.Sleep(30 * time.Second)
	//fmt.Println(m)

	//go func() {
	//	for i := 0; i < 10000; i++ {
	//		m[i] = fmt.Sprintf("put_%d", i)
	//	}
	//}()
	for i := 0; i < 10000; i++ {
		m[i] = fmt.Sprintf("put_%d", i)
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func read(m map[int]string) {
	for {
		_ = m[1]
		time.Sleep(1)
	}
}

func write(m map[int]string) {
	for {
		m[1] = "write"
		time.Sleep(1)
	}
}
