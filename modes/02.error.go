package modes

/*
Go 编程模式：错误处理

C 语言的错误检查
	基本上来说，其通过函数的返回值标识是否有错，然后通过全局的 errno 变量加一个 errstr 的数组来告诉你为什么出错
	为什么是这样的设计呢？
		道理很简单，除了可以共用一些错误，更重要的是这其实是一种妥协，比如：read()、 write()、 open() 这些函数的返回值其实是返回有业务逻辑的值，也就是说，这些函数的返回值有两种语义：
		1.一种是成功的值，比如 open() 返回的文件句柄指针 FILE*
		2.另一种是错误 NULL。这会导致调用者并不知道是什么原因出错了，需要去检查 errno 以获得出错的原因，从而正确地处理错误
	像atoi()、 atof()、 atol() 或 atoll() 这样的函数，是不会设置 errno的，而且，如果结果无法计算的话，行为是 undefined
		如果一个要转的字符串是非法的（不是数字的格式），如 “ABC” 或者整型溢出了，那么这个函数应该返回什么呢？
		出错返回，返回什么数都不合理，因为这会和正常的结果混淆在一起。比如，如果返回 0，就会和正常的对 “0” 字符的返回值完全混淆在一起，这样就无法判断出错的情况了
	用返回值 + errno 的错误检查方式会有一些问题：
		程序员一不小心就会忘记检查返回值，从而造成代码的 Bug
		函数接口非常不纯洁，正常值和错误值混淆在一起，导致语义有问题

Java 的错误处理
	Java 语言使用 try-catch-finally 通过使用异常的方式来处理错误，其实，这比起 C 语言的错误处理进了一大步
	使用抛异常和抓异常的方式可以让我们的代码有这样一些好处
		函数接口在 input（参数）和 output（返回值）以及错误处理的语义是比较清楚的
		正常逻辑的代码可以跟错误处理和资源清理的代码分开，提高了代码的可读性
		异常不能被忽略（如果要忽略也需要 catch 住，这是显式忽略）
		在面向对象的语言中（如 Java），异常是个对象，所以，可以实现多态式的 catch
		与状态返回码相比，异常捕捉有一个显著的好处，那就是函数可以嵌套调用，或是链式调用

Go 语言的错误处理
	Go 语言的函数支持多返回值，所以，可以在返回接口把业务语义（业务返回值）和控制语义（出错返回值）区分开
	Go 语言的很多函数都会返回 result、err 两个值，于是就有这样几点：
		参数上基本上就是入参，而返回接口把结果和错误分离，这样使得函数的接口语义清晰
		而且，Go 语言中的错误参数如果要忽略，需要显式地忽略，用 _ 这样的变量来忽略
		另外，因为返回的 error 是个接口（其中只有一个方法 Error()，返回一个 string ），所以你可以扩展自定义的错误处理
	一个函数返回了多个不同类型的 error
		if err != nil {
		  switch err.(type) {
			case *json.SyntaxError:
			  ...
			case *ZeroDivisionError:
			  ...
			case *NullPointerError:
			  ...
			default:
			  ...
		  }
		}
Go 语言的错误处理的方式，本质上是返回值检查，但是它也兼顾了异常的一些好处——对错误的扩展
	资源清理
		出错后是需要做资源清理的，不同的编程语言有不同的资源清理的编程模式
		C 语言：使用的是 goto fail; 的方式到一个集中的地方进行清理（推荐一篇有意思的文章《由苹果的低级 BUG 想到的》
			https://coolshell.cn/articles/11112.html
		C++ 语言：一般来说使用 RAII 模式，通过面向对象的代理模式，把需要清理的资源交给一个代理类，然后再析构函数来解决
			https://en.wikipedia.org/wiki/Resource_acquisition_is_initialization
		Java 语言：可以在 finally 语句块里进行清理
		Go 语言：使用 defer 关键词进行清理
	Error Check Hell
		原代码：
			func parse(r io.Reader) (*Point, error) {

				var p Point

				if err := binary.Read(r, binary.BigEndian, &p.Longitude); err != nil {
					return nil, err
				}
				if err := binary.Read(r, binary.BigEndian, &p.Latitude); err != nil {
					return nil, err
				}
				if err := binary.Read(r, binary.BigEndian, &p.Distance); err != nil {
					return nil, err
				}
				if err := binary.Read(r, binary.BigEndian, &p.ElevationGain); err != nil {
					return nil, err
				}
				if err := binary.Read(r, binary.BigEndian, &p.ElevationLoss); err != nil {
					return nil, err
				}
			}
		函数式编程：
			func parse(r io.Reader) (*Point, error) {
				var p Point
				var err error
				read := func(data interface{}) {
					if err != nil {
						return
					}
					err = binary.Read(r, binary.BigEndian, data)
				}

				read(&p.Longitude)
				read(&p.Latitude)
				read(&p.Distance)
				read(&p.ElevationGain)
				read(&p.ElevationLoss)

				if err != nil {
					return &p, err
				}
				return &p, nil
			}
		Go 语言的 bufio.Scanner()
			scanner := bufio.NewScanner(input)

			for scanner.Scan() {
				token := scanner.Text()
				// process token
			}

			if err := scanner.Err(); err != nil {
				// process the error
			}

			scanner在操作底层的 I/O 的时候，那个 for-loop 中没有任何的 if err !=nil 的情况，退出循环后有一个 scanner.Err() 的检查，看来使用了结构体的方式
		模仿：
			首先，定义一个结构体和一个成员函数：
				type Reader struct {
					r   io.Reader
					err error
				}

				func (r *Reader) read(data interface{}) {
					if r.err == nil {
						r.err = binary.Read(r.r, binary.BigEndian, data)
					}
				}
			然后，修改代码：
				func parse(input io.Reader) (*Point, error) {
					var p Point
					r := Reader{r: input}

					r.read(&p.Longitude)
					r.read(&p.Latitude)
					r.read(&p.Distance)
					r.read(&p.ElevationGain)
					r.read(&p.ElevationLoss)

					if r.err != nil {
						return nil, r.err
					}

					return &p, nil
				}
		流式接口 Fluent Interface：https://martinfowler.com/bliki/FluentInterface.html
			package main

			import (
			  "bytes"
			  "encoding/binary"
			  "fmt"
			)

			// 长度不够，少一个Weight
			var b = []byte {0x48, 0x61, 0x6f, 0x20, 0x43, 0x68, 0x65, 0x6e, 0x00, 0x00, 0x2c}
			var r = bytes.NewReader(b)

			type Person struct {
			  Name [10]byte
			  Age uint8
			  Weight uint8
			  err error
			}
			func (p *Person) read(data interface{}) {
			  if p.err == nil {
				p.err = binary.Read(r, binary.BigEndian, data)
			  }
			}

			func (p *Person) ReadName() *Person {
			  p.read(&p.Name)
			  return p
			}
			func (p *Person) ReadAge() *Person {
			  p.read(&p.Age)
			  return p
			}
			func (p *Person) ReadWeight() *Person {
			  p.read(&p.Weight)
			  return p
			}
			func (p *Person) Print() *Person {
			  if p.err == nil {
				fmt.Printf("Name=%s, Age=%d, Weight=%d\n",p.Name, p.Age, p.Weight)
			  }
			  return p
			}

			func main() {
			  p := Person{}
			  p.ReadName().ReadAge().ReadWeight().Print()
			  fmt.Println(p.err)  // EOF 错误
			}
		if err != nil
			这个技巧的使用场景是有局限的，也就只能在对于同一个业务对象的不断操作下可以简化错误处理
			如果是多个业务对象，还是得需要各种 if err != nil的方式
	包装错误
		需要包装一下错误，而不是干巴巴地把err返回到上层，我们需要把一些执行的上下文加入
			通常来说，我们会使用 fmt.Errorf()来完成这个事，比如：
			if err != nil {
			   return fmt.Errorf("something failed: %v", err)
			}
		在 Go 语言的开发者中，更为普遍的做法是将错误包装在另一个错误中，同时保留原始内容：
			type authorizationError struct {
				operation string
				err error   // original error
			}

			func (e *authorizationError) Error() string {
				return fmt.Sprintf("authorization failed during %s: %v", e.operation, e.err)
			}

		更好的方式是通过一种标准的访问方法，这样，我们最好使用一个接口
			比如 causer接口中实现 Cause() 方法来暴露原始错误，以供进一步检查：
			type causer interface {
				Cause() error
			}

			func (e *authorizationError) Cause() error {
				return e.err
			}
		三方错误库：https://github.com/pkg/errors
			import "github.com/pkg/errors"

			//错误包装
			if err != nil {
				return errors.Wrap(err, "read failed")
			}

			// Cause接口
			switch err := errors.Cause(err).(type) {
			case *MyError:
				// handle specifically
			default:
				// unknown error
			}
参考
	https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
	https://blog.golang.org/errors-are-values
*/
