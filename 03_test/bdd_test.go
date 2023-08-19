package _03_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
BDD
	Behavior Driven Development
	BDD是一组编写优秀自动化测试的最佳实践，可以单独使用，但是更多情况下是与TDD 单元测试配合使用的
	BDD解决的一个关键问题就是如何定义TDD或单元测试过程中的细节
	重要是考虑方案而不是实现，它可以引导您设计更好的测试
	需求
		一些不良的单元测试的一个常见问题是过于依赖被测试功能的实现逻辑
		这通常意味着如果你要修改实现逻辑，即使输入输出没有变，通常也需要去更新测试代码
		这就造成了一个问题，让开发人员对测试代码的维护感觉乏味和厌烦
	解决
		BDD则通过向你展示如何测试来解决这个问题，你不需要再面向实现细节设计测试，取而代之的是面向行为来测试
		BDD建议针对行为进行测试，我们不考虑如何实现代码，取而代之的是我们花时间考虑场景是什么，会有什么行为，针对行为代码应该有什么反应

单元测试、TDD、BDD
	https://zhuanlan.zhihu.com/p/91136759
	单元测试回答的是What的问题，TDD回答的是When的问题，BDD回答的是How的问题
	把BDD看作是在需求与TDD之间架起一座桥梁，它将需求进一步场景化，更具体的描述系统应该满足哪些行为和场景，让TDD的输入更优雅、更可靠
	你可以选择单独使用其中一种方法，也可以综合使用这几个方法以取得更好的效果
TDD的一般过程是：
	1.写一个测试
	2.运行这个测试，看到预期的失败
	3.编写尽可能少的业务代码，让测试通过
	4.重构代码
	5.不断重复以上过程
	TDD的最大障碍在于你需要先写测试代码，然后才是产品代码，这是个思维转换和习惯养成的过程，需要不断的重复练习才能逐步掌握

BDD in Go
	项目：https://github.com/smartystreets/goconvey
	安装：$ go get -u github.com/smartystreets/goconvey/convey
		$ go install github.com/smartystreets/goconvey
	启动 Web UI：$GOPATH/bin/goconvey

. "github.com/smartystreets/goconvey/convey"
	. 可以不写 包名

Web 界面
	$ ~/go/bin/goconvey
	window：$ goconvey
默认端口
	http://127.0.0.1:8080
*/

func TestSpec(t *testing.T) {
	Convey("Given 2 even numbers", t, func() {
		a, b := 2, 3
		Convey("When add the two numbers", func() {
			c := a + b
			Convey("Then the result is still even", func() {
				So(c&1, ShouldEqual, 0)
			})
		})
	})
}
