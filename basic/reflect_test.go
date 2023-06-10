package basic

import (
	"fmt"
	"reflect"
	"testing"
)

/*
1.reflect.DeepEqual()
	比较两个结构体中的数据是否相同，就要使用深度比较，而不只是简单地做浅度比较
*/

func TestReflect(t *testing.T) {
	v1 := data{1, "v1"}
	v2 := data{1, "v1"}
	fmt.Println(v1 == v2) // true
	fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))
	//prints: v1 == v2: true

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	//fmt.Println(m1 == m2) // 编译错误
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))
	//prints: m1 == m2: true

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	//fmt.Println(s1 == s2) // 编译错误
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))
	//prints: s1 == s2: true
}

type data struct {
	v    int
	name string
}
