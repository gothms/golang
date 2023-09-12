package main

import "fmt"

func main() {
	ts := TestStructB{}
	a := ts.A() // panic: runtime error: invalid memory address or nil pointer dereference
	fmt.Println(a)
}

type TestInterface interface {
	A() <-chan string
	B() chan string
}
type TestStructA struct {
	TestInterface
}

var _ TestInterface = (*TestStructA)(nil)

func (ti *TestStructA) A() <-chan string {
	fmt.Println("a")
	return nil
}

type TestStructB struct {
	TestStructA
}

var _ TestInterface = (*TestStructB)(nil)
