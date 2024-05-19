package _1_basic

import (
	"fmt"
	"testing"
)

/*
接口：定义对象之间交互的协议

Go接口与其他主要编程语言的差异：Duck Typing 鸭子类型
	1.接口为非入侵性，实现不依赖于接口定义
	2.所以接口的定义可以包含在接口使用者包内

1.接口编程
2.接口完整性检查
	var _ Interface = (*Implement)(nil)

3.空接口与断言
	空接口可以表示任何类型
	通过断言来将空接口转换为指定类型
		v, ok := p.(int)
		p.(type)

4.Go接口最佳实践
	4.1.倾向于使用小的接口定义，很多接口只包含一个方法
	4.2.较大的接口定义，可以由多个小接口定义组合而成
		type C interface {
			A
			B
		}
		A B：都是接口
	4.3.只依赖于必要功能的最小接口
*/

func TestAssert(t *testing.T) {
	var v interface{} = 10
	//v = "10"
	switch tp := v.(type) {
	case int:
		fmt.Println("int", tp)
	case string:
		fmt.Println("string", tp)
	default:
		fmt.Println("Unknown Type")
	}
}

var _ Interface = (*Implement)(nil)

type Interface interface {
	A()
}
type Implement struct{}

func (i *Implement) A() {}

// Country1 Country 版本 1
type Country1 struct {
	Name string
}
type City1 struct {
	Name string
}
type Printable1 interface {
	PrintStr()
}

func (c Country) PrintStr() {
	fmt.Println(c.Name)
}
func (c City) PrintStr() {
	fmt.Println(c.Name)
}

// WithName 版本 2
type WithName struct {
	Name string
}
type Country2 struct {
	WithName
}
type City2 struct {
	WithName
}
type Printable interface {
	PrintStr()
}

func (w WithName) PrintStr() {
	fmt.Println(w.Name)
}

// Country 版本 3
// Stringable 接口把“业务类型” Country 和 City 和“控制逻辑” Print() 给解耦了
// 于是，只要实现了Stringable 接口，都可以传给 PrintStr() 来使用
// 这种编程模式在 Go 的标准库有很多的示例，最著名的就是 io.Read 和 ioutil.ReadAll 的玩法
// 其中 io.Read 是一个接口，你需要实现它的一个 Read(p []byte) (n int, err error) 接口方法
// 只要满足这个规则，就可以被 ioutil.ReadAll这个方法所使用
// 这就是面向对象编程方法的黄金法则——“Program to an interface not an implementation”
type Country struct { // 业务逻辑
	Name string
}
type City struct {
	Name string
}
type Stringable interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}
func (c City) ToString() string {
	return "City = " + c.Name
}
func PrintStr(p Stringable) { // 控制逻辑
	fmt.Println(p.ToString())
}
