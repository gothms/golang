package generic

/*
驯服泛型：定义泛型约束

类型参数约束（type parameter constraint）
	虽然泛型是开发人员表达“通用代码”的一种重要方式，但这并不意味着所有泛型代码对所有类型都适用
	更多的时候，我们需要对泛型函数的类型参数以及泛型函数中的实现代码设置限制
	泛型函数调用者只能传递满足限制条件的类型实参，泛型函数内部也只能以类型参数允许的方式使用这些类型实参值
	在 Go 泛型语法中，我们使用类型参数约束（type parameter constraint）（以下简称约束）来表达这种限制条件
	函数普通参数在函数实现代码中可以表现出来的性质与可以参与的运算由参数类型限制，而泛型函数的类型参数就由约束（constraint）来限制
interface 关键字
	参见：E:\gothmslee\golang\generic\generic_01.go Go 泛型设计的简史
	使用 interface 类型作为约束的定义方法能够最大程度地复用已有语法，并抑制语言引入泛型后的复杂度
	但原有的 interface 语法尚不能满足定义约束的要求。所以，在 Go 泛型版本中，interface 语法也得到了一些扩展

	Go 原生内置的约束、如何定义自己的约束、新引入的类型集合概念等

最宽松的约束：any
	无论是泛型函数还是泛型类型，其所有类型参数声明中都必须显式包含约束，即便你允许类型形参接受所有类型作为类型实参传入也是一样
		表达“所有类型”这种约束，可以使用空接口类型（interface{}）来作为类型参数的约束
		不足：
			如果存在多个这类约束时，泛型函数声明部分会显得很冗长，比如上面示例中的 doSomething 的声明部分
			interface{} 包含 {} 这样的符号，会让本已经很复杂的类型参数声明部分显得更加复杂
			和 comparable、Sortable、ordered 这样的约束命名相比，interface{} 作为约束的表意不那么直接
		为此，Go 团队在 Go 1.18 泛型落地的同时又引入了一个预定义标识符：any。any 本质上是 interface{} 的一个类型别名
	any 约束的类型参数意味着可以接受所有类型作为类型实参。在函数体内，使用 any 约束的形参 T 可以用来做如下操作：
		声明变量
		同类型赋值
		将变量传给其他函数或从函数返回
		取变量地址
		转换或赋值给 interface{} 类型变量
		用在类型断言或 type switch 中
		作为复合类型中的元素类型
		传递给预定义的函数，比如 new
		func doSomething[T1, T2 any](t1 T1, t2 T2) T1 {
			var a T1        // 声明变量
			var b T2
			a, b = t1, t2   // 同类型赋值
			_ = b

			f := func(t T1) {
			}
			f(a)            // 传给其他函数

			p := &a         // 取变量地址
			_ = p

			var i interface{} = a  // 转换或赋值给interface{}类型变量
			_ = i

			c := new(T1)    // 传递给预定义函数
			_ = c

			f(a)            // 将变量传给其他函数

			sl := make([]T1, 0, 10) // 作为复合类型中的元素类型
			_ = sl

			j, ok := i.(T1) // 用在类型断言中
			_ = ok
			_ = j

			switch i.(type) { // 作为type switch中的case类型
			case T1:
			case T2:
			}
			return a        // 从函数返回
		}
	如果对 any 约束的类型参数进行了非上述允许的操作，比如相等性或不等性比较，那么 Go 编译器就会报错：
		func doSomething[T1, T2 any](t1 T1, t2 T2) T1 {
			var a T1
			if a == t1 { // 编译器报错：invalid operation: a == t1 (incomparable types in type set)
			}

			if a != t1 { // 编译器报错：invalid operation: a != t1 (incomparable types in type set)
			}
			... ...
		}
支持比较操作的内置约束：comparable
	Go 编译器会在编译期间判断某个类型是否实现了 comparable 接口
		根据其注释说明，所有可比较的类型都实现了 comparable 这个接口，包括：布尔类型、数值类型、字符串类型、指针类型、channel 类型、元素类型实现了 comparable 的数组和成员类型均实现了 comparable 接口的结构体类型
	comparable 虽然也是一个 interface，但它不能像普通 interface 类型那样来用，比如下面代码会导致编译器报错：
		var i comparable = 5 // 编译器错误：cannot use type comparable outside a type constraint: interface is (or embeds) comparable
		comparable 只能用作修饰类型参数的约束

自定义约束
	凡是接口类型均可作为类型参数的约束
		例如使用 fmt.Stringer 接口作为约束
			func Stringify[T fmt.Stringer](s []T) (ret []string) {}
		一方面，这要求类型参数 T 的实参必须实现 fmt.Stringer 接口的所有方法
		另一方面，泛型函数 Stringify 的实现代码中，声明的 T 类型实例也仅被允许调用 fmt.Stringer 的 String 方法
	扩展 Stringify 这个示例
		将 Stringify 的语义改为只处理非零值的元素
		还要为之加上对排序行为的支持
		1.comparable 虽然不能像普通接口类型那样声明变量，但它却可以作为类型嵌入到其他接口类型中
			自定义了一个 Stringer 接口类型作为约束。在该类型中，我们不仅定义了 String 方法，还嵌入了 comparable
			这样在泛型函数中，我们用 Stringer 约束的类型参数就具备了进行相等性和不等性比较的能力了

			type Stringer interface {
				comparable
				String() string
			}

			func StringifyWithoutZero[T Stringer](s []T) (ret []string) {
				var zero T
				for _, v := range s {
					if v == zero {
						continue
					}
					ret = append(ret, v.String())
				}
				return ret
			}

			type MyString string

			func (s MyString) String() string {
				return string(s)
			}

			func main() {
				sl := StringifyWithoutZero([]MyString{"I", "", "love", "", "golang"}) // 输出：[I love golang]
				fmt.Println(sl)
			}
		2.假如比较 if v == zero || v >= max，编译器会报错：
			invalid operation: v >= max (type parameter T is not comparable with >=)
			Go 编译器认为 Stringer 约束的类型参数 T 不具备排序比较能力
			如果连排序比较性都无法支持，这将大大限制我们泛型函数的表达能力
			但是 Go 又不支持运算符重载（operator overloading），不允许我们定义出下面这样的接口类型作为类型参数的约束：
				type Stringer[T any] interface {
					String() string
					comparable
				  >(t T) bool
				  >=(t T) bool
				  <(t T) bool
				  <=(t T) bool
				}
			于是对 Go 接口类型声明语法做了扩展，支持在接口类型中放入类型元素（type element）信息，比如下面的 ordered 接口类型：
				type ordered interface {
				  ~int | ~int8 | ~int16 | ~int32 | ~int64 |
				  ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
				  ~float32 | ~float64 | ~string
				}
			“|” 和 “~”
				表示以它们为底层类型（underlying type）的类型都满足 ordered 约束，都可以作为以 ordered 为约束的类型参数的类型实参，传入泛型函数
			此时便可以使用 v >= max
				type Stringer interface {
					ordered
					comparable
					String() string
				}
				func StringifyLessThan[T Stringer](s []T, max T) (ret []string) {
					var zero T
					for _, v := range s {
						if v == zero || v >= max {
							continue
						}
						ret = append(ret, v.String())
					}
					return ret
				}
				func main() {
					sl := StringifyLessThan([]MyString{"I", "", "love", "", "golang"}, MyString("cpp")) // 输出：[I]
					fmt.Println(sl)
				}
	Go 接口类型语法的扩展
		图示：generic_03_interface.jpg
		新的接口类型依然可以嵌入其他接口类型，满足组合的设计哲学
			除了嵌入的其他接口类型外，其余的组成元素被称为接口元素（interface element）
		接口元素也有两类
			一类就是常规的方法元素（method element），每个方法元素对应一个方法原型
			另一类则是此次扩展新增的类型元素（type element），即在接口类型中，我们可以放入一些类型信息，就像前面的 ordered 接口那样
		类型元素可以是单个类型，也可以是一组由竖线 “|” 连接的类型，竖线 “|” 的含义是“并”，这样的一组类型被称为 union element
			如果类型中不带有 “~” 符号的类型就代表其自身
			而带有 “~” 符号的类型则代表以该类型为底层类型（underlying type）的所有类型，这类带有 “~” 的类型也被称为 approximation element
		注意：
			union element 中不能包含带有方法元素的接口类型，也不能包含预定义的约束类型，如 comparable
		扩展后，Go 将接口类型分成了两类
			一类是基本接口类型（basic interface type），即其自身和其嵌入的接口类型都只包含方法元素，而不包含类型元素
				基本接口类型不仅可以当做常规接口类型来用，即声明接口类型变量、接口类型变量赋值等，还可以作为泛型类型参数的约束
			除此之外的非空接口类型都属于非基本接口类型，即直接或间接（通过嵌入其他接口类型）包含了类型元素的接口类型
				这类接口类型仅可以用作泛型类型参数的约束，或被嵌入到其他仅作为约束的接口类型中
		示例
			type BasicInterface interface { // 基本接口类型
				M1()
			}

			type NonBasicInterface interface { // 非基本接口类型
				BasicInterface
				~int | ~string // 包含类型元素
			}

			type MyString string

			func (MyString) M1() {
			}

			func foo[T NonBasicInterface](a T) { // 非基本接口类型作为约束
			}

			func bar[T BasicInterface](a T) { // 基本接口类型作为约束
			}

			func main() {
				var s = MyString("hello")
				var bi BasicInterface = s // 基本接口类型支持常规用法
				var nbi NonBasicInterface = s // 非基本接口不支持常规用法，导致编译器错误：cannot use type NonBasicInterface outside a type constraint: interface contains type constraints
				bi.M1()
				nbi.M1()
				foo(s)
				bar(s)
			}
	如何判断一个类型是否满足约束，并作为类型实参传给类型形参呢？
		由于其仅包含方法元素，依旧可以基于方法集合，来确定一个类型是否实现了接口，以及是否可以作为类型实参传递给约束下的类型形参
		但对于只能作为约束的非基本接口类型，既有方法元素，也有类型元素，我们如何判断一个类型是否满足约束，并作为类型实参传给类型形参呢？

类型集合（type set）
	类型集合（type set）的概念是 Go 核心团队在 2021 年 4 月更新 Go 泛型设计方案时引入的
		https://github.com/golang/go/issues/45346
		一旦确定了一个接口类型的类型集合，类型集合中的元素就可以满足以该接口类型作为的类型约束
		也就是可以将该集合中的元素作为类型实参传递给该接口类型约束的类型参数
	结合 Go 泛型设计方案以及 Go 语法规范，我们可以这么来理解类型集合：https://go.dev/ref/spec（Go 语法规范）
		每个类型都有一个类型集合
		非接口类型的类型的类型集合中仅包含其自身
			比如非接口类型 T，它的类型集合为 {T}，即集合中仅有一个元素且这唯一的元素就是它自身
		空接口类型（any 或 interface{}）的类型集合是一个无限集合，该集合中的元素为所有非接口类型
			这个与我们之前的认知也是一致的，所有非接口类型都实现了空接口类型
		非空接口类型的类型集合则是其定义中接口元素的类型集合的交集
			图示：generic_03_type-set.jpg
	接口元素可以是其他嵌入接口类型，可以是常规方法元素，也可以是类型元素
		当接口元素为其他嵌入接口类型时，该接口元素的类型集合就为该嵌入接口类型的类型集合
		而当接口元素为常规方法元素时，接口元素的类型集合就为该方法的类型集合
	一个方法也有自己的类型集合
		Go 规定一个方法的类型集合为所有实现了该方法的非接口类型的集合，这显然也是一个无限集合
			图示：generic_03_type-set_method.jpg
		仅包含多个方法的常规接口类型的类型集合，那就是这些方法元素的类型集合的交集，即所有实现了这三个方法的类型所组成的集合
	每个类型元素的类型集合就是其表示的所有类型组成的集合
		如果是 ~T 形式，则集合中不仅包含 T 本身，还包含所有以 T 为底层类型的类型
		如果使用 Union element，则类型集合是所有竖线 “|” 连接的类型的类型集合的并集
	综合示例
		// TODO
		类型 I 的类型集合
			type Intf1 interface {
				~int | string
			  F1()
			  F2()
			}

			type Intf2 interface {
			  ~int | ~float64
			}

			type I interface {
				Intf1
				M1()
				M2()
				int | ~string | Intf2
			}
		接口类型 I 由四个接口元素组成，分别是 Intf1、M1、M2 和 Union element “int | ~string | Intf2”，我们只要分别求出这四个元素的类型集合，再取一个交集即可
			Intf1 的类型集合
				Intf1 是接口类型 I 的一个嵌入接口，它自身也是由三个接口元素组成，它的类型集合为这三个接口元素的交集
				即 {以 int 为底层类型的所有类型、string、实现了 F1 和 F2 方法的所有类型}
			M1 和 M2 的类型集合
				方法的类型集合是由所有实现该方法的类型组成的
				因此 M1 的方法集合为 {实现了 M1 的所有类型}，M2 的方法集合为 {实现了 M2 的所有类型}
			int | ~string | Intf2 的类型集合
				这是一个类型元素，它的类型集合为 int、~string 和 Intf2 类型集合的并集
				int 类型集合就是 {int}
				~string 的类型集合为 {以 string 为底层类型的所有类型}
				而 Intf2 的类型集合为 {以 int 为底层类型的所有类型，以 float64 为底层类型的所有类型}
			取一下上面集合的交集
				也就是 {以 int 为底层类型的且实现了 F1、F2、M1、M2 这个四个方法的所有类型}
		验证
			func doSomething[T I](t T) {
			}

			type MyInt int

			func (MyInt) F1() {
			}
			func (MyInt) F2() {
			}
			func (MyInt) M1() {
			}
			func (MyInt) M2() {
			}

			func main() {
				var a int = 11
				//doSomething(a) //int does not implement I (missing F1 method)

				var b = MyInt(a)
				doSomething(b) // ok
			}

简化版的约束形式
	在约束对应的接口类型中仅有一个接口元素，且该元素为类型元素时，Go 提供了简化版的约束形式，我们不必将约束独立定义为一个接口类型
		func doSomething[T interface {T1 | T2 | ... | Tn}](t T)

		等价于下面简化版的约束形式：
		func doSomething[T T1 | T2 | ... | Tn](t T)
	特殊情况：定义仅包含一个类型参数的泛型类型时，如果约束中仅有一个 *int 型类型元素
		编译错误
			type MyStruct [T * int]struct{} // 编译错误：undefined: T
											// 编译错误：int (type) is not an expression
		分析
			Go 编译器会将该语句理解为一个类型声明：MyStruct 为新类型的名字，而其底层类型为 [T * int]struct{}
			即一个元素为空结构体类型的数组
		两种解决方案
			完整形式的约束：type MyStruct[T interface{*int}] struct{}
			在简化版约束的 *int 类型后面加上一个逗号：type MyStruct[T *int,] struct{}
约束的类型推断
	在大多数情况下，我们都可以使用类型推断避免在调用泛型函数时显式传入类型实参，Go 泛型可以根据泛型函数的实参推断出类型实参
		但当我们遇到下面示例中的泛型函数时，光依靠函数类型实参的推断是无法完全推断出所有类型实参的：
			func DoubleDefined[S ~[]E, E constraints.Integer](s S) S {
		因为像 DoubleDefined 这样的泛型函数，其类型参数 E 在其常规参数列表中并未被用来声明输入参数，函数类型实参推断仅能根据传入的 S 的类型，推断出类型参数 S 的类型实参，E 是无法推断出来的
	约束类型推断（constraint type inference）
		所以为了进一步避免开发者显式传入类型实参，Go 泛型支持了约束类型推断（constraint type inference）
		即基于一个已知的类型实参（已经由函数类型实参推断判断出来了），来推断其他类型参数的类型
		当通过实参推断得到类型 S 后，Go 会尝试启动约束类型推断来推断类型参数 E 的类型。但约束类型推断可成功应用的前提是 S 是由 E 所表示的

小结
	自定义约束
		因为 Go 不支持操作符重载，单纯依赖基于行为的接口类型 (仅包含方法元素) 作约束是无法满足泛型函数的要求的
		所以“引入” Go 接口类型的扩展语法：支持类型元素
	类型元素
		既有方法元素，也有类型元素，对于作为约束的非基本接口类型，我们就不能像以前那样仅凭是否实现方法集合来判断是否实现了该接口，新的判定手段为类型集合
		类型集合并没有改变什么，只是对哪些类型实现了某接口类型进行了重新解释
		并且，类型集合不是一个运行时概念，我们目前还无法通过运行时反射直观看到一个接口类型的类型集合是什么
	Go 内置了像 any、comparable 的约束
		后续随着 Go 核心团队在 Go 泛型使用上的经验的逐渐丰富，Go 标准库中会增加更多可直接使用的约束
		原计划在 Go 1.18 版本加入 Go 标准库的一些泛型约束的定义暂放在了 Go 实验仓库中，你可以自行参考
		实验仓库：https://github.com/golang/exp/blob/master/constraints/constraints.go

思考
	如果将 Intf1 改为
		type Intf1 interface {
			int | string
		  F1()
		  F2()
		}
	那么接口类型 I 的类型集合变成了什么呢？
		Intf1 的类型集合是个空集合，因为 int 和 string 都没有实现F1和F2方法
*/
