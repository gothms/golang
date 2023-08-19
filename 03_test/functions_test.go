package _03_test

import (
	. "golang/03_test"
	"testing"
)

/*
内置单元测试框架
	1.Fail Error：该测试失败，该测试继续，其他测试继续执行
	2.FailNow Fatal：该测试失败，该测试终止，其他测试继续执行

	代码覆盖率
		_test 文件目录下：$ go test -v -cover：-cover
	断言
		Go 中的断言：https://github.com/stretchr/testify
			$ go get -u github.com/stretchr/testify/assert

测试文件包名
	当测试文件和源文件在一个目录的同层级时，包名可以不同，规则示例：
		xx：源文件包名
		xx_test：测试文件包名
	问题
		此时源文件和测试文件在不同包下，必须 大写 暴漏源文件方法，才能调用
*/

func TestSquare(t *testing.T) {
	inputs := [...]int{1, 2, 3}
	expected := [...]int{1, 4, 9}
	for i := 0; i < len(inputs); i++ {
		if sq := Square(inputs[i]); sq != expected[i] {
			t.Errorf("input is %d, the expected is %d, the actual %d\n",
				inputs[i], expected[i], sq)
		}
	}
}

func TestError(t *testing.T) {
	t.Log("Start...")
	t.Error("Error") // 继续执行
	t.Log("End!")
}

func TestFail(t *testing.T) {
	t.Log("Start...")
	t.Fatal("Error") // 终止
	t.Log("End!")
}
