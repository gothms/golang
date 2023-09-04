package _1_basic

import (
	"testing"
)

/*
defer规则
	1.延迟函数的参数在defer语句出现时就已经确定下来了
	2.延迟函数执行按后进先出顺序执行，即先出现的defer最后执行
	3.延迟函数可能操作主函数的具名返回值
*/

func TestDefer(t *testing.T) {
	i := 0
	defer t.Log(i)                       // 0
	defer func() { t.Log("func:", i) }() // 5：闭包
	for ; i < 5; i++ {
		//defer fmt.Println(i)                       // 0
		//defer func() { fmt.Println("func:", i) }() // 5：闭包
	}
}
