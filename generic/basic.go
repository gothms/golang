package generic

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

/*
https://xie.infoq.cn/article/4aa886a7e62783c7e3c99caf3?utm_campaign=geektime_search&utm_content=geektime_search&utm_medium=geektime_search&utm_source=geektime_search&utm_term=geektime_search

泛型
	Go 1.18 版本中首次引入了泛型
	泛型编程是一种计算机编程风格，在这种编程风格中，算法的具体类型可以在以后指定
	泛型是一种可以与多种类型结合使用的类型，泛型函数是一种可以与多种类型结合使用的函数
	尽管"泛型"在过去和现在都可以通过 interface{}、反射包或代码生成器在 Go 中实现，但还是要提一下在使用这三种方法之前需要仔细考虑
	既然 Go 泛型已经可用，就可以消除模板代码，不必担心向后兼容问题，同时还能编写可重用、类型安全和可维护的代码
Go 泛型
	为 Go 语言增加了三个主要组件：
		函数和类型的类型参数
		将接口类型定义为类型集，包括没有方法的类型
		类型推导，允许在调用函数时省略类型参数
	Go 1.18 之前一些处理"泛型"的方法
		使用"泛型"代码生成器生成 Go 软件包，如 https://github.com/cheekybits/genny
		使用带有switch语句和类型转换的接口
		使用带有参数验证的反射软件包
	与正式的 Go 泛型相比，这些方法还远远不够，有如下缺点
		使用类型switch和转换时性能较低
		类型安全损耗：接口和反射不是类型安全的，这意味着代码可能会传递任何类型，而这些类型在编译过程中会被忽略，从而在运行时引起panic
		Go 项目构建更复杂，编译时间更长
		可能需要对调用代码和函数代码进行类型断言
		缺乏对自定义派生类型的支持
		代码可读性差（使用反射时更明显）
	目前在 Go 中的泛型实现
		类型安全 （运行时不会丢失类型，也不需要类型验证、切换或转换）
		高性能
		Go IDE 的支持
		向后兼容 （使用 Go 1.18+ 重构后，旧版代码仍可运行）
		对自定义数据类型的高度支持
入门：使用 Go 泛型
	给定一个整型（int 或 in32 或 int64）数组 nums，如果任何值在数组中至少出现两次，则返回 true；如果每个元素都不同，则返回 false
		https://leetcode.com/problems/contains-duplicate
		创建一个 Go 模块：go mod init github.com/username/leetcode1（将 username 替换为 Github 用户名）
	不使用 Go 泛型的情况下解决这个问题
		程序应该检查输入的数组（可以是 INT、INT32 或 INT64），并找出是否有重复数据，如果有则返回 true，否则返回 false
		分别提供了 int、int32 和 int64 类型数据的示例数组
		...
		如果我们要查找 float32、float64 或字符串的重复内容，该怎么办？
		可以为每种类型编写一个实现，为不同类型明确编写多个函数，或者使用接口，或者通过包生成"泛型"代码。这就是"泛型"诞生的过程
		通过泛型，我们可以编写泛型函数来替代多个函数，或使用带有类型转换的接口
		// ====================Leetcode 示例，不使用 Go 泛型====================
	用泛型来重构代码
		// ====================Go 泛型基础知识====================
		// ====================Leetcode 示例，用泛型来重构代码====================

泛型基础知识
	1.类型参数
	2.类型推导
	3.约束
	4.波浪线(Tilde)运算符和基础类型
	5.预定义约束
	6.可比较（comparable）约束
	7.约束类型链和类型推导
	8.多类型参数和约束
1.类型参数
	func funcName[T any](data T) bool
	type typeName[T comparable] map[T]bool
	T 是类型参数，any 是类型参数的约束条件
	类型参数就像一个抽象的数据层，通常用紧跟函数或类型名称的方括号中的大写字母（多为字母 T）来表示
2.类型推导
	泛型函数必须了解其支持的数据类型，才能正常运行
		要点：泛型类型参数的约束条件是在编译时由调用代码确定的代表单一类型的一组类型
	进一步来说，类型参数的约束代表了一系列可允许的类型，但在编译时，类型参数只代表一种类型，因为 Go 是一种强类型的静态检查语言
		提醒：由于 Go 是一种强类型的静态语言，因此会在应用程序编译期间而非运行时检查类型。Go 泛型解决了这个问题
	类型由调用代码类型推导提供，如果泛型类型参数的约束条件不允许使用该类型，代码将无法编译
		func funcName[T any](data T)
		data := []int{1,3,4,4,5,8,7,3,2}
		fn := funcName(data)
		由于类型是通过约束知道的，因此在大多数情况下，编译器可以在编译时推断出参数类型
	通过类型推导，可以避免从调用代码中为泛型函数或泛型类型实例化进行人工类型推导
		注意：如果编译器无法推断类型（即类型推导失败），可以在实例化时或在调用代码中手动指定类型
		可以忽略调用代码中的 [[]int]，因为编译器会推断出[[]int]，但我更倾向于加入[[]int]以提高代码的可读性
3.约束
	在引入泛型之前，Go 接口用于定义方法集。然而，随着泛型约束的引入，接口现在既可以定义类型集，也可以定义方法集
	约束是用于指定允许使用的泛型的接口
		func funcName[T any](data T) bool
		函数中使用了 any 约束
		Pro 提示：除非必要，否则避免使用 any 接口约束
		在底层实现上，any关键字只是一个空接口，这意味着可以用 interface{} 替换，编译时不会出现任何错误
	接口约束允许使用 int、int16、int32 和 int64 类型。这些类型是约束联合体，用管道符 | 分隔类型
		type example interface {
			int | int16 | int32 | int64
		}
	约束在以下几个方面有好处
		通过类型参数定义了一组允许的类型
		明确发现泛型函数的误用
		提高代码可读性
		有助于编写更具可维护性、可重用性和可测试性的代码
	使用约束时有一个小问题
		示例
			type CustomType int16

			func main() {
			  var value CustomType
			  value = 2
			  printValue(value)
			}

			func printValue[T int16](value T) {
			  fmt.Printf("Value %d", value)
			}
		在终端执行 go run custom-generics.go，就会出现这样的错误
			./custom-type-generics.go:10:12: CustomType does not implement int16 (possibly missing ~ for int16 in constraint int16)
			尽管自定义类型 CustomType 是 int16 类型，但 printValue 泛型函数的类型参数约束无法识别
			鉴于函数约束不允许使用该类型，这也是合理的。不过，可以修改 printValue 函数，使其接受我们的自定义类型
		使用管道操作符，我们将自定义类型 CustomType 添加到 printValue 泛型函数类型参数的约束中，现在有了一个联合约束
			func printValue[T int16 | CustomType](value T) {
				fmt.Println(value)
			}
		但是，等等！为什么需要 int16 类型和"int16"类型的约束联合？
			波浪线(Tilde)运算符和基础类型
4.波浪线(Tilde)运算符和基础类型
	Go 1.18 通过波浪线运算符引入了底层类型，波浪线运算符允许约束支持底层类型
		简单来说，~ 告诉约束接受任何 int16 类型以及任何以 int16 作为底层类型的类型
	删除了约束联合，并在约束中的 int16 类型前用 ~ 波浪线运算符替换了 CustomType
		编译器现在可以理解，CustomType 类型之所以可以使用，仅仅是因为它的底层类型是 int16
		type CustomType int16

		func main() {
		  var value CustomType
		  value = 2
		  printValue(value)
		}

		func printValue[T ~int16](value T) {
		  fmt.Printf("Value %d", value)
		}
5.预定义约束
	Go 团队非常慷慨的为我们提供了一个常用约束的预定义包，可在 golang.org/x/exp/constraints 找到
		go get -u golang.org/x/exp/constraints
		E:\gospace\pkg\mod\golang.org\x\exp@v0.0.0-20240506185415-9bf2ced13842\constraints\constraints.go
		记住：不要忘记导入预定义约束包 golang.org/x/exp/constraints
	重构 Leetcode 示例
		// ====================Leetcode 示例，用泛型来重构代码====================
	具体修改为：
		创建允许使用整数、浮点和字符串及其底层类型的接口约束
		使用 go get 将约束包下载到项目中，在终端的 Leetcode 根目录中执行如下指令：
			go get -u golang.org/x/exp/constraints
		添加到项目中后，在主函数上方创建名为 AllowedData 的约束，如下所示：
			type allowedData interface {
			   constraints.Ordered
			}

			constraints.Ordered 是一种约束，允许任何使用支持比较运算符（如 ≤=≥===）的有序类型
			注：可以在泛型函数中使用 constraint.Ordered，而无需创建新的接口约束
		接下来，删除类型 map 中的所有 filterIntX 类型，创建一个名为 filter 的新类型，如下所示，该类型以 T 为类型参数，以 allowedData 为约束条件：
			type filter[T allowedData] map[T]bool
			在泛型类型 filter 前面，声明了 T 类型参数，并指定 map 键只接受类型参数的约束 allowedData 作为键类型
		现在，删除所有 FindDuplicateIntX 函数。然后使用 Go 泛型创建一个新的 FindDuplicate 函数，代码如下：
			func findDuplicate[T allowedData](data []T) bool {
			   inArray := Filter[T]{}
			   for _, datum := range data {
				  if inArray.has(datum) {
					 return true
				  }
				  inArray.add(datum)
			   }
			   return false
			}

			findDuplicate 函数是一个泛型函数，添加了类型参数 T，并在函数名后面的方括号中指定了 allowedData 约束，然后用类型参数 T 定义了切片类型的函数参数，并用类型参数 T 初始化了 inArray
			注：在函数中声明泛型参数时使用方括号
		接下来，更新 has 以及 add 方法，如下所示
			func (r filter[T]) add(datum T) {
			   r[datum] = true
			}

			func (r filter[T]) has(datum T) bool {
			   _, ok := r[datum]
			   return ok
			}

			因为我们在定义类型 filter 时已经声明了约束，因此方法中只包含类型参数
		最后，更新调用 findDuplicateIntX 的调用代码，使用新的泛型函数 findDuplicate
6.可比较（comparable）约束
	可比较约束与相等运算符（即 == 和≠）相关联
	在 Go 1.18 中引入的一个接口，由结构体、指针、接口、管道等类似类型实现
		注：Comparable 不用作任何变量的类型
		func sort[K comparable, T Data](values map[K]T) error {
			for k, t := range values {
				// code
			}
			return nil
		}
7.约束类型链和类型推导
	类型链
		允许一个已定义的类型参数与另一个类型参数复合的做法被称为类型链。当在泛型结构或函数中定义辅助类型时，这种方法就派上用场了
		func example07[U ~T, T any](t U) <-chan T {
			c := make(chan T)
			// ...
			return c
		}
		...
		c := example07(2)
	约束类型推导
		由于 ~T 是类型参数 T 与任意约束条件的复合体，因此在调用 Example 函数时可以推断出类型参数 U
		注：2 是整数，是 T 的底层类型
8.多类型参数和约束
	Go 泛型支持多类型参数，但有一个问题
		示例
			func test08() {
				printValues(1, 2.1, 3, "c")
			}
			func printValues[A, B any, C comparable](a, a1 A, b B, c C) {
				fmt.Println(a, a1, b, c)
			}
		编译失败
			Cannot use '2.1' (type untyped float) as the type A
	分析
		我们到底有没有声明 int 类型？
		在编译过程中，编译器会根据函数括号中的类型参数约束进行推断
		可以看到，a 和 a1 共享同一个类型参数 A，约束条件是 any（允许所有类型）
		编译器会根据调用代码的变量类型进行推断，并在编译过程中使用函数括号中的类型参数约束来检查类型
		可以看到，a 和 a1 具有相同的类型参数 A，并带有 any 约束。因此，a 和 a1 必须具有相同的类型，因为它们在用于类型推导的函数括号中共享相同的类型参数
		尽管类型参数 A 和 B 共享同一个约束条件，但 b 在函数括号中是独立的

何时使用（或不使用）泛型
	总之，请记住一点：大多数用例并不需要 Go 泛型
	一些指导原则：
		何时使用 Go 泛型
			替换多个类型执行相同逻辑的重复代码，或者替换处理切片、映射和管道等多个类型的重复代码
			在处理容器型数据结构（如链表、树和堆）时
			当代码逻辑需要对多种类型进行排序、比较和/或打印时
		何时不使用 Go 泛型
			当 Go 泛型会让代码变得更复杂时
			当指定函数参数类型时
			当有可能滥用 Go 泛型时。避免使用 Go 泛型/类型参数，除非确定有使用多种类型的重复逻辑
			当不同类型的实现不同时
			使用 io.Reader 等读取器时
	局限性
		目前，匿名函数和闭包不支持类型参数
	Go 泛型的测试
		由于 Go 泛型支持编写多种类型的泛型代码，测试用例将与函数支持的类型数量成正比增长

结论
	如果使用得当，Go 泛型的功能会非常强大；但要谨慎，因为能力越大，责任越大
	Go 泛型将提高代码的灵活性和可重用性，同时保持向后兼容，从而为 Go 语言增添价值
	它简单易用，直接明了，学习周期短，练习有助于更好的理解 Go 泛型及其局限性
	过度使用、借用其他语言的泛型实现以及误解会导致 Go 社区出现反模式和复杂性，风险自担
*/

// ====================Go 泛型基础知识====================

// 1.类型参数
type filter01[T comparable] map[T]bool

func findDuplicate01[T any](data T) bool {
	// find duplicate code
	return false
}

// 2.类型推导
func test02() {
	data := []int{1, 3, 4, 4, 5, 8, 7, 3, 2}
	hasDuplicate := findDuplicate01(data) // Type inferred by variable type

	hasDuplicate = findDuplicate01[[]int](data) // Explicitly declare the type to be used
	fmt.Println(hasDuplicate)
}

// 3.约束
type example03 interface { // example: Interface Constraint
	int | int16 | int32 | int64 // permissible types
}

func findDuplicate03[T example03](data T) bool { // example: The type parameter's Constraint
	return false
}

// 4.波浪线(Tilde)运算符和基础类型
func printValue04[T ~int16](value T) {
	fmt.Println(value)
}

// 5.预定义约束
func printValue05[T constraints.Integer](value T) {
	fmt.Println(value)
}

// 6.可比较（comparable）约束
//func sort06[K comparable, T Data](values map[K]T) error {
//	for k, t := range values {
//		// code
//	}
//	return nil
//}

// 7.约束类型链和类型推导
//func example07[U ~T, T any](t U) <-chan T {	// 编译错误？
//	c := make(chan T)
//	// ...
//	return c
//}

// 8.多类型参数和约束
//
//	func test08() {
//		printValues(1, 2.1, 3, "c")
//	}
func printValues[A, B any, C comparable](a, a1 A, b B, c C) {
	fmt.Println(a, a1, b, c)
}

// ====================Leetcode 示例，用泛型来重构代码====================

type allowedData interface {
	constraints.Ordered
}
type filter[T allowedData] map[T]bool

func findDuplicate[T allowedData](data []T) bool {
	inArray := filter[T]{}
	for _, datum := range data {
		if inArray.has(datum) {
			return true
		}
		inArray.add(datum)
	}
	return false
}
func (r filter[T]) has(datum T) bool {
	_, ok := r[datum]
	return ok
}
func (r filter[T]) add(datum T) {
	r[datum] = true
}
func test05() {
	data := []int{1, 3, 4, 4, 5, 8, 7, 3, 2}     // sample array
	data32 := []int32{1, 3, 4, 4, 5, 8, 7, 3, 2} // sample array
	data64 := []int64{1, 3, 4, 4, 5, 8, 7, 3, 2} // sample array
	fmt.Printf("Duplicate found %t\n", findDuplicate(data))
	fmt.Printf("Duplicate found %t\n", findDuplicate(data32))
	fmt.Printf("Duplicate found %t\n", findDuplicate(data64))
}

// ====================Leetcode 示例，不使用 Go 泛型====================

//type FilterInt map[int]bool
//type FilterInt32 map[int32]bool
//type FilterInt64 map[int64]bool
//
//func test01() {
//	data := []int{1, 3, 4, 4, 5, 8, 7, 3, 2}     // sample array
//	data32 := []int32{1, 3, 4, 4, 5, 8, 7, 3, 2} // sample array
//	data64 := []int64{1, 3, 4, 4, 5, 8, 7, 3, 2} // sample array
//	fmt.Printf("Duplicate found %t\n", FindDuplicateInt(data))
//	fmt.Printf("Duplicate found %t\n", FindDuplicateInt32(data32))
//	fmt.Printf("Duplicate found %t\n", FindDuplicateInt64(data64))
//}
//
//func FindDuplicateInt(data []int) bool {
//	inArray := FilterInt{}
//	for _, datum := range data {
//		if inArray.has(datum) {
//			return true
//		}
//		inArray.add(datum)
//	}
//	return false
//}
//
//func FindDuplicateInt32(data []int32) bool {
//	inArray := FilterInt32{}
//	for _, datum := range data {
//		if inArray.has(datum) {
//			return true
//		}
//		inArray.add(datum)
//	}
//	return false
//}
//
//func FindDuplicateInt64(data []int64) bool {
//	inArray := FilterInt64{}
//	for _, datum := range data {
//		if inArray.has(datum) {
//			return true
//		}
//		inArray.add(datum)
//	}
//	return false
//}
//
//func (r FilterInt) add(datum int) {
//	r[datum] = true
//}
//
//func (r FilterInt32) add(datum int32) {
//	r[datum] = true
//}
//
//func (r FilterInt64) add(datum int64) {
//	r[datum] = true
//}
//
//func (r FilterInt) has(datum int) bool {
//	_, ok := r[datum]
//	return ok
//}
//
//func (r FilterInt32) has(datum int32) bool {
//	_, ok := r[datum]
//	return ok
//}
//
//func (r FilterInt64) has(datum int64) bool {
//	_, ok := r[datum]
//	return ok
//}
