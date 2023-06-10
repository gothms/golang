package basic

import (
	"fmt"
	"testing"
)

/*
1.数据共享
2.append()：只要扩容，地址就变了
	在 cap 不够用的时候，会重新分配内存以扩大容量
3.Limited Capacity
	arr1=arr[:i:i]：后一个参数 i 叫Limited Capacity，后续的 append() 操作会导致重新分配内存

*/
func TestSlice(t *testing.T) {
	// 1
	s := make([]int, 5)
	s[3] = 9
	s1 := s[1:4]
	s1[0] = 4
	fmt.Println(s, s1) // [0 4 0 9 0] [4 0 9]

	// 2
	a := make([]int, 8)
	b := a[1:5]
	a = append(a, 1)
	a[2] = 3
	fmt.Println(a, b) // [0 0 3 0 0 0 0 0 1] [0 0 0 0]

	// 3
	arr := []int{1, 2, 3, 4, 5}
	arr2 := arr[:3:3]
	arr2 = append(arr2, 0)
	fmt.Println(arr) // [1 2 3 4 5]

	arr1 := arr[:3]
	arr1 = append(arr1, 0)
	fmt.Println(arr) // [1 2 3 0 5]
}
