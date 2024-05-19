package modes

/*
Go编程模式：切片、接口、时间和性能

Slice
	append()这个函数在 cap 不够用的时候，就会重新分配内存以扩大容量，如果够用，就不会重新分配内存了！
	示例：
		func main() {
			path := []byte("AAAA/BBBBBBBBB")
			sepIndex := bytes.IndexByte(path,'/')

			dir1 := path[:sepIndex]
			dir2 := path[sepIndex+1:]

			fmt.Println("dir1 =>",string(dir1)) //prints: dir1 => AAAA
			fmt.Println("dir2 =>",string(dir2)) //prints: dir2 => BBBBBBBBB

			dir1 = append(dir1,"suffix"...)

			fmt.Println("dir1 =>",string(dir1)) //prints: dir1 => AAAAsuffix
			fmt.Println("dir2 =>",string(dir2)) //prints: dir2 => uffixBBBB
		}
		解决问题：
		dir1 := path[:sepIndex:sepIndex]
		使用了 Full Slice Expression，最后一个参数叫“Limited Capacity”，于是，后续的 append() 操作会导致重新分配内存
深度比较
	复制一个对象时，这个对象可以是内建数据类型、数组、结构体、Map……
	在复制结构体的时候，如果我们需要比较两个结构体中的数据是否相同，就要使用深度比较，而不只是简单地做浅度比较
	使用反射 reflect.DeepEqual()
接口编程
	“成员函数”（“Receiver”）
		示例
			func PrintPerson(p *Person) {
				fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
			  p.Name, p.Sexual, p.Age)
			}

			func (p *Person) Print() {
				fmt.Printf("Name=%s, Sexual=%s, Age=%d\n",
			  p.Name, p.Sexual, p.Age)
			}

			func main() {
				var p = Person{
					Name: "Hao Chen",
					Sexual: "Male",
					Age: 44,
				}

				PrintPerson(&p)
				p.Print()
			}
		这种方式是一种封装，因为 PrintPerson()本来就是和 Person强耦合的，所以理应放在一起
		更重要的是，这种方式可以进行接口编程，对于接口编程来说，也就是一种抽象，主要是用在“多态”
			https://coolshell.cn/articles/8460.html#%E6%8E%A5%E5%8F%A3%E5%92%8C%E5%A4%9A%E6%80%81
	Go 语言接口的编程模式
		版本一：
			type Country struct {
				Name string
			}

			type City struct {
				Name string
			}

			type Printable interface {
				PrintStr()
			}
			func (c Country) PrintStr() {
				fmt.Println(c.Name)
			}
			func (c City) PrintStr() {
				fmt.Println(c.Name)
			}

			c1 := Country {"China"}
			c2 := City {"Beijing"}
			c1.PrintStr()
			c2.PrintStr()
		版本二：
			使用“结构体嵌入”的方式
			引入一个叫 WithName的结构体，但是这会带来一个问题：在初始化的时候变得有点乱

			type WithName struct {
				Name string
			}

			type Country struct {
				WithName
			}

			type City struct {
				WithName
			}

			type Printable interface {
				PrintStr()
			}

			func (w WithName) PrintStr() {
				fmt.Println(w.Name)
			}

			c1 := Country {WithName{ "China"}}
			c2 := City { WithName{"Beijing"}}
			c1.PrintStr()
			c2.PrintStr()
		版本三：
			使用了一个叫Stringable 的接口，我们用这个接口把“业务类型” Country 和 City 和“控制逻辑” Print() 给解耦了
			于是，只要实现了Stringable 接口，都可以传给 PrintStr() 来使用

			type Country struct {
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
			func (c City) ToString() string{
				return "City = " + c.Name
			}

			func PrintStr(p Stringable) {
				fmt.Println(p.ToString())
			}

			d1 := Country {"USA"}
			d2 := City{"Los Angeles"}
			PrintStr(d1)
			PrintStr(d2)
		Go 标准库
			io.Read 和 ioutil.ReadAll
			其中 io.Read 是一个接口，你需要实现它的一个 Read(p []byte) (n int, err error) 接口方法，只要满足这个规则，就可以被 ioutil.ReadAll这个方法所使用
			这就是面向对象编程方法的黄金法则——“Program to an interface not an implementation”
接口完整性检查
	var _ Shape = (*Square)(nil)
	声明一个 _ 变量（没人用）会把一个 nil 的空指针从 Square 转成 Shape，这样，如果没有实现完相关的接口方法，编译器就会报错
时间
	《你确信你了解时间吗？》《关于闰秒》
		https://coolshell.cn/articles/5075.html
		https://coolshell.cn/articles/7804.html
		时间有时区、格式、精度等问题，其复杂度不是一般人能处理的
	在 Go 语言中，一定要使用 time.Time 和 time.Duration 这两个类型
		在命令行上，flag 通过 time.ParseDuration 支持了 time.Duration
		JSON 中的 encoding/json 中也可以把time.Time 编码成 RFC 3339 的格式
		数据库使用的 database/sql 也支持把 DATATIME 或 TIMESTAMP 类型转成 time.Time
		YAML 也可以使用 gopkg.in/yaml.v2 支持 time.Time 、time.Duration 和 RFC 3339 格式

		如果你要和第三方交互，实在没有办法，也请使用 RFC 3339 的格式
		最后，如果你要做全球化跨时区的应用，一定要把所有服务器和时间全部使用 UTC 时间
性能提示
	1.数字转字符串，使用 strconv.Itoa() 比 fmt.Sprintf 快一倍左右
	2.尽可能避免把 string 转成 []Byte，会导致性能下降
	3.在 for-loop 里对 Slice 使用 append()，请先把 Slice 容量扩充到位
		避免系统自动按 2^n 进行扩展，但又用不到的情况，从而避免浪费内存
	4.使用 StringBuffer 或 StringBuild 拼接字符串，性能比使用 + 或 += 高三到四个数量级
	5.尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作
	6.避免在热代码中进行内存分配，会导致 gc 很忙。尽可能使用 sync.Pool 来重用对象
	7.使用 lock-free 的操作，避免使用 mutex，尽可能使用 sync/Atomic 包。关于无锁编程，参考：
		无锁队列实现：https://coolshell.cn/articles/8239.html
		无锁Hashmap实现：https://coolshell.cn/articles/9703.html
	8.关于 I/O 缓冲，I/O 是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可带来更高的性能
	9.在 for-loop 里的固定的正则表达式，一定使用 regexp.Compile() 编译正则表达式，性能会提升两个数量级
	10.如果需要更高性能的协议，就考虑使用 protobuf 或 msgp，而不是JSON，因为JSON的序列化和反序列化里使用了反射
		protobuf：https://github.com/golang/protobuf
		msgp：https://github.com/tinylib/msgp
	11.使用 map 时，使用整型 key 会比字符串快，因为整型比较比字符串比较快
参考文档
	更多技巧，写出更好的 Go，必读：
	Effective Go：https://golang.org/doc/effective_go.html
	Uber Go Style：https://github.com/uber-go/guide/blob/master/style.md
	50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs：http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	Go Advice：https://github.com/cristaloleg/go-advice
	Practical Go Benchmarks：https://www.instana.com/blog/practical-golang-benchmarks/
	Benchmarks of Go serialization methods：https://github.com/alecthomas/go_serialization_benchmarks
	Debugging performance issues in Go programs：https://github.com/golang/go/wiki/Performance
	Go code refactoring: the 23x performance hunt：https://medium.com/@val_deleplace/go-code-refactoring-the-23x-performance-hunt-156746b522f7
*/
