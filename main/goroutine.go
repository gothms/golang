package main

import (
	"fmt"
	"runtime"
)

func main() {
	gomaxprocs := runtime.GOMAXPROCS(0) // 8
	fmt.Println(gomaxprocs)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}
