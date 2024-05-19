package generic

/*
驯服泛型：明确使用时机

何时适合使用泛型？
	场景一：编写通用数据结构时
		比如一个先入后出的 stack 数据结构
		第一种方案是为每种要使用的元素类型单独实现一套栈结构
			优点是便于编译器的静态类型检查，保证类型安全，且运行性能很好，因为 Go 编译器可以对代码做出很好的优化
			缺点也很明显，那就是会有大量的重复代码
		第二种方案是使用 interface{} 实现通用数据结构
			在泛型之前，Go 语言中唯一具有“通用”语义的语法就是 interface{} 了
			无论 Go 标准库还是第三方实现的通用数据结构都是基于 interface{} 实现的
			比如下面标准库中 ring 包中 Ring 结构就是使用 interface{} 作为元素类型的：
				// $GOROOT/src/container/ring/ring.go
				type Ring struct {
					next, prev *Ring
					Value      interface{}
				}
			不足：
			Go 编译器无法在编译阶段对进入数据结构中的元素的类型进行静态类型检查
			要想得到元素的真实类型，不可避免要进行类型断言或 type switch 操作
			不同类型数据赋值给 interface{} 或从 interface{} 还原时执行的装箱和拆箱操作带来的额外开销
		比较理想的方案，使用 Go 泛型
			泛型版实现基本消除了前面两种方案的不足
			如果非要说和 IntStack、StringStack 等的差异，那可能就是在执行性能上要差一些了
			泛型版本性能略差与泛型的实现原理有关
	场景二：函数操作的是 Go 原生的容器类型时
		如果函数具有切片、map 或 channel 这些 Go 内置容器类型的参数，并且函数代码未对容器中的元素类型做任何特定假设，那我们使用类型参数可能很有帮助
		类型参数使得此类容器算法与容器内元素类型彻底解耦
		在没有泛型语法之前，实现这样的函数通常需要使用反射
		不过使用反射，会让代码可读性大幅下降，编译器也无法做静态类型检查，并且运行时开销也大得很
	场景三：不同类型实现一些方法的逻辑相同时
		示例：某个函数接受一个自定义接口类型作为参数
			type MyInterface interface {
				M1()
				M2()
				M3()
			}

			func doSomething(i MyInterface) {
			}
		分析
			只有实现了 MyInterface 中全部三个方法的类型，才被允许作为实参传递给 doSomething 函数
			当这些类型实现 M1、M2 和 M3 的逻辑看起来都相同时，我们就可以使用类型参数来帮助实现 M1~M3 这些方法了
		泛型改造：一
			type commonMethod[T any] struct{}

			func (commonMethod[T]) M1() {}
			func (commonMethod[T]) M2() {}
			func (commonMethod[T]) M3() {}

			func main() {
				var intThings commonMethod[int]
				var stringThings commonMethod[string]
				doSomething(intThings)
				doSomething(stringThings)
			}

			使用不同类型，比如 int、string 等作为 commonMethod 的类型实参就可以得到相应实现了 M1~M3 的类型的变量，比如 intThings、stringThings
			这些变量可以直接作为实参传递给 doSomething 函数
		泛型改造：二
			再封装一个泛型函数来简化上述调用
			func doSomethingCM[T any]() {
				doSomething(commonMethod[T]{})
			}

			func main() {
				doSomethingCM[int]()
				doSomethingCM[string]()
			}

			doSomethingCM 泛型函数将 commonMethod 泛型类型实例化与调用 doSomething 函数的过程封装到一起，使得 commonMethod 泛型类型的使用进一步简化了
		Go 标准库的 sort.Sort 就是这样的情况，其参数类型为 sort.Interface，而 sort.Interface 接口中定义了三个方法：
			func Sort(data Interface)

			type Interface interface {
			  Len() int
			  Less(i, j int) bool
			  Swap(i, j int)
			}
		思考
			所有实现 sort.Interface 类型接口的类型，在实现 Len、Less 和 Swap 这三个通用方法的逻辑看起来都相同
			在这样的情况下，我们就可以通过类型参数来给出这三个方法的通用实现
			// TODO
		注意
			如果多个类型实现上述方法的逻辑并不相同，那么我们就不应该使用类型参数
	使用泛型的时机，如果非要总结为一条：
		如果你发现自己多次编写完全相同的代码，其中副本之间的唯一区别是代码使用不同的类型，那么可考虑使用类型参数了

Go 泛型实现原理简介
	泛型窘境
		参见：E:\gothmslee\golang\generic\generic_01.go
		C 语言路径：不实现泛型，不会引入复杂性，但这会“拖慢程序员”，因为可能需要程序员花费精力做很多重复实现
		C++ 语言路径：就像 C++ 的泛型实现方案那样，通过增加编译器负担为每个类型实参生成一份单独的泛型函数的实现
			这种方案产生了大量的代码，其中大部分是多余的，有时候还需要一个好的链接器来消除重复的拷贝，显然这个实现路径会“拖慢编译器”
		Java 路径：就像 Java 的泛型实现方案那样，通过隐式的装箱和拆箱操作消除类型差异，虽然节省了空间，但代码执行效率低，即“拖慢执行性能”
	Go 泛型的实现方案
		Go 核心团队在评估 Go 泛型实现方案时是非常谨慎的，负责泛型实现设计的 Keith Randall 博士一口气提交了三个实现方案，供大家讨论和选择：
		Stenciling 方案
			https://github.com/golang/proposal/blob/master/design/generics-implementation-stenciling.md
		Dictionaries 方案
			https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries.md
		GC Shape Stenciling 方案
			https://github.com/golang/proposal/blob/master/design/generics-implementation-gcshape.md
	Stenciling 方案
		Stenciling 方案也称为模板方案， 它也是 C++、Rust 等语言使用的实现方案
			其主要思路就是在编译阶段，根据泛型函数调用时类型实参或约束中的类型元素，为每个实参类型或类型元素中的类型生成一份单独实现
		图示：generic_04_stenciling.jpg
			Go 编译器为每个调用生成一个单独的函数副本（图中函数名称并非真实的，仅为便于说明而做的命名），相同类型实参的函数只生成一次，或通过链接器消除不同包的相同函数实现
			图示的这一过程在其他编程语言中也被称为“单态化（monomorphization）”。单态是相对于泛型函数的参数化多态（parametric polymorphism）而言的
		Randall 博士也提到了这种方案的不足，那就是拖慢编译器
			泛型函数需要针对不同类型进行单独编译并生成一份独立的代码。如果类型非常多，那么编译出来的最终文件可能会非常大
			同时由于 CPU 缓存无法命中、指令分支预测等问题，可能导致生成的代码运行效率不高
			对于性能不高持保留态度，模板方案在其他编程语言中基本上是没有额外的运行时开销的，并且是应该是对编译器优化友好的
			很多面向系统编程的语言都选择该方案，比如 C++、D 语言、Rust 等
	Dictionaries 方案
		与 Stenciling 方案的实现思路正相反，它不会为每个类型实参单独创建一套代码，反之它仅会有一套函数逻辑
			但这个函数会多出一个参数 dict，这个参数会作为该函数的第一个参数，这和 Go 方法的 receiver 参数在方法调用时自动作为第一个参数有些类似
			这个 dict 参数中保存泛型函数调用时的类型实参的类型相关信息
		图示：generic_04_dictionaries.jpg
			包含类型信息的字典是 Go 编译器在编译期间生成的，并且被保存在 ELF 的只读数据区段（.data）中
			传给函数的 dict 参数中包含了到特定字典的指针
			从方案描述来看，每个 dict 中的类型信息还是十分复杂的
		这种方案也有自身的问题
			比如字典递归的问题，如果调用某个泛型函数的类型实参有很多，那么 dict 信息也会过多等等
			更重要的是它对性能可能有比较大的影响
				比如通过 dict 的指针的间接类型信息和方法的访问导致运行时开销较大
				再比如，如果泛型函数调用时的类型实参是 int，那么如果使用 Stenciling 方案，我们可以通过寄存器复制即可实现 x=y 的操作，但在 Dictionaries 方案中，必须通过 memmove 了
	Go 最终采用的方案：GC Shape Stenciling 方案
		它基于 Stenciling 方案，但又没有为所有类型实参生成单独的函数代码，而是以一个类型的 GC Shape 为单元进行函数代码生成
			一个类型的 GC Shape 是指该类型在 Go 内存分配器 / 垃圾收集器中的表示，这个表示由类型的大小、所需的对齐方式以及类型中包含指针的部分所决定
			这样一来势必就有 GC Shape 相同的类型共享一个实例化后的函数代码，那么泛型调用时又是如何区分这些类型的呢？
		泛型调用时如何区分这些类型
			答案就是字典
			该方案同样在每个实例化后的函数代码中自动增加了一个 dict 参数，用于区别 GC Shape 相同的不同类型
			可见，GC Shape Stenciling 方案本质上是 Stenciling 方案和 Dictionaries 方案的混合版，它也是 Go 1.18 泛型最终采用的实现方案
			为此 Go 团队还给出一个更细化、更接近于实现的 GC Shape Stenciling 实现方案
			新版 GC Shape 方案：https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md
		图示：generic_04_gc-shape-stenciling.jpg
		那么如今的 Go 版本（Go 1.19.x）究竟会为哪些类型实例化出一份独立的函数代码呢？
			示例：声明了一个简单的泛型函数 f，然后分别用不同的 Go 原生类型、自定义类型以及指针类型作为类型实参对 f 进行调用
				func f[T any](t T) T {
					var zero T
					return zero
				}

				type MyInt int

				func main() {
					f[int](5)
					f[MyInt](15)
					f[int64](6)
					f[uint64](7)
					f[int32](8)
					f[rune](18)
					f[uint32](9)
					f[float64](3.14)
					f[string]("golang")

					var a int = 5
					f[*int](&a)
					var b int32 = 15
					f[*int32](&b)
					var c float64 = 8.88
					f[*float64](&c)
					var s string = "hello"
					f[*string](&s)
				}
			通过工具为上述 goshape.go 生成的汇编代码如下：
				generic_04_gc-shape-stenciling_test.jpg
				Go 编译器为每个底层类型相同的类型生成一份函数代码，像 MyInt 和 int、rune 和 int32
				对于所有指针类型，像上面的 *float64、*int 和 *int32，仅生成一份名为 main.f[go.shape.*uint8_0] 的函数代码
			这与新版 GC Shape 方案中的描述是一致的：
				“我们目前正在以一种相当精细的方式实现 GC Shapes。当且仅当两个具体类型具有相同的底层类型或者它们都是指针类型时，它们才会在同一个 GC Shape 分组中”

泛型对执行效率的影响
	Go 泛型实现选择了一条折中的路线
		既没有选择纯 Stenciling 方案，避免了对 Go 编译性能带去较大影响
		也没有选择像 Java 那样泛型那样的纯装箱和拆箱方案，给运行时带去较大开销
	测试：GC Shape + Dictionaries 的混合方案
		PS E:\gothmslee\golang> go test -bench="." generic\test\generic_04_test.go -benchmem
		goos: windows
		goarch: amd64
		cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
		BenchmarkAddInt-8               1000000000               0.2533 ns/op          0 B/op          0 allocs/op
		BenchmarkAddIntGeneric-8        1000000000               0.2543 ns/op          0 B/op          0 allocs/op
	在 Go 1.20 版本中，由于将使用 Unified IR（中间代码表示）替换现有的 IR 表示
		Go 泛型函数的执行性能将得到进一步优化，上述的 benchmark 中两个函数的执行性能将不分伯仲
		Go 1.19 中也可使用 GOEXPERIMENT=unified 来开启 Unified IR 试验性功能
			$GOEXPERIMENT=unified go test -bench .
	建议在一些性能敏感的系统中，还是要慎用尚未得到足够性能优化的泛型
		而在性能不那么敏感的情况下，在符合前面泛型使用时机的时候，我们还是可以大胆使用泛型语法的

小结
	为了防止 Gopher 滥用泛型，给出了几个 Go 泛型最适合应用的场景
		编写通用数据结构时
		编写操作 Go 原生容器类型时
		以及不同类型实现一些方法的逻辑看起来相同时
		除此之外的其他场景下，如果你要使用泛型，务必慎重并深思熟虑
	Go 核心团队在 Go 泛型实现方案的选择上也是煞费苦心，最终选择了 GC Shape Stenciling 的混合方案
		目前这个方案很大程度避免了对 Go 编译性能的影响，但对 Go 泛型代码的执行效率依然存在不小影响
		相信经过几个版本打磨和优化后，Go 泛型的执行性能会有提升，甚至能接近于非泛型的单态版
		从目前看，本讲中的内容仅针对 Go 1.18 和 Go 1.19 的 GC Shape Stenciling 方案适用

思考
	为 Go 标准库 sort.Interface 接口类型提供一个像文中示例 common_method.go 中那样的通用方法的泛型实现
		E:\gothmslee\golang\generic\gc_shape.go
		go API
		E:\Go\src\slices\sort.go
*/
