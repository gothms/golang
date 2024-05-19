package _1_basic

import (
	"fmt"
	"testing"
)

/*
&^：按位置零
	a &^ 1：第一位清0
	a &^ 10：第二位清0，第一位保留
*/

func TestBitClear(t *testing.T) {
	a := 7
	a &^= 5
	fmt.Println(a) // 2
}
