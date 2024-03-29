package main

import (
	"fmt"
	"unsafe"
)

func main() {
	arr := [3]int{}
	s := arr[:2]
	s = arr[:]
	fmt.Println(s)
	fmt.Printf("%T\n", s)

	s1 := make([]int, 5)
	fmt.Printf("The length of s1: %d\n", len(s1))
	fmt.Printf("The capacity of s1: %d\n", cap(s1))
	fmt.Printf("The value of s1: %d\n", s1)
	s2 := make([]int, 5, 8)
	fmt.Printf("The length of s2: %d\n", len(s2))
	fmt.Printf("The capacity of s2: %d\n", cap(s2))
	fmt.Printf("The value of s2: %d\n", s2)

	s3 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s4 := s3[3:6]
	fmt.Printf("The length of s4: %d\n", len(s4))
	fmt.Printf("The capacity of s4: %d\n", cap(s4))
	fmt.Printf("The value of s4: %d\n", s4)

	//fmt.Println(s3[8])
	fmt.Println(len("测试")) // 6

	fmt.Println(32 * 1024)
	fmt.Println(2 << 10)

	var i int32
	size := unsafe.Sizeof(i)
	fmt.Println(size)
}
