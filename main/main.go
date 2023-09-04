package main

import (
	"fmt"
	"golang/core/basic"
	//"golang/core/basic/internal"
	"os"
)

func main() {
	//osArgs()

	//basic.TestFlag()
	basic.TestFlagUsage()

	// Use of the internal package is not allowed
	//var name string
	//internal.Hello(os.Stdout, name)

	//i := new([3]int)
	//i[0] = 2
	//fmt.Println(i)
	//fmt.Printf("%T", *i)

	//ch1 := make(chan int, 3)
	//ch1 <- 2
	//ch1 <- 1
	//ch1 <- 3
	//elem1 := <-ch1
	//fmt.Printf("The first element received from channel ch1: %v\n",
	//	elem1)
}

// osArgs 编译时，传入 os.Args 参数
// cd main
// go run main.go lee
// 输出：hello,world! lee
func osArgs() {
	// go run main.go lee
	if len(os.Args) > 1 { // hello,world! lee
		fmt.Println("hello,world!", os.Args[1])
	}
	os.Exit(-1) // exit status 4294967295
}
