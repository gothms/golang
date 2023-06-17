package main

import (
	"fmt"
	"os"
)

func main() {
	// go run main.go lee
	if len(os.Args) > 1 { // hello,world! lee
		fmt.Println("hello,world!", os.Args[1])
	}
	os.Exit(-1) // exit status 4294967295
}
