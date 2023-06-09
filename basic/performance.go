package main

/*
1.数字转字符串，使用 strconv.Itoa() 比 fmt.Sprintf 快一倍左右
	2.尽可能避免把 string 转成 []Byte，会导致性能下降
	3.在 for-loop 里对 Slice 使用 append()，请先把 Slice 容量扩充到位
		避免系统自动按 2^n 进行扩展，但又用不到的情况，从而避免浪费内存
	4.使用 StringBuffer 或 StringBuild 拼接字符串，性能比使用 + 或 += 高三到四个数量级
	5.尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作
	6.避免在热代码中进行内存分配，会导致 gc 很忙。尽可能使用 sync.Pool 来重用对象
	7.使用 lock-free 的操作，避免使用 mutex，尽可能使用 sync/Atomic 包。关于无锁编程，参考：
		无锁队列实现：https://coolshell.cn/articles/8239.html
		无锁Hashmap实现：https://coolshell.cn/articles/9703.html
	8.关于 I/O 缓冲，I/O 是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可带来更高的性能
	9.在 for-loop 里的固定的正则表达式，一定使用 regexp.Compile() 编译正则表达式，性能会提升两个数量级
	10.如果需要更高性能的协议，就考虑使用 protobuf 或 msgp，而不是JSON，因为JSON的序列化和反序列化里使用了反射
		protobuf：https://github.com/golang/protobuf
		msgp：https://github.com/tinylib/msgp
	11.使用 map 时，使用整型 key 会比字符串快，因为整型比较比字符串比较快
	12.更多技巧，写出更好的 Go，必读：
		Effective Go：https://golang.org/doc/effective_go.html
		Uber Go Style：https://github.com/uber-go/guide/blob/master/style.md
		50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs：http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
		Go Advice：https://github.com/cristaloleg/go-advice
		Practical Go Benchmarks：https://www.instana.com/blog/practical-golang-benchmarks/
		Benchmarks of Go serialization methods：https://github.com/alecthomas/go_serialization_benchmarks
		Debugging performance issues in Go programs：https://github.com/golang/go/wiki/Performance
		Go code refactoring: the 23x performance hunt：https://medium.com/@val_deleplace/go-code-refactoring-the-23x-performance-hunt-156746b522f7
*/
