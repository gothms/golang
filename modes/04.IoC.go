package modes

import "errors"

/*
Go编程模式：委托和反转控制

控制反转（Inversion of Control，IoC）
	一种软件设计的方法，它的主要思想是把控制逻辑与业务逻辑分开，不要在业务逻辑里写控制逻辑
	因为这样会让控制逻辑依赖于业务逻辑，而是反过来，让业务逻辑依赖控制逻辑
		Inversion of Control：https://en.wikipedia.org/wiki/Inversion_of_control
		IoC：https://en.wikipedia.org/wiki/Inversion_of_control
	IoC/DIP 其实是一种管理思想：https://coolshell.cn/articles/9949.html
		开关和电灯的例子
		开关就是控制逻辑，电器是业务逻辑。我们不要在电器中实现开关，而是要把开关抽象成一种协议，让电器都依赖它
		这样的编程方式可以有效降低程序复杂度，并提升代码重用度

嵌入和委托
	结构体嵌入
	方法重写
	嵌入结构多态
		也可以使用泛型的 interface{} 来多态，但是需要有一个类型转换

反转控制
	有一个存放整数的数据结构，实现了 Add() 、Delete() 和 Contains() 三个操作
		type IntSet struct {
			data map[int]bool
		}
		func NewIntSet() IntSet {
			return IntSet{make(map[int]bool)}
		}
		func (set *IntSet) Add(x int) {
			set.data[x] = true
		}
		func (set *IntSet) Delete(x int) {
			delete(set.data, x)
		}
		func (set *IntSet) Contains(x int) bool {
			return set.data[x]
		}
	实现 Undo 功能
		再包装一下 IntSet ，变成 UndoableIntSet
		type UndoableIntSet struct { // Poor style
			IntSet    // Embedding (delegation)
			functions []func()
		}

		func NewUndoableIntSet() UndoableIntSet {
			return UndoableIntSet{NewIntSet(), nil}
		}


		func (set *UndoableIntSet) Add(x int) { // Override
			if !set.Contains(x) {
				set.data[x] = true
				set.functions = append(set.functions, func() { set.Delete(x) })
			} else {
				set.functions = append(set.functions, nil)
			}
		}


		func (set *UndoableIntSet) Delete(x int) { // Override
			if set.Contains(x) {
				delete(set.data, x)
				set.functions = append(set.functions, func() { set.Add(x) })
			} else {
				set.functions = append(set.functions, nil)
			}
		}

		func (set *UndoableIntSet) Undo() error {
			if len(set.functions) == 0 {
				return errors.New("No functions to undo")
			}
			index := len(set.functions) - 1
			if function := set.functions[index]; function != nil {
				function()
				set.functions[index] = nil // For garbage collection
			}
			set.functions = set.functions[:index]
			return nil
		}
	解释下这段代码
		在 UndoableIntSet 中嵌入了IntSet ，然后 Override 了 它的 Add()和 Delete() 方法
		Contains() 方法没有 Override，所以，就被带到 UndoableInSet 中来了
		在 Override 的 Add()中，记录 Delete 操作
		在 Override 的 Delete() 中，记录 Add 操作
		在新加入的 Undo() 中进行 Undo 操作
	用这样的方式为已有的代码扩展新的功能是一个很好的选择
		这样，就可以在重用原有代码功能和新的功能中达到一个平衡
	但是，这种方式最大的问题是，Undo 操作其实是一种控制逻辑，并不是业务逻辑
		所以，在复用 Undo 这个功能时，是有问题的，因为其中加入了大量跟 IntSet 相关的业务逻辑

反转依赖
	示例代码如下
		先声明一种函数接口，表示我们的 Undo 控制可以接受的函数签名是什么样的
		有了这个协议之后，实现 Undo 的控制逻辑
		在 IntSet 里嵌入 Undo，接着在 Add() 和 Delete() 里使用方法
	不是由控制逻辑 Undo 来依赖业务逻辑 IntSet，而是由业务逻辑 IntSet 依赖 Undo
	这里依赖的是其实是一个协议，这个协议是一个没有参数的函数数组
*/

// IntSet 业务逻辑
type IntSet struct {
	data map[int]bool
	undo Undo
}

func NewIntSet() IntSet {
	return IntSet{data: make(map[int]bool)}
}
func (set *IntSet) Undo() error {
	return set.undo.Undo()
}
func (set *IntSet) Contains(v int) bool {
	return set.data[v]
}
func (set *IntSet) Add(v int) {
	if !set.Contains(v) {
		set.data[v] = true
		set.undo.Add(func() {
			set.Delete(v)
		})
	} else {
		set.undo.Add(nil)
	}
}
func (set *IntSet) Delete(v int) {
	if set.Contains(v) {
		delete(set.data, v)
		set.undo.Add(func() {
			set.Add(v)
		})
	} else {
		set.undo.Add(nil)
	}
}

// Undo 控制逻辑：IntSet 依赖 Undo（一个没有参数的函数数组 协议）
type Undo []func()

func (undo *Undo) Add(f func()) {
	*undo = append(*undo, f)
}
func (undo *Undo) Undo() error {
	if len(*undo) == 0 {
		return errors.New("no functions to undo")
	}
	i := len(*undo) - 1
	if f := (*undo)[i]; f != nil {
		f()
		(*undo)[i] = nil
	}
	*undo = (*undo)[:i]
	return nil
}
