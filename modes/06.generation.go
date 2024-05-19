package modes

/*
Go 编程模式：Go Generation
	官方文章：
		The Laws of Reflection：https://blog.golang.org/laws-of-reflection
泛型编程：
	Go 语言的代码生成主要还是用来解决编程泛型的问题
	主要解决问题：
		因为静态类型语言有类型，所以相关的算法或是对数据处理的程序会因为类型不同而需要复制一份
		这会导致数据类型和算法功能耦合
	泛型编程可以解决这样的问题，泛型编程在写代码时，不用关心处理数据的类型，只需要关心相关的处理逻辑
	泛型编程是静态语言中非常非常重要的特性，如果没有泛型，就很难做到多态，也很难完成抽象，导致代码冗余量很大
类比
	螺丝刀本来只有一个拧螺丝的作用，但是因为螺丝的类型太多，有平口的，有十字口的，有六角的……螺丝还有不同的尺寸，这就导致我们的螺丝刀为了要适配各种千奇百怪的螺丝类型（样式和尺寸），也是各种样式的
	而真正的抽象是，螺丝刀不应该关心螺丝的类型，它只要关注自己的功能是不是完备，并且让自己可以适配不同类型的螺丝就行了
	这就是所谓的泛型编程要解决的实际问题

Go语言的类型检查
	两种技术
		Type Assert
		Reflection
	Type Assert
		variable, error := 变量.(type)
		variable：被转好的类型
		error：能否转换
	Reflection
		示例
			type Container struct {
				s reflect.Value
			}
			func NewContainer(t reflect.Type, size int) *Container {
				if size <=0  { size=64 }
				return &Container{
					s: reflect.MakeSlice(reflect.SliceOf(t), 0, size),
				}
			}
			func (c *Container) Put(val interface{})  error {
				if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
					return fmt.Errorf(“Put: cannot put a %T into a slice of %s",
						val, c.s.Type().Elem()))
				}
				c.s = reflect.Append(c.s, reflect.ValueOf(val))
				return nil
			}
			func (c *Container) Get(refval interface{}) error {
				if reflect.ValueOf(refval).Kind() != reflect.Ptr ||
					reflect.ValueOf(refval).Elem().Type() != c.s.Type().Elem() {
					return fmt.Errorf("Get: needs *%s but got %T", c.s.Type().Elem(), refval)
				}
				reflect.ValueOf(refval).Elem().Set( c.s.Index(0) )
				c.s = c.s.Slice(1, c.s.Len())
				return nil
			}
		Reflection 的玩法
			在 NewContainer()时，会根据参数的类型初始化一个 Slice
			在 Put()时，会检查 val 是否和 Slice 的类型一致
			在 Get()时，我们需要用一个入参的方式，因为我们没有办法返回 reflect.Value 或 interface{}，不然还要做 Type Assert
			不过有类型检查，所以，必然会有检查不对的时候，因此，需要返回 error
		弊端
			用反射写出来的代码还是有点复杂的
	他山之石
		对于泛型编程最牛的语言 C++ 来说，这类问题都是使用 Template 解决的

Go Generator
	C++ 的编译器会在编译时分析代码，根据不同的变量类型来自动化生成相关类型的函数或类，在 C++ 里，叫模板的具体化
		这个技术是编译时问题，所以不需要在运行时进行任何的类型识别
		Go 也可以使用这种技术，只是 Go 的编译器不会帮你干，你需要自己动手
	Go 的代码生成，需要三点：
		1)一个函数模板，在里面设置好相应的占位符
		2)一个脚本，用于按规则来替换文本并生成新的代码
		3)一行注释代码
	函数模板
		把上面的示例改成模板，取名为 container.tmp.go 放在 ./template/下
			package PACKAGE_NAME
			type GENERIC_NAMEContainer struct {
				s []GENERIC_TYPE
			}
			func NewGENERIC_NAMEContainer() *GENERIC_NAMEContainer {
				return &GENERIC_NAMEContainer{s: []GENERIC_TYPE{}}
			}
			func (c *GENERIC_NAMEContainer) Put(val GENERIC_TYPE) {
				c.s = append(c.s, val)
			}
			func (c *GENERIC_NAMEContainer) Get() GENERIC_TYPE {
				r := c.s[0]
				c.s = c.s[1:]
				return r
			}
		函数模板中我们有如下的占位符
			PACKAGE_NAME：包名
			GENERIC_NAME ：名字
			GENERIC_TYPE ：实际的类型
	gen.sh的生成脚本
		如下
			#!/bin/bash

			set -e

			SRC_FILE=${1}
			PACKAGE=${2}
			TYPE=${3}
			DES=${4}
			#uppcase the first char
			PREFIX="$(tr '[:lower:]' '[:upper:]' <<< ${TYPE:0:1})${TYPE:1}"

			DES_FILE=$(echo ${TYPE}| tr '[:upper:]' '[:lower:]')_${DES}.go

			sed 's/PACKAGE_NAME/'"${PACKAGE}"'/g' ${SRC_FILE} | \
				sed 's/GENERIC_TYPE/'"${TYPE}"'/g' | \
				sed 's/GENERIC_NAME/'"${PREFIX}"'/g' > ${DES_FILE}
		需要 4 个参数
			模板源文件
			包名
			实际需要具体化的类型
			用于构造目标文件名的后缀
		sed 命令
			用 sed 命令去替换刚刚的函数模板，并生成到目标文件中
			教程：https://coolshell.cn/articles/9104.html
	生成代码
		只需要在代码中打一个特殊的注释
			//go:generate ./gen.sh ./template/container.tmp.go gen uint32 container
			func generateUint32Example() {
				var u uint32 = 42
				c := NewUint32Container()
				c.Put(u)
				v := c.Get()
				fmt.Printf("generateExample: %d (%T)\n", v, v)
			}

			//go:generate ./gen.sh ./template/container.tmp.go gen string container
			func generateStringExample() {
				var s string = "Hello"
				c := NewStringContainer()
				c.Put(s)
				v := c.Get()
				fmt.Printf("generateExample: %s (%T)\n", v, v)
			}
		注释
			第一个注释是生成包名 gen，类型是 uint32，目标文件名以 container 为后缀
			第二个注释是生成包名 gen，类型是 string，目标文件名是以 container 为后缀
		在工程目录中直接执行 go generate 命令，就会生成两份代码
			这两份代码可以让我们的代码完全编译通过，付出的代价就是需要多执行一步 go generate 命令
		一份文件名为 uint32_container.go
			package gen

			type Uint32Container struct {
				s []uint32
			}
			func NewUint32Container() *Uint32Container {
				return &Uint32Container{s: []uint32{}}
			}
			func (c *Uint32Container) Put(val uint32) {
				c.s = append(c.s, val)
			}
			func (c *Uint32Container) Get() uint32 {
				r := c.s[0]
				c.s = c.s[1:]
				return r
			}
		另一份的文件名为 string_container.go
			package gen

			type StringContainer struct {
				s []string
			}
			func NewStringContainer() *StringContainer {
				return &StringContainer{s: []string{}}
			}
			func (c *StringContainer) Put(val string) {
				c.s = append(c.s, val)
			}
			func (c *StringContainer) Get() string {
				r := c.s[0]
				c.s = c.s[1:]
				return r
			}
新版 Filter
	有了这样的技术，我们就不用在代码里，用那些晦涩难懂的反射来做运行时的类型检查了
		我们可以写出很干净的代码，让编译器在编译时检查类型对不对
	一个 Fitler 的模板文件 filter.tmp.go
		package PACKAGE_NAME

		type GENERIC_NAMEList []GENERIC_TYPE

		type GENERIC_NAMEToBool func(*GENERIC_TYPE) bool

		func (al GENERIC_NAMEList) Filter(f GENERIC_NAMEToBool) GENERIC_NAMEList {
			var ret GENERIC_NAMEList
			for _, a := range al {
				if f(&a) {
					ret = append(ret, a)
				}
			}
			return ret
		}
	可以在需要使用这个的地方，加上相关的 Go Generate 的注释
		代码如下

第三方工具
	可以直接使用第三方已经写好的 gen.sh 工具
	Genny：https://github.com/cheekybits/genny
	Generic：https://github.com/taylorchu/generic
	GenGen：https://github.com/joeshaw/gengen
	Gen：https://github.com/clipperhouse/gen
*/

//type Employee struct {
//	Name     string
//	Age      int
//	Vacation int
//	Salary   int
//}
//
////go:generate ./gen.sh ./template/filter.tmp.go gen Employee filter
//func filterEmployeeExample() {
//
//	var list = EmployeeList{
//		{"Hao", 44, 0, 8000},
//		{"Bob", 34, 10, 5000},
//		{"Alice", 23, 5, 9000},
//		{"Jack", 26, 0, 4000},
//		{"Tom", 48, 9, 7500},
//	}
//
//	var filter EmployeeList
//	filter = list.Filter(func(e *Employee) bool {
//		return e.Age > 40
//	})
//
//	fmt.Println("----- Employee.Age > 40 ------")
//	for _, e := range filter {
//		fmt.Println(e)
//	}
//
//	filter = list.Filter(func(e *Employee) bool {
//		return e.Salary <= 5000
//	})
//
//	fmt.Println("----- Employee.Salary <= 5000 ------")
//	for _, e := range filter {
//		fmt.Println(e)
//	}
//}
