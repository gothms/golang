package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func decorator(decoPtr, fn interface{}) (err error) {
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
	fmt.Printf("%s %s \n", a, b)
	return a + b
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

func main() {
	sum1 := timedSumFunc(Sum1)
	sum2 := timedSumFunc(Sum2)

	fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))

	//http.HandleFunc("/v1/hello", WithServerHeader(hello))
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	//type MyFoo func(int, int, int) int
	//var myfoo MyFoo
	//decorator(&myfoo, foo)
	//i := myfoo(1, 2, 3)
	//fmt.Println(i)

	mybar := bar
	decorator(&mybar, bar)
	s := mybar("hello,", "world!")
	fmt.Println(s)
}
