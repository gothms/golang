package modes

import "net/http"

/*
	1.Python 修饰器的函数式编程
		https://coolshell.cn/articles/11265.html
	2.函数式编程
		https://coolshell.cn/articles/10822.html
	3.Go语言的修饰器编程模式，也就是函数式编程模式
	4.Pipeline：多个修饰器
	5.泛型修饰器：通过反射机制实现
*/

//HttpHandlerDecorator 函数类型：被调用时，方法名 WithServerHeader 直接作为参数，如下
//http.HandleFunc("/v4/hello", Handler(hello, WithServerHeader, WithBasicAuth, WithDebugLog))
type HttpHandlerDecorator func(handlerFunc http.HandlerFunc) http.HandlerFunc

// Handler 实现多个修饰器的 Pipeline
func Handler(f http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		fn := decors[len(decors)-1-i]
		f = fn(f)
	}
	return f
}
