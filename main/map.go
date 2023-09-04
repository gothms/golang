package main

import (
	"fmt"
	"time"
)

func main() {
	m := map[int]string{
		1: "haha",
	}

	go read(m)
	time.Sleep(time.Second)
	go write(m)
	time.Sleep(30 * time.Second)
	fmt.Println(m)
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
