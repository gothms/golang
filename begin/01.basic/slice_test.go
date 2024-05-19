package _1_basic

import (
	"fmt"
	"testing"
)

/*
1.共享存储结构
2.append()：只要扩容，地址就变了
	在 cap 不够用的时候，会重新分配内存以扩大容量

3.Limited Capacity
	arr1=arr[:i:i]：后一个参数 i 叫Limited Capacity，后续的 append() 操作会导致重新分配内存

4.for range 性能问题
	for _, v := range arr {
		data = append(data, &v)
	}
	解释为：
	for i, v := 0, 0; i < len(arr); i++ {
		v = arr[i]
		data = append(data, &v)
	}
*/

func TestCapacity(t *testing.T) {
	arr := make([]int, 3, 5)
	arr[0], arr[1], arr[2] = 1, 2, 3
	i := len(arr)
	arr = (arr)[:i+1] // 容量capacity > len，否则越界
	t.Log(arr)

	brr := arr[:]
	brr[2] = 99
	t.Logf("%p, %v\n", brr, brr) // 0xc000014510, [1 2 99 0]
	t.Logf("%p, %v\n", arr, arr) // 0xc000014510, [1 2 99 0]
}

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
func TestRange(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	var data []*int
	for _, v := range arr {
		data = append(data, &v)
	}
	//for i, v := 0, 0; i < len(arr); i++ {
	//	v = arr[i]
	//	data = append(data, &v)
	//}
	for _, p := range data {
		t.Log(*p) // 5,5,5,5,5
	}
}
