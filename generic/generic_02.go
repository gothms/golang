package generic

/*
驯服泛型：了解类型参数
	在经历了 2022 年 3 月 Go 1.18 版本的泛型落地以及 8 月份 Go 1.19 对泛型问题的一轮修复后，我认为是时候开讲 Go 泛型篇了
	虽说目前的 Go 泛型实现和最后一版的泛型设计方案相比还有差距，依旧不是完全版，还有一些特性没有加入，还有问题亟待解决，但对于入门 Go 泛型语法来说，我认为已经是足够了
	和支持泛型的主流编程语言之间的泛型设计与实现存在差异一样，Go 的泛型与其他主流编程语言的泛型也是不同的
Go 泛型设计方案已经明确不支持的若干特性
	https://github.com/golang/proposal/blob/master/design/43651-type-parameters.md#omissions
	不支持泛型特化（specialization），即不支持编写一个泛型函数针对某个具体类型的特殊版本
	不支持元编程（metaprogramming），即不支持编写在编译时执行的代码来生成在运行时执行的代码
	不支持操作符方法（operator method），即只能用普通的方法（method）操作类型实例（比如：getIndex(k)），而不能将操作符视为方法并自定义其实现，比如一个容器类型的下标访问 c[k]
	不支持变长的类型参数（type parameters）
	...
泛型篇的内容共有三讲
	从泛型的基本语法，也就是类型参数（type parameter）开启驯服泛型之旅
	接下来再搞定泛型的难点定义约束（constraints）
	最后我们再来谈谈 Go 泛型的使用时机

例子：返回切片中值最大的元素
	实现一个函数，该函数接受一个切片作为输入参数，然后返回该切片中值最大的那个元素（并没有明确使用什么元素类型的切片）
	非泛型版本
		maxAny 利用 any、type switch 和类型断言（type assertion）实现了我们预期的目标
			E:\gothmslee\golang\generic\max_any.go
		存在问题：
		若要支持其他元素类型的切片，我们需对该函数进行修改
		maxAny 的返回值类型为 any（interface{}），要得到其实际类型的值还需要通过类型断言转换
		使用 any（interface{}）作为输入参数的元素类型和返回值的类型，由于存在装箱和拆箱操作，其性能与 maxInt 等比起来要逊色不少，实测数据如下：
			PS E:\gothmslee\golang> go test -bench="." generic\test\generic_o2_test.go -benchmem
			goos: windows
			goarch: amd64
			cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
			BenchmarkMaxInt-8       442599296                2.765 ns/op           0 B/op          0 allocs/op
			BenchmarkMaxAny-8       82808566                14.93 ns/op            0 B/op          0 allocs/op
	泛型版本
		func MaxGenerics[T constraints.Ordered](sl []T) T {
			if len(sl) == 0 {
				panic("slice is empty")
			}
			ans := sl[0]
			for _, v := range sl[1:] {
				ans = max(ans, v)
			}
			return ans
		}

		对于 ordered 接口中声明的那些原生类型以及以这些原生类型为底层类型（underlying type）的类型（比如示例中的 myString），maxGenerics 都可以无缝支持
		并且，maxGenerics 返回的类型与传入的切片的元素类型一致，调用者也无需通过类型断言做转换
		通过性能基准测试我们也可以看出，与 maxAny 相比，泛型版本的 maxGenerics 性能要好很多，但与原生版函数如 maxInt 等还有差距
			E:\gothmslee\golang\generic\test\generic_o2_test.go
	Go 泛型十分适合实现一些操作容器类型（比如切片、map 等）的算法
		这也是Go 官方推荐的第一种泛型应用场景，此类容器算法的泛型实现使得容器算法与容器内元素类型彻底解耦！

类型参数（type parameters）
	根据官方说法，由于泛型（generic）一词在 Go 社区中被广泛使用，所以官方也就接纳了这一说法
	但 Go 泛型方案的实质是对类型参数（type parameter）的支持，包括：
		泛型函数（generic function）：带有类型参数的函数
		泛型类型（generic type）：带有类型参数的自定义类型
		泛型方法（generic method）：泛型类型的方法
泛型函数
	MaxGenerics 与普通 Go 函数（ordinary function）相比，至少有两点不同：
		MaxGenerics 函数在函数名称与函数参数列表之间多了一段由方括号括起的代码：[T ordered]
		MaxGenerics 参数列表中的参数类型以及返回值列表中的返回值类型都是 T，而不是某个具体的类型
	maxGenerics 函数原型中多出的这段代码[T ordered]就是 Go 泛型的类型参数列表（type parameters list）
		示例中这个列表中仅有一个类型参数 T，ordered 为类型参数的类型约束（type constraint）
		类型约束之于类型参数，就好比常规参数列表中的类型之于常规参数
	Go 语言规范规定：
		函数的类型参数列表位于函数名与函数参数列表之间，由方括号括起的固定个数的、由逗号分隔的类型参数声明组成
		函数一旦拥有类型参数，就可以用该参数作为常规参数列表和返回值列表中修饰参数和返回值的类型
	按 Go 惯例，类型参数名的首字母通常采用大写形式，并且类型参数必须是具名的
		即便你在后续的函数参数列表、返回值列表和函数体中没有使用该类型参数，也是这样
		和常规参数列表中的参数名唯一一样，在同一个类型参数列表中，类型参数名字也要唯一
	作用域
		常规参数列表中的参数有其特定作用域，即从参数声明处开始到函数体结束
		和常规参数类似，泛型函数中类型参数也有其作用域范围，这个范围从类型参数列表左侧的方括号[开始，一直持续到函数体结束
		类型参数的作用域也决定了类型参数的声明顺序并不重要，也不会影响泛型函数的行为
	调用泛型函数
		和普通函数有形式参数与实际参数一样，类型参数也有类型形参（type parameter）和类型实参（type argument）之分
			其中类型形参就是泛型函数声明中的类型参数，以前面示例中的 maxGenerics 泛型函数为例
			如下面代码，maxGenerics 的类型形参就是 T，而类型实参则是在调用 maxGenerics 时实际传递的类型 int：
			// 泛型函数声明：T为类型形参
			func maxGenerics[T ordered](sl []T) T

			// 调用泛型函数：int为类型实参
			m := maxGenerics[int]([]int{1, 2, -4, -6, 7, 0})
		调用泛型函数与调用普通函数的区别
			在调用泛型函数时，除了要传递普通参数列表对应的实参之外，还要显式传递类型实参，比如这里的 int
			并且，显式传递的类型实参要放在函数名和普通参数列表前的方括号中
		问题
			如果泛型函数的类型形参较多，那么逐一显式传入类型实参会让泛型函数的调用显得十分冗长
				foo[int, string, uint32, float64](1, "hello", 17, 3.14)
			解决方法：函数类型实参的自动推断（function argument type inference）
		函数类型实参的自动推断
			通过判断传递的函数实参的类型来推断出类型实参的类型，从而允许开发者不必显式提供类型实参
			类型实参自动腿短有一个前提，它必须是函数的参数列表中使用了的类型形参
			否则编译器将报无法推断类型实参的错误
		在编译器无法推断出结果时，我们可以给予编译器“部分提示”
			比如既然编译器无法推断出 T 的实参类型，那我们就显式告诉编译器 T 的实参类型，即在泛型函数调用时，在类型实参列表中显式传入 T 的实参类型，但 E 的实参类型依然由编译器自动推断，示例代码如下：
			var s = "hello"
			foo[int](5, s)  //ok
			foo[int,](5, s) //ok
		不能通过返回值类型来推断类型实参
			func foo[T any](a int) T {
				var zero T
				return zero
			}

			var a int = foo(5) // 编译器错误：cannot infer T
			println(a)
	泛型函数实例化（instantiation）
		其实泛型函数调用是一个不同于普通函数调用的过程
		Go 对这段泛型函数调用代码的处理分为两个阶段
			图示：generic_02_instantiation.jpg
			Go 首先会对泛型函数进行实例化（instantiation），即根据自动推断出的类型实参生成一个新函数（当然这一过程是在编译阶段完成的，不会对运行时性能产生影响）
			然后才会调用这个新函数对输入的函数参数进行处理
		也可以用一种更形象的方式来描述上述泛型函数的实例化过程。实例化就好比一家生产“求最大值”机器的工厂，它会根据要比较大小的对象的类型将这样的机器生产出来
			工厂接单：调用 maxGenerics([]int{…})，工厂师傅发现要比较大小的对象类型为 int
			模具检查与匹配：检查 int 类型是否满足模具的约束要求，即 int 是否满足 ordered 约束
				如满足，则将其作为类型实参替换 maxGenerics 函数中的类型形参 T，结果为 maxGenerics[int]
			生产机器：将泛型函数 maxGenerics 实例化为一个新函数，这里将其起名为 maxGenericsInt，其函数原型为 func([]int)int
				本质上 maxGenericsInt := maxGenerics[int]
		实际的 Go 代码也可以真实得到这台新生产出的“机器”
			func TestMaxGenerics(t *testing.T) {
				maxGenericsInt := generic.MaxGenerics[int] // 实例化后得到的新“机器”：maxGenericsInt
				fmt.Printf("%T\n", maxGenericsInt)         // func([]int) int
				genericsInt := maxGenericsInt([]int{1, 2, 3, 4, 7, 8, 9, 0})
				fmt.Println(genericsInt) // 输出：9
			}

			整个过程只需检查传入的函数实参（[]int{1, 2, …}）的类型与 maxGenericsInt 函数原型中的形参类型（[]int）是否匹配即可
		另外要注意，当我们使用相同类型实参对泛型函数进行多次调用时，Go 仅会做一次实例化，并复用实例化后的函数
			maxGenerics([]int{1, 2, -4, -6, 7, 0})
			maxGenerics([]int{11, 12, 14, -36,27, 0}) // 复用第一次调用后生成的原型为func([]int) int的函数
泛型类型
	所谓泛型类型，就是在类型声明中带有类型参数的 Go 类型
		示例
			// maxable_slice.go

			type maxableSlice[T ordered] struct {
				elems []T
			}
		maxableSlice 是一个自定义切片类型，这个类型的特点是总可以获取其内部元素的最大值
			其唯一的要求是其内部元素是可排序的，它通过带有 ordered 约束的类型参数来明确这一要求
			像这样在定义中带有类型参数的类型就被称为泛型类型（generic type）
	泛型类型中，类型参数列表放在类型名字后面的方括号中
		和泛型函数一样，泛型类型可以有多个类型参数，类型参数名通常是首字母大写的，这些类型参数也必须是具名的，且命名唯一
		和泛型函数中类型参数有其作用域一样，泛型类型中类型参数的作用域范围也是从类型参数列表左侧的方括号[开始，一直持续到类型定义结束的位置
		示例
			type Set[T comparable] map[T]struct{}

			type sliceFn[T any] struct {
			  s   []T
			  cmp func(T, T) bool
			}

			type Map[K, V any] struct {
			  root    *node[K, V]
			  compare func(K, K) int
			}

			type element[T any] struct {
			  next *element[T]
			  val  T
			}

			type Numeric interface {
			  ~int | ~int8 | ~int16 | ~int32 | ~int64 |
				~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
				~float32 | ~float64 |
				~complex64 | ~complex128
			}

			type NumericAbs[T Numeric] interface {
			  Abs() T
			}
		泛型类型中的类型参数可以用来作为类型声明中
			字段的类型（比如上面的 element 类型）
			复合类型的元素类型（比如上面的 Set 和 Map 类型）
			或方法的参数和返回值类型（如 NumericAbs 接口类型）等
		如果要在泛型类型声明的内部引用该类型名，必须要带上类型参数
			如上面的 element 结构体中的 next 字段的类型：*element[T]
		类型参数的顺序也要与声明中类型参数列表中的顺序一致
			按照泛型设计方案，如果泛型类型有不止一个类型参数，那么在其声明内部引用该类型名时，不仅要带上所有类型参数，类型参数的顺序也要与声明中类型参数列表中的顺序一致
			不过从实测结果来看，Go 1.19 版本对于下面不符合技术方案的泛型类型声明也并未报错
			type P[T1, T2 any] struct {
				F *P[T2, T1] // 不符合技术方案，但Go 1.19编译器并未报错
			}
	使用泛型类型
		和泛型函数一样，使用泛型类型时也会有一个实例化（instantiation）过程
			var sl = maxableSlice[int]{
				elems: []int{1, 2, -4, -6, 7, 0},
			}

			Go 会根据传入的类型实参（int）生成一个新的类型并创建该类型的变量实例，sl 的类型等价于下面代码：
			type maxableIntSlice struct {
				elems []int
			}
		泛型类型是否可以像泛型函数那样实现类型实参的自动推断呢？
			目前的 Go 1.19 尚不支持，下面代码会遭到 Go 编译器的报错
			var sl = maxableSlice {
				elems: []int{1, 2, -4, -6, 7, 0}, // 编译器错误：cannot use generic type maxableSlice[T ordered] without instantiation
			}
		泛型类型与类型别名
			类型别名与其绑定的原类型是完全等价的，但这仅限于原类型是一个直接类型，即可直接用于声明变量的类型
			那么将类型别名与泛型类型绑定是否可行呢？
				type foo[T1 any, T2 comparable] struct {
					a T1
					b T2
				}

				type fooAlias = foo // 编译器错误：cannot use generic type foo[T1 any, T2 comparable] without instantiation
			编译器报错原因
				泛型类型只是一个生产真实类型的“工厂”，它自身在未实例化之前是不能直接用于声明变量的，因此不符合类型别名机制的要求
				泛型类型只有实例化后才能得到一个真实类型，例如下面的代码就是合法的：
					type fooAlias = foo[int, string]
			只能为泛型类型实例化后的类型创建类型别名，实际上上述 fooAlias 等价于实例化后的类型 fooInstantiation：
				type fooInstantiation struct {
					a int
					b string
				}
		泛型类型与类型嵌入
			类型嵌入是运用 Go 组合设计哲学的一个重要手段。引入泛型类型之后，我们依然可以在泛型类型定义中嵌入普通类型
			比如下面示例中 Lockable 类型中嵌入的 sync.Mutex：
				type Lockable[T any] struct {
					t T
					sync.Mutex
				}

				func (l *Lockable[T]) Get() T {
					l.Lock()
					defer l.Unlock()
					return l.t
				}

				func (l *Lockable[T]) Set(v T) {
					l.Lock()
					defer l.Unlock()
					l.t = v
				}
			在泛型类型定义中，我们也可以将其他泛型类型实例化后的类型作为成员
				改写一下上面的 Lockable，为其嵌入另外一个泛型类型实例化后的类型 Slice[int]：
				代码使用泛型类型名（Slice）作为嵌入后的字段名，并且 Slice[int]的方法 String 被提升为 Lockable 实例化后的类型的方法了
				type Slice[T any] []T

				func (s Slice[T]) String() string {
					if len(s) == 0 {
						return ""
					}
					var result = fmt.Sprintf("%v", s[0])
					for _, v := range s[1:] {
						result = fmt.Sprintf("%v, %v", result, v)
					}
					return result
				}

				type Lockable[T any] struct {
					t T
					Slice[int]
					sync.Mutex
				}

				func main() {
					n := Lockable[string]{
						t:     "hello",
						Slice: []int{1, 2, 3},
					}
					println(n.String()) // 输出：1, 2, 3
				}
			同理，在普通类型定义中，我们也可以使用实例化后的泛型类型作为成员
				type Foo struct {
					Slice[int]
				}

				func main() {
					f := Foo{
						Slice: []int{1, 2, 3},
					}
					println(f.String()) // 输出：1, 2, 3
				}
			此外，Go 泛型设计方案支持在泛型类型定义中嵌入类型参数作为成员
				type Lockable[T any] struct {
					T
					sync.Mutex
				}

				不过，Go 1.19 版本编译上述代码时会针对嵌入 T 的那一行报如下错误：
					编译器报错：embedded field type cannot be a (pointer to a) type parameter
				关于这个错误，Go 官方在其 issue 中给出了临时的结论：暂不支持
				https://github.com/golang/go/issues/49030
泛型方法
	在定义泛型类型的方法时，方法的 receiver 部分不仅要带上类型名称，还需要带上完整的类型形参列表（如 maxableSlice[T]），这些类型形参后续可以用在方法的参数列表和返回值列表中
		func (sl *maxableSlice[T]) max() T {
			if len(sl.elems) == 0 {
				panic("slice is empty")
			}

			max := sl.elems[0]
			for _, v := range sl.elems[1:] {
				if v > max {
					max = v
				}
			}
			return max
		}
	在 Go 泛型目前的设计中，泛型方法自身不可以再支持类型参数了
		func (f *foo[T]) M1[E any](e E) T { // 编译器错误：syntax error: method must have no type parameters
			//... ...
		}
	在泛型方法中，receiver 中某个类型参数如果没有在方法参数列表和返回值中使用，可以用“_”代替，但不能不写
		type foo[A comparable, B any] struct{}

		func (foo[A, B]) M1() { // ok
		}

		或
		func (foo[_, _]) M1() { // ok
		}

		或
		func (foo[A, _]) M1() { // ok
		}

		但
		func (foo[]) M1() { // 错误：receiver部分缺少类型参数

		}
	另外，泛型方法中的 receiver 中类型参数名字可以与泛型类型中的类型形参名字不同，位置和数量对上即可
		type foo[A comparable, B any] struct{}

		func (foo[First, Second]) M1(a First, b Second) { // First对应类型参数A，Second对应类型参数B

		}

小结
	类型参数是 Go 泛型方案的具体实现，通过类型参数，我们可以定义泛型函数、泛型类型以及对应的泛型方法
	泛型函数是带有类型参数的函数，在函数名称与参数列表之间声明的类型参数列表使得泛型函数的运行逻辑与参数 / 返回值类型解耦
		调用泛型函数与普通函数略有不同，泛型函数需要进行实例化后才能生成真正执行的、带有类型信息的函数
		同时，Go 泛型支持的类型实参推断也使得开发者在大多数情况下无需显式传递类型实参，获得与普通函数调用几乎一致的体验
	泛型类型是带有类型参数的类型，泛型类型的类型参数放在类型名称后面的类型参数列表中声明，类型参数后续可以在泛型类型声明中用作成员字段的类型或复合类型成员元素的类型
		不过目前（Go 1.19 版本）Go 尚不支持泛型类型的类型实参的自动推断，我们在泛型类型实例化时需要显式传入类型实参
	与泛型类型绑定的方法被称为泛型方法，泛型方法的参数列表和返回值列表中可以使用泛型类型的类型参数，但泛型方法目前尚不支持声明自己的类型参数列表
	Go 泛型的引入，使得 Go 开发人员在 interface{}之后又拥有了一种编写“通用代码”的手段
		并且这种新手段因其更多在编译阶段的检查而变得更加安全，也因其减少了运行时的额外开销使得代码性能更好

思考
	为什么 Go 在方括号“[]”中声明类型参数，而不是使用其他语言都用的尖括号“<>”呢？
		go语言范型不使用 <>，解析的时候容易与 大于 或者 小于 符号混淆？
*/
