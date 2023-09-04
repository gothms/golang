package main

import (
	"fmt"
	"golang/concurrent/gomem/p1"
)

func init() {
	fmt.Println("init func in main")
}

func main() {
	fmt.Println("V1_p1:", p1.V1_p1)
	fmt.Println("V2_p1:", p1.V2_p1)
}
