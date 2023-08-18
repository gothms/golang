package basic

import "testing"

/*
Is Go an object-oriented language?
	Yes and no. Although Go has types and methods and allows an object-oriented style of programming, there is no type hierarchy.
	The concept of “interface” in Go provides a different approach that we believe is easy to use and in some ways more general.
	There are also ways to embed types in other types to provide something analogous—but not identical—to subclassing.
	Moreover, methods in Go are more general than in C++ or Java: they can be defined for any sort of data, even built-in types such as plain, “unboxed” integers.
	They are not restricted to structs (classes).

	Also, the lack of a type hierarchy makes “objects” in Go feel much more lightweight than in languages such as C++ or Java.

1.值接收器/指针接收器
	值接收器：复制一份结构体，更大的内存开销
2.自定义类型：简单可读
3.内嵌结构体
	不是继承，不支持方法的重载，不支持LSP
	变量是什么类型，就优先调用该类型的方法
*/

type IntConv func(op int) int

func TestStruct(t *testing.T) {

}
