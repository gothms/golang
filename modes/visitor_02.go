package modes

import "fmt"

type VisitorFunc func(*Info, error) error

//type VisitorFunc func(Visitor, error) error

// Visitor
// 1.Visitor 模式定义
// 首先，kubectl 主要是用来处理 Info结构体
//
// 有一个 VisitorFunc 的函数类型的定义
// 一个 Visitor 的接口，其中需要 Visit(VisitorFunc) error 的方法（这就像是我们上面那个例子的 Shape）
// 最后，为Info 实现 Visitor 接口中的 Visit() 方法，实现就是直接调用传进来的方法（与前面的例子相仿）
type Visitor interface {
	Visit(VisitorFunc) error
}

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

// NameVisitor
// 2.Name Visitor 再来定义几种不同类型的 Visitor
// 这个 Visitor 主要是用来访问 Info 结构中的 Name 和 NameSpace 成员
//
// 声明了一个 NameVisitor 的结构体，这个结构体里有一个 Visitor 接口成员，这里意味着多态
// 在实现 Visit() 方法时，调用了自己结构体内的那个 Visitor的 Visitor() 方法，这其实是一种修饰器的模式，用另一个 Visitor 修饰了自己
// 关于修饰器模式，E:\gothmslee\golang\modes\07.decoration.go
type NameVisitor struct {
	visitor Visitor
}

// Visit 用到修饰器模式：07.decoration
func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

// OtherThingsVisitor
// 3.Other Visitor 这个 Visitor 主要用来访问 Info 结构中的 OtherThings 成员
type OtherThingsVisitor struct {
	visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

// LogVisitor
// 4.Log Visitor
type LogVisitor struct {
	visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

// Visitor02Test
// 5.使用方代码
//
// Visitor 们一层套一层
// 我用 loadFile 假装从文件中读取数据
// 最后执行 v.Visit(loadfile)，这样，我们上面的代码就全部开始激活工作了
//
// 这段代码输出如下的信息
// LogVisitor() before call function
// NameVisitor() before call function
// OtherThingsVisitor() before call function
// ==> OtherThings=We are running as remote team.
// OtherThingsVisitor() after call function
// ==> Name=Hao Chen, NameSpace=MegaEase
// NameVisitor() after call function
// LogVisitor() after call function
func Visitor02Test() {
	info := Info{}
	var v Visitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	v.Visit(loadFile)
}

// 6.上面的代码有以下几种功效：
// 解耦了数据和程序
// 使用了修饰器模式
// 还做出了 Pipeline 的模式

// DecoratedVisitor 用 修饰器模式 重构
// 7.重构一下上面的代码
// 需要注意的是，这个DecoratedVisitor 同样可以成为一个 Visitor 来使用
//
// 用一个 DecoratedVisitor 的结构来存放所有的VistorFunc函数
// NewDecoratedVisitor 可以把所有的 VisitorFunc转给它，构造 DecoratedVisitor 对象
// DecoratedVisitor实现了 Visit() 方法，里面就是来做一个 for-loop，顺着调用所有的 VisitorFunc
type DecoratedVisitor struct {
	visitor    Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}
	return DecoratedVisitor{v, fn}
}

// Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		if err := fn(info, nil); err != nil {
			return err
		}
		for i := range v.decorators {
			if err := v.decorators[i](info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

// Visitor02TestDecoratedVisitor
// 8.DecoratedVisitor 测试
func Visitor02TestDecoratedVisitor() {
	//info := Info{}
	//var v Visitor = &info
	//v = NewDecoratedVisitor(v, NameVisitor, OtherVisitor)
	//
	//v.Visit(LoadFile)
}
