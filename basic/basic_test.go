package basic

import "testing"

/*
main：
	1.os.Exit(-1)：退出时，返回 -1
	2.os.Args
		func main() {
			// go run main.go lee
			if len(os.Args) > 1 { // hello,world! lee
				fmt.Println("hello,world!", os.Args[1])
			}
			os.Exit(-1) // exit status 4294967295
		}

1.不支持类型隐式转换，也不支持别名到原类型的隐式转换
2.指针不支持任何运算
3.string是数据类型，不是引用/指针类型，默认值是空字符串，而不是 nil
*/
func TestConst(t *testing.T) {
	t.Log(Day8)
	admin := 7
	t.Log(admin&Readable,
		admin&Writable,
		admin&Executable)
}

const (
	Monday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	//Saturday
	Sunday = iota + 2
	Day8
)
const (
	Readable = 1 << iota
	Writable
	Executable
)
