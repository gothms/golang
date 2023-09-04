package _1_basic

/*
Go的错误机制：与其他主要编程语言的差异
	1.没有异常机制
	2.error类型实现了error接口
	3.可以通过 errors.New 来快速创建错误实例

最佳实践：
	及早失败，避免嵌套

panic
	用于不可以恢复的错误
	退出前会执行defer指定的内容
os.Exit
	退出时不会调用defer指定的函数
	退出时不输出当前调用栈信息
recover：当心 recover 成为恶魔
	形成僵尸服务进程，导致 health check 失效
	"Let it Crash!" 往往是我们恢复不确定性错误的最好方法
*/
