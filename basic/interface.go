package main

import "fmt"

/*
1.接口编程
2.接口完整性检查
	var _ Interface = (*Implement)(nil)
*/
func main() {

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
