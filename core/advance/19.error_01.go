package advance

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

/*
错误处理（上）

error 接口
	error类型其实是一个接口类型，也是一个 Go 语言的内建类型
	在这个接口类型的声明中只包含了一个方法Error。这个方法不接受任何参数，但是会返回一个string类型的结果
	它的作用是返回错误信息的字符串表示形式
通常使用方式
	使用error类型的方式通常是，在函数声明的结果列表的最后，声明一个该类型的结果
	同时在调用这个函数之后，先判断它返回的最后一个结果值是否“不为nil”
	如果这个值“不为nil”，那么就进入错误处理流程，否则就继续进行正常的流程
卫述语句
	12.func.go
	被用来检查后续操作的前置条件并进行相应处理的语句
	在进行错误处理的时候经常会用到卫述语句，以至于程序满屏都是卫述语句，简直是太难看了！

errors.New函数
	一种最基本的生成错误值的方式
	我们调用它的时候传入一个由字符串代表的错误信息，它会给返回给我们一个包含了这个错误信息的error类型值
error & *errorString
	error类型值的静态类型当然是error，而动态类型则是一个在errors包中的，包级私有的类型*errorString
	显然，errorString类型拥有的一个指针方法实现了error接口中的Error方法
	这个方法在被调用后，会原封不动地返回我们之前传入的错误信息
	实际上，error类型值的Error方法就相当于其他类型值的String方法
Printf & Error
	通过调用fmt.Printf函数，并给定占位符%s就可以打印出某个值的字符串表示形式
	对于其他类型的值来说，只要我们能为这个类型编写一个String方法，就可以自定义它的字符串表示形式
	而对于error类型值，它的字符串表示形式则取决于它的Error方法
	fmt.Printf函数如果发现被打印的值是一个error类型的值，那么就会去调用它的Error方法。fmt包中的这类打印函数其实都是这么做的
fmt.Errorf
	当我们想通过模板化的方式生成错误信息，并得到错误值时，可以使用fmt.Errorf函数
	该函数所做的其实就是先调用fmt.Sprintf函数，得到确切的错误信息
	再调用errors.New函数，得到包含该错误信息的error类型值，最后返回该值

问题：对于具体错误的判断，Go 语言中都有哪些惯用法？
	即：怎样判断一个错误值具体代表的是哪一类错误？
	由于error是一个接口类型，所以即使同为error类型的错误值，它们的实际类型也可能不同
典型回答
	1. 对于类型在已知范围内的一系列错误值，一般使用类型断言表达式或类型switch语句来判断
	2. 对于已有相应变量且类型相同的一系列错误值，一般直接使用判等操作来判断
	3. 对于没有相应变量且类型未知的一系列错误值，只能使用其错误信息的字符串表示形式来做判断
问题解析
	类型不同：switch x.(type)
		如os包中的几个代表错误的类型os.PathError、os.LinkError、os.SyscallError和os/exec.Error
		它们的指针类型都是error接口的实现类型，同时它们也都包含了一个名叫Err，类型为error接口类型的代表潜在错误的字段
		如果我们得到一个error类型值，并且知道该值的实际类型肯定是它们中的某一个，那么就可以用类型switch语句去做判断
	类型相同：switch error
		在 Go 语言的标准库中也有不少以相同方式创建的同类型的错误值
		os包中不少的错误值都是通过调用errors.New函数来初始化的，比如：os.ErrClosed、os.ErrInvalid以及os.ErrPermission 等
			与前面讲到的那些错误类型不同，这几个都是已经定义好的、确切的错误值
			os包中的代码有时候会把它们当做潜在错误值，封装进前面那些错误类型的值中
		如果我们在操作文件系统的时候得到了一个错误值，并且知道该值的潜在错误值肯定是上述值中的某一个，那么就可以用普通的switch语句去做判断
	未知类型
		只能通过它拥有的错误信息去做判断
		好在我们总是能通过错误值的Error方法，拿到它的错误信息
		其实os包中就有做这种判断的函数，比如：os.IsExist、os.IsNotExist和os.IsPermission

思考
	请列举出你经常用到或者看到的 3 个错误类型，它们所在的错误类型体系都是怎样的？你能画出一棵树来描述它们吗？
*/

// underlyingError 会返回已知的操作系统相关错误的潜在错误值
func underlyingError(err error) error {
	switch err := err.(type) {
	case *os.PathError:
		return err.Err
	case *os.LinkError:
		return err.Err
	case *os.SyscallError:
		return err.Err
	case *exec.Error:
		return err.Err
	}
	return err
}

// switchErr 类型相同
func switchErr(i int, err error) {
	err = underlyingError(err)
	switch err {
	case os.ErrClosed:
		fmt.Printf("error(closed)[%d]: %s\n", i, err)
	case os.ErrInvalid:
		fmt.Printf("error(invalid)[%d]: %s\n", i, err)
	case os.ErrPermission:
		fmt.Printf("error(permission)[%d]: %s\n", i, err)
	}
}

func SwitchTypeError() {
	// 示例1
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Printf("unexpected error: %s\n", err)
		return
	}
	// 人为制造 *os.PathError 类型的错误。
	r.Close()
	_, err = w.Write([]byte("hi"))
	uError := underlyingError(err)
	fmt.Printf("underlying error: %s (type: %T)\n",
		uError, uError)
	fmt.Println()

	// 示例2
	paths := []string{
		os.Args[0],           // 当前的源码文件或可执行文件。
		"/it/must/not/exist", // 肯定不存在的目录。
		os.DevNull,           // 肯定存在的目录。
	}
	printError := func(i int, err error) {
		if err == nil {
			fmt.Println("nil error")
			return
		}
		err = underlyingError(err)
		switch err {
		case os.ErrClosed:
			fmt.Printf("error(closed)[%d]: %s\n", i, err)
		case os.ErrInvalid:
			fmt.Printf("error(invalid)[%d]: %s\n", i, err)
		case os.ErrPermission:
			fmt.Printf("error(permission)[%d]: %s\n", i, err)
		}
	}
	var f *os.File
	var index int
	{
		index = 0
		f, err = os.Open(paths[index])
		if err != nil {
			fmt.Printf("unexpected error: %s\n", err)
			return
		}
		// 人为制造潜在错误为 os.ErrClosed 的错误。
		f.Close()
		_, err = f.Read([]byte{})
		printError(index, err)
	}
	{
		index = 1
		// 人为制造 os.ErrInvalid 错误。
		f, _ = os.Open(paths[index])
		_, err = f.Stat()
		printError(index, err)
	}
	{
		index = 2
		// 人为制造潜在错误为 os.ErrPermission 的错误。
		_, err = exec.LookPath(paths[index])
		printError(index, err)
	}
	if f != nil {
		f.Close()
	}
	fmt.Println()

	// 示例3
	paths2 := []string{
		runtime.GOROOT(),     // 当前环境下的Go语言根目录。
		"/it/must/not/exist", // 肯定不存在的目录。
		os.DevNull,           // 肯定存在的目录。
	}
	printError2 := func(i int, err error) {
		if err == nil {
			fmt.Println("nil error")
			return
		}
		err = underlyingError(err)
		if os.IsExist(err) {
			fmt.Printf("error(exist)[%d]: %s\n", i, err)
		} else if os.IsNotExist(err) {
			fmt.Printf("error(not exist)[%d]: %s\n", i, err)
		} else if os.IsPermission(err) {
			fmt.Printf("error(permission)[%d]: %s\n", i, err)
		} else {
			fmt.Printf("error(other)[%d]: %s\n", i, err)
		}
	}
	{
		index = 0
		err = os.Mkdir(paths2[index], 0700)
		printError2(index, err)
	}
	{
		index = 1
		f, err = os.Open(paths[index])
		printError2(index, err)
	}
	{
		index = 2
		_, err = exec.LookPath(paths[index])
		printError2(index, err)
	}
	if f != nil {
		f.Close()
	}
}
