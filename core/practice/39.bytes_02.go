package practice

import (
	"bytes"
	"fmt"
)

/*
bytes包与字节串操作（下）

知识扩展
问题 1：bytes.Buffer的扩容策略是怎样的？
	Buffer值既可以被手动扩容，也可以进行自动扩容
		这两种扩容方式的策略是基本一致的。所以，除非我们完全确定后续内容所需的字节数，否则让Buffer值自动去扩容就好了
		在扩容的时候，Buffer值中相应的代码会先判断内容容器的剩余容量，是否可以满足调用方的要求，或者是否足够容纳新的内容
		如果可以，那么扩容代码会在当前的内容容器之上，进行长度扩充
			b.buf = b.buf[:length+need]
		反之，如果内容容器的剩余容量不够了，那么扩容代码可能就会用新的内容容器去替代原有的内容容器，从而实现扩容
			扩容策略如下
	扩容策略
		如果当前内容容器的容量的一半仍然大于或等于其现有长度再加上另需的字节数的和，即：
			cap(b.buf)/2 >= len(b.buf)+need
			那么，扩容代码就会复用现有的内容容器，并把容器中的未读内容拷贝到它的头部位置
			这也意味着其中的已读内容，将会全部被未读内容和之后的新内容覆盖掉
		如果当前内容容器的容量小于新长度的二倍，那么扩容代码就只能再创建一个新的内容容器，并把原有容器中的未读内容拷贝进去，最后再用新的容器替换掉原有的容器
			这个新容器的容量将会等于原有容量的二倍再加上另需字节数的和
			新容器的容量 =2* 原有容量 + 所需字节数
		扩容后还会把已读计数置为0，并再对内容容器做一下切片操作，以掩盖掉原有的已读内容
	对于处在零值状态的Buffer值
		如果第一次扩容时的另需字节数不大于64，那么该值就会基于一个预先定义好的、长度为64的字节数组来创建内容容器
问题 2：bytes.Buffer中的哪些方法可能会造成内容的泄露？
	内容泄漏
		使用Buffer值的一方通过某种非标准的方式得到了本不该得到的内容
		比如说，我通过调用Buffer值的某个用于读取内容的方法，得到了一部分未读内容
		但是，在这个Buffer值又有了一些新内容之后，我却可以通过当时得到的结果值，直接获得新的内容，而不需要再次调用相应的方法
		这种读取方式是不应该存在的，即使存在，我们也不应该使用。因为它是在无意中暴露出来的，其行为很可能是不稳定的
	典型的非标准读取方式
		在bytes.Buffer中，Bytes方法和Next方法都可能会造成内容的泄露
		原因在于，它们都把基于内容容器的切片直接返回给了方法的调用方
		通过切片，我们可以直接访问和操纵它的底层数组。不论这个切片是基于某个数组得来的，还是通过对另一个切片做切片操作获得的
		Bytes方法和Next方法返回的字节切片，都是通过对内容容器做切片操作得到的。也就是说，它们与内容容器共用了同一个底层数组，起码在一段时期之内是这样的
	示例：BufferContentLeak()
		为什么该值的容量却变为了8
		源码分析：runtime包中 string.stringtoslicebyte 函数
			stringtoslicebyte
			rawbyteslice
			roundupsize
			return uintptr(class_to_size[size_to_class8[divRoundUp(size, smallSizeDiv)]])
		向该值写入了字符串值"cdefg"，此时，其容量仍然是8
			unreadBytes与buffer1的内容容器在此时还共用着同一个底层数组
			所以，只需通过简单的再切片操作，就可以利用这个结果值拿到buffer1在此时的所有未读内容
			如此一来，buffer1的新内容就被泄露出来了
		如果把unreadBytes的值传到了外界，那么外界就可以通过该值操纵buffer1的内容了
	示例：demo80.go
		对于Buffer值的Next方法，也存在相同的问题
		如果经过扩容，Buffer值的内容容器或者它的底层数组被重新设定了，那么之前的内容泄露问题就无法再进一步发展了

总结
	Buffer值的扩容方法并不一定会为了获得更大的容量，替换掉现有的内容容器
		而是先会本着尽量减少内存分配和内容拷贝的原则，对当前的内容容器进行重用
		并且，只有在容量实在无法满足要求的时候，它才会去创建新的内容容器
	Buffer值的某些方法可能会造成内容的泄露
		这主要是由于这些方法返回的结果值，在一段时期内会与其所属值的内容容器共用同一个底层数组
		如果我们有意或无意地把这些结果值传到了外界，那么外界就有可能通过它们操纵相关联Buffer值的内容
		这属于很严重的数据安全问题
		最彻底的做法是，在传出切片这类值之前要做好隔离。比如，先对它们进行深度拷贝，然后再把副本传出去

思考
	对比strings.Builder和bytes.Buffer的String方法，并判断哪一个更高效？原因是什么？
*/

func BufferContentLeak() {
	contents := "ab"
	buffer1 := bytes.NewBufferString(contents)
	fmt.Printf("The capacity of new buffer with contents %q: %d\n",
		contents, buffer1.Cap()) // 内容容器的容量为：8
	unreadBytes := buffer1.Bytes()
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // 未读内容为：[97 98]

	buffer1.WriteString("cdefg")
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 内容容器的容量仍为：8
	unreadBytes = unreadBytes[:cap(unreadBytes)]
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes) // 基于前面获取到的结果值可得，未读内容为：[97 98 99 100 101 102 103 0]

	unreadBytes[len(unreadBytes)-2] = byte('X')                         // 'X'的 ASCII 编码为 88
	fmt.Printf("The unread bytes of the buffer: %v\n", buffer1.Bytes()) // 未读内容变为了：[97 98 99 100 101 102 88]
}
