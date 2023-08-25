package main

import (
	"fmt"
	"golang/core/basic"
	"os"
)

func main() {
	//osArgs()

	//basic.TestFlag()
	basic.TestFlagUsage()

	//i := new([3]int)
	//i[0] = 2
	//fmt.Println(i)
	//fmt.Printf("%T", *i)
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
