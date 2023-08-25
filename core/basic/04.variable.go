package basic

import (
	"flag"
	"fmt"
)

/*
程序实体的那些事儿（上）

程序实体
	Go 语言中的程序实体包括变量、常量、函数、结构体和接口
	Go 语言是静态类型的编程语言，所以我们在声明变量或常量的时候，都需要指定它们的类型，或者给予足够的信息
	在 Go 语言中，变量的类型可以是其预定义的那些类型，也可以是程序自定义的函数、结构体或接口
	常量的合法类型不多，只能是那些 Go 语言预定义的基本类型
问题：声明变量有几种方式？
典型回答
	var
	=
		直接赋予变量类型，这就是“推断”一词所指代的操作了
		它真正的好处，往往会体现在我们写代码之后的那些事情上，比如代码重构
	:=
		短变量声明，实际上就是 Go 语言的类型推断再加上一点点语法糖
		只能在函数体内部使用短变量声明
问题解析
	Go 语言中的类型推断，以及它在代码中的基本体现，另一个是短变量声明的用法
	简单地说，类型推断是一种编程语言在编译期自动解释表达式类型的能力
		Go语言规范
		表达式：https://golang.google.cn/ref/spec#Expressions
		表达式语句：https://golang.google.cn/ref/spec#Expression_statements

知识扩展
1. Go 语言的类型推断可以带来哪些好处？
	它真正的好处，往往会体现在我们写代码之后的那些事情上，比如代码重构
	重构
		通常把不改变某个程序与外界的任何交互方式和规则，而只改变其内部实现”的代码修改方式，叫做对该程序的重构
		重构的对象可以是一行代码、一个函数、一个功能模块，甚至一个软件系统
		这是一个关于程序灵活性的质变
	示例：mainTest
		我们不显式地指定变量name的类型，使得它可以被赋予任何类型的值
		也就是说，变量name的类型可以在其初始化时，由其他程序动态地确定
		在你改变getTheFlag函数的结果类型之后，Go 语言的编译器会在你再次构建该程序的时候，自动地更新变量name的类型
	Go vs 动态类型编程语言
		通过这种类型推断，你可以体验到动态类型编程语言所带来的一部分优势，即程序灵活性的明显提升
		但在那些编程语言中，如Python或Ruby，这种提升可以说是用程序的可维护性和运行效率换来的
		Go 语言是静态类型的，所以一旦在初始化变量时确定了它的类型，之后就不可能再改变。这就避免了在后面维护程序时的一些问题
		重要，这种类型的确定是在编译期完成的，因此不会对程序的运行效率产生任何影响
	小结
		Go 语言的类型推断可以明显提升程序的灵活性，使得代码重构变得更加容易
		同时又不会给代码的维护带来额外负担（实际上，它恰恰可以避免散弹式的代码修改），更不会损失程序的运行效率
2. 变量的重声明是什么意思？
	通过使用短变量声明，我们可以对同一个代码块中的变量进行重声明
	代码块
		在 Go 语言中，代码块一般就是一个由花括号括起来的区域，里面可以包含表达式和语句
		一个代码块可以有若干个子代码块；但对于每个代码块，最多只会有一个直接包含它的代码块（后者可以简称为前者的外层代码块）
	代码块大小
		Go 语言本身以及我们编写的代码共同形成了一个非常大的代码块，也叫全域代码块
			这主要体现在，只要是公开的全局变量，都可以被任何代码所使用
		相对小一些的代码块是代码包，一个代码包可以包含许多子代码包，所以这样的代码块也可以很大
		每个源码文件也都是一个代码块
		每个函数也是一个代码块
		每个if语句、for语句、switch语句和select语句都是一个代码块。甚至，switch或select语句中的case子句也都是独立的代码块
		一对紧挨着的花括号也是代码块，叫“空代码块”
	变量重声明：对已经声明过的变量再次声明。变量重声明的前提条件如下：
		1. 由于变量的类型在其初始化时就已经确定了，所以对它再次声明时赋予的类型必须与其原本的类型相同，否则会产生编译错误
		2. 变量的重声明只可能发生在某一个代码块中
			如果与当前的变量重名的是外层代码块中的变量，那么就是另外一种含义了
		3. 变量的重声明只有在使用短变量声明时才会发生，否则也无法通过编译
			如果要在此处声明全新的变量，那么就应该使用包含关键字var的声明语句，但是这时就不能与同一个代码块中的任何变量有重名了
		4. 被“声明并赋值”的变量必须是多个，并且其中至少有一个是新的变量。这时我们才可以说对其中的旧变量进行了重声明
		变量重声明其实算是一个语法糖（或者叫便利措施）。它允许我们在使用短变量声明时不用理会被赋值的多个变量中是否包含旧变量

总结
	var
		可以被用在任何地方
		无法对已有的变量进行重声明，也就是说它无法处理新旧变量混在一起的情况
	短变量声明
		只能被用在函数或者其他更小的代码块中
	共同点
		基于类型推断，Go 语言的类型推断只应用在了对变量或常量的初始化方面

思考
	如果与当前的变量重名的是外层代码块中的变量，那么这意味着什么？
*/

// 重构示例
func mainTest() {
	var name = getTheFlag()
	flag.Parse()
	fmt.Printf("Hello, %v!\n", name)
}

// 可以随意改变getTheFlag函数的内部实现，及其返回结果的类型，而不用修改main函数中的任何代码
func getTheFlag() *string {
	return flag.String("name", "everyone", "The greeting object.")
}