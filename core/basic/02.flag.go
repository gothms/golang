package basic

import (
	"flag"
	"fmt"
	"os"
)

/*
命令源码文件

三种源码文件
	环境变量 GOPATH 指向的是一个或多个工作区，每个工作区中都会有以代码包为基本组织形式的源码文件
	这里的源码文件又分为三种，即：命令源码文件、库源码文件和测试源码文件，它们都有着不同的用途和编写规则

问题
	命令源码文件的用途是什么，怎样编写它？
	A：
		命令源码文件是程序的运行入口，是每个可独立运行的程序必须拥有的
		我们可以通过构建或安装，生成与其对应的可执行文件，后者一般会与该命令源码文件的直接父目录同名
		如果一个源码文件声明属于main包，并且包含一个无参数声明且无结果声明的main函数，那么它就是命令源码文件
		go run main.go
	补充
		当需要模块化编程时，我们往往会将代码拆分到多个文件，甚至拆分到不同的代码包中
		但无论怎样，对于一个独立的程序来说，命令源码文件永远只会也只能有一个
		如果有与命令源码文件同包的源码文件，那么它们也应该声明属于main包
	接收参数
		无论是 Linux 还是 Windows，几乎所有命令（command）都是可以接收参数（argument）的
		通过构建或安装命令源码文件，生成的可执行文件就可以被视为“命令”，既然是命令，那么就应该具备接收参数的能力

知识精讲
1. 命令源码文件怎样接收参数
	flag
		Go 语言标准库中有一个代码包专门用于接收和解析命令参数
		这个代码包的名字叫flag
	flag.StringVar(&name, "name", "everyone", "The greeting object.")
		第 1 个参数是用于存储该命令参数的值的地址，具体到这里就是声明的变量name的地址了，由表达式&name表示
			var name string
		第 2 个参数是为了指定该命令参数的名称，这里是name
		第 3 个参数是为了指定在未追加该命令参数时的默认值，这里是everyone
		第 4 个函数参数，即是该命令参数的简短说明了，这在打印命令说明时会用到
	flag.String
		区别是，flag.String 会直接返回一个已经分配好的用于存储命令参数值的地址
		var name = *flag.String("name", "everyone", "The greeting object.")
	flag.Parse()
		用于真正解析命令参数，并把它们的值赋给相应的变量
		对该函数的调用必须在所有命令参数存储载体的声明（这里是对变量name的声明）和设置（这里是对flag.StringVar函数的调用）之后，并且在读取任何命令参数值之前进行
		正因为如此，我们最好把flag.Parse()放在main函数的函数体的第一行】
	注意：是不同的
		flag.Parsed()
		flag.Parse()
2. 怎样在运行命令源码文件的时候传入参数，又怎样查看参数的使用说明
	命令
		go run main.go -name="lee"
		go test -v 02_test.go -name="lee"
		go test -v -run TestFlag 02_test.go -name="lee"
	查看该命令源码文件的参数说明
		go run main.go --help
			Usage of C:\Users\sc\AppData\Local\Temp\go-build3180127776\b001\exe\main.exe:
			  -name string
					The greeting object. (default "everyone")
		Usage of C:\...
			go run命令构建上述命令源码文件时临时生成的可执行文件的完整路径
		go run main.go --help 等同于：
			go build main.go
			./main --help
			但是输出不同
				Usage of E:\gothmslee\golang\main\main.exe:
				  -name string
						The greeting object. (default "everyone")
3. 怎样自定义命令源码文件的参数使用说明
	flag.Usage
		flag.Usage的类型是func()，即一种无参数声明且无结果声明的函数类型
		flag.Usage变量在声明时就已经被赋值了，所以我们才能够在运行命令go run main.go --help时看到正确的结果
		对flag.Usage的赋值必须在调用flag.Parse函数之前
	测试一
		Usage -> of question:
	flag.CommandLine
		在调用flag包中的一些函数（比如StringVar、Parse等等）的时候，实际上是在调用flag.CommandLine变量的对应方法
		flag.CommandLine相当于默认情况下的命令参数容器
		通过对flag.CommandLine重新赋值，我们可以更深层次地定制当前命令源码文件的参数使用说明
	测试二
		注意：
			此时 flag.StringVar 的调用要在 flag.CommandLine 之后
		flag.PanicOnError
			修改 panic 错误说明
		Usage --> of question:
	flag.PanicOnError和flag.ExitOnError都是预定义在flag包中的常量
		flag.ExitOnError的含义是，告诉命令参数容器，当命令后跟--help或者参数设置的不正确的时候，在打印命令参数使用说明后以状态码2结束当前程序
		状态码2代表用户错误地使用了命令，而flag.PanicOnError与之的区别是在最后抛出“运行时恐慌（panic）”
		两种情况都会在我们调用flag.Parse函数时被触发
	测试三
		索性不用全局的flag.CommandLine变量，转而自己创建一个私有的命令参数容器
		把对flag.StringVar的调用替换为对cmdLine.StringVar调用，再把flag.Parse()替换为cmdLine.Parse(os.Args[1:])
		其中的os.Args[1:]指的就是我们给定的那些命令参数
		这样做就完全脱离了flag.CommandLine
		注意：
			此时就可以注释 flag.Parse()
		Usage of question:
	好处
		*flag.FlagSet类型的变量cmdLine拥有很多有意思的方法
		更灵活地定制命令参数容器
		更重要的是，你的定制完全不会影响到那个全局变量flag.CommandLine

总结
	用 Go 编写命令，可以让它们像众多操作系统命令那样被使用，甚至可以把它们嵌入到各种脚本中
	flag包的用法
		https://golang.google.cn/pkg/flag/
	使用godoc命令在本地启动一个 Go 语言文档服务器
		https://github.com/hyper0x/go_command_tutorial/blob/master/0.5.md

思考
	1.默认情况下，我们可以让命令源码文件接受哪些类型的参数值？
		查阅文档获得答案
		string, bool, int
		中文文档：https://studygolang.com/pkgdoc
	2.我们可以把自定义的数据类型作为参数值的类型吗？如果可以，怎样做？
		可以自定义一个用于flag的类型（满足Value接口）并将该类型用于flag解析
		flag.Var(&flagVal, "name", "help message for flagname")
*/

var name string

//var id int

// 测试三
var cmdLine = flag.NewFlagSet("question", flag.ExitOnError)

func init() {
	// 测试二
	//flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	//flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	//flag.CommandLine.Usage = func() {
	//	fmt.Fprintf(os.Stderr, "Usage --> of %s:\n", "question")
	//	flag.PrintDefaults()
	//}
	//flag.StringVar(&name, "name", "everyone", "The greeting object.")
}
func init() {
	//flag.StringVar(&name, "name", "everyone", "The greeting object.")
	// 测试三
	cmdLine.StringVar(&name, "name", "everyone", "The greeting object.")
	//cmdLine.IntVar(&id, "id", 0, "The greeting id.")
}

// TestFlag go run main.go -name="lee"
// go test -v 02_test.go -name="lee"
// go test -v -run TestFlag 02_test.go -name="lee"
func TestFlag() {
	//flag.Parsed()
	flag.Parse()
	fmt.Printf("Hello, %s!\n", name)
}

// TestFlagUsage 需要在 main 中测试，在 test 中测试无效
func TestFlagUsage() {
	// 测试一
	//flag.Usage = func() {
	//	fmt.Fprintf(os.Stderr, "Usage -> of %s:\n", "question")
	//	flag.PrintDefaults()
	//}

	// 测试三
	cmdLine.Parse(os.Args[1:])

	//flag.Parse()
	fmt.Printf("Hello, %s!\n", name)
	// go run main.go -name="abc" -id=3
	//fmt.Printf("Hello, %s! %d\n", name, id)
}
