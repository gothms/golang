package modes

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

/*
Go编程模式：修饰器
	Python 修饰器的函数式编程
		https://coolshell.cn/articles/11265.html
		这种模式可以很轻松地把一些函数装配到另外一些函数上，让你的代码更加简单
		也可以让一些“小功能型”的代码复用性更高，让代码中的函数可以像乐高玩具那样自由地拼装
	函数式编程
		https://coolshell.cn/articles/10822.html
		其实，Go 语言的修饰器编程模式，也就是函数式编程的模式
		Go 语言的“糖”不多，而且又是强类型的静态无虚拟机的语言，所以，没有办法做到像 Java 和 Python 那样写出优雅的修饰器的代码

简单示例
	decorator 简单示例
		动用了一个高阶函数 decorator()，在调用的时候，先把 Hello() 函数传进去，然后会返回一个匿名函数
		这个匿名函数中除了运行了自己的代码，也调用了被传入的 Hello() 函数
		这个玩法和 Python 的异曲同工，只不过，有些遗憾的是，Go 并不支持像 Python 那样的 @decorator 语法糖
	decorator 计算运行时间
		有两个 Sum 函数，Sum1() 函数就是简单地做个循环，Sum2() 函数动用了数据公式（注意：start 和 end 有可能有负数）
		代码中使用了 Go 语言的反射机制来获取函数名
		修饰器函数是 timedSumFunc()
	decorator HTTP 相关的一个示例
		WithServerHeader() 函数就是一个 Decorator，它会传入一个 http.HandlerFunc，然后返回一个改写的版本
		用 WithServerHeader() 就可以加入一个 Response 的 Header
		这样的函数我们可以写出好多。如下所示，有写 HTTP 响应头的，有写认证 Cookie 的，有检查认证 Cookie 的，有打日志的...
		在使用上，需要对函数一层层地套起来，看上去好像不是很好看，如果需要修饰器比较多的话，代码就会比较难看了

多个修饰器的 Pipeline
	原来的代码
		http.HandleFunc("/v1/hello", WithServerHeader(WithAuthCookie(hello)))
		http.HandleFunc("/v2/hello", WithServerHeader(WithBasicAuth(hello)))
		http.HandleFunc("/v3/hello", WithServerHeader(WithBasicAuth(WithDebugLog(hello))))
	重构时，我们需要先写一个工具函数，用来遍历并调用各个修饰器（通过一个代理函数）
		type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

		func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
			for i := range decors {
				d := decors[len(decors)-1-i] // iterate in reverse
				h = d(h)
			}
			return h
		}
	可以移除不断的嵌套使用
		http.HandleFunc("/v4/hello", Handler(hello,
				WithServerHeader, WithBasicAuth, WithDebugLog))
	无法做到泛型
		代码耦合了需要被修饰的函数的接口类型，无法做到非常通用

泛型的修饰器
	静态语言
		因为 Go 语言不像 Python 和 Java，Python 是动态语言，而 Java 有语言虚拟机，所以它们可以实现一些比较“变态”的事
		但是，Go 语言是一个静态的语言，这就意味着类型需要在编译时就搞定，否则无法编译
		不过，Go 语言支持的最大的泛型是 interface{} ，还有比较简单的 Reflection 机制
	一个比较通用的修饰器：为了便于阅读，删除了出错判断代码
		如下：泛型的修饰器
			这样写是不是有些“傻”？的确是的。不过，这是我个人在 Go 语言里所能写出来的最好的代码了
		使用
			type MyFoo func(int, int, int) int
			var myfoo MyFoo
			Decorator(&myfoo, foo)
			myfoo(1, 2, 3)
		使用 Decorator() 时，还需要先声明一个函数签名，感觉好傻啊，一点都不泛型，不是吗？如果你不想声明函数签名，就可以这样：
			mybar := bar
			Decorator(&mybar, bar)
			mybar("hello,", "world!")
	Go 语言反射机制
		https://blog.golang.org/laws-of-reflection
*/

// ====================泛型的修饰器====================

func Decorator(decoPtr, fn interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value
	// Decorator() 需要两个参数
	decoratedFunc = reflect.ValueOf(decoPtr).Elem() // 出参 decoPtr ，就是完成修饰后的函数
	targetFunc = reflect.ValueOf(fn)                // 入参 fn ，就是需要修饰的函数

	v := reflect.MakeFunc(targetFunc.Type(), // 动用了 reflect.MakeFunc() 函数，创造了一个新的函数
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("before")
			out = targetFunc.Call(in) // 其中的 targetFunc.Call(in) 调用了被修饰的函数
			fmt.Println("after")
			return
		})

	decoratedFunc.Set(v)
	return
}

func foo(a, b, c int) int {
	fmt.Printf("%d, %d, %d \n", a, b, c)
	return a + b + c
}

func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

// ====================多个修饰器的 Pipeline====================

// HttpHandlerDecorator 函数类型：被调用时，方法名 WithServerHeader 直接作为参数，如下
// http.HandleFunc("/v4/hello", Handler(hello, WithServerHeader, WithBasicAuth, WithDebugLog))
type HttpHandlerDecorator func(handlerFunc http.HandlerFunc) http.HandlerFunc

// Handler 实现多个修饰器的 Pipeline
func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}

// ====================decorator HTTP 相关的一个示例====================

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader()")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithAuthCookie()")
		cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
		http.SetCookie(w, cookie)
		h(w, r)
	}
}

func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithBasicAuth()")
		cookie, err := r.Cookie("Auth")
		if err != nil || cookie.Value != "Pass" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithDebugLog")
		r.ParseForm()
		log.Println(r.Form)
		log.Println("path", r.URL.Path)
		log.Println("scheme", r.URL.Scheme)
		log.Println(r.Form["url_long"])
		for k, v := range r.Form {
			log.Println("key:", k)
			log.Println("val:", strings.Join(v, ""))
		}
		h(w, r)
	}
}
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

func main04() {
	http.HandleFunc("/v1/hello", WithServerHeader(WithAuthCookie(hello)))
	http.HandleFunc("/v2/hello", WithServerHeader(WithBasicAuth(hello)))
	http.HandleFunc("/v3/hello", WithServerHeader(WithBasicAuth(WithDebugLog(hello))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader()")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

func main03() {
	http.HandleFunc("/v1/hello", WithServerHeader(hello))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// ====================decorator 计算运行时间====================

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {

		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n",
				getFunctionName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		sum += i
	}
	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

func main02() {
	sum1 := timedSumFunc(Sum1)
	sum2 := timedSumFunc(Sum2)
	// --- Time Elapsed (main.Sum1): 5.2109ms ---
	// --- Time Elapsed (main.Sum2): 0s ---
	// 49999954995000, 49999954995000
	fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))
}

// ====================decorator 简单示例====================

func decorator(f func(s string)) func(s string) {

	return func(s string) {
		fmt.Println("Started")
		f(s)
		fmt.Println("Done")
	}
}
func Hello(s string) {
	fmt.Println(s)
}
func main01() {
	decorator(Hello)("Hello, World!")
}
