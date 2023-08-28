package practice

import (
	"fmt"
	"strings"
	"unsafe"
)

/*
strings包与字符串操作

问题：与string值相比，strings.Builder类型的值有哪些优势？
	strings.Builder类型的值的优势有下面的三种：
		已存在的内容不可变，但可以拼接更多的内容
		减少了内存分配和内容拷贝的次数
		可将内容重置，可重用值
问题解析
	string 底层
		在 Go 语言中，string类型的值是不可变的
		只能基于原字符串进行裁剪、拼接等操作，从而生成一个新的字符串。裁剪操作可以使用切片表达式，而拼接操作可以用操作符+实现
		在底层，一个string值的内容会被存储到一块连续的内存空间中。同时，这块内存容纳的字节数量也会被记录下来，并用于表示该string值的长度
		相应的string值则包含了指向一个字节数组头部的指针值
		在一个string值上应用切片表达式，就相当于在对其底层的字节数组做切片
	拼接
		Go 语言会把所有被拼接的字符串依次拷贝到一个崭新且足够大的连续内存空间中，并把持有相应指针值的string值作为结果返回
		当程序中存在过多的字符串拼接操作的时候，会对内存的分配产生非常大的压力
	值类型
		虽然string值在内部持有一个指针值，但其类型仍然属于值类型
		不过，由于string值的不可变，其中的指针值也为内存空间的节省做出了贡献
		一个string值会在底层与它的所有副本共用同一个字节数组。由于这里的字节数组永远不会被改变，所以这样做是绝对安全的
	Builder 底层
		与string值相比，Builder值的优势其实主要体现在字符串拼接方面
		Builder值中有一个用于承载内容的容器，它是一个以byte为元素类型的切片
		由于这样的字节切片的底层数组就是一个字节数组，所以它与string值存储内容的方式是一样的
		实际上，它们都是通过一个unsafe.Pointer类型的字段来持有那个指向了底层字节数组的指针值的
			正是因为这样的内部构造，Builder值同样拥有高效利用内存的前提条件
		虽然，对于字节切片本身来说，它包含的任何元素值都可以被修改，但是Builder值并不允许这样做，其中的内容只能够被拼接或者完全重置
			即已存在于Builder值中的内容是不可变的
	API：一系列指针方法（拼接方法）
		Write、WriteByte、WriteRune和WriteString
	自动扩容
		有必要，Builder值会自动地对自身的内容容器进行扩容，自动扩容策略与切片的扩容策略一致
		只要内容容器的容量够用，扩容就不会进行，针对于此的内存分配也不会发生。同时，只要没有扩容，Builder值中已存在的内容就不会再被拷贝
	Grow：手动扩容
		Grow方法也可以被称为扩容方法，它接受一个int类型的参数n，该参数用于代表将要扩充的字节数量
		Grow方法会把其所属值中内容容器的容量增加n个字节
		实际上，它会生成一个字节切片作为新的内容容器，该切片的容量会是原容器容量的二倍再加上n
		之后，它会把原容器中的所有字节全部拷贝到新容器中
		当然，Grow方法还可能什么都不做。前提条件是：当前的内容容器中的未用容量已经够用了，即：未用容量大于或等于n
	Reset：重用 Builder 值
		通过调用它的Reset方法，我们可以让Builder值重新回到零值状态，就像它从未被使用过那样
		一旦被重用，Builder值中原有的内容容器会被直接丢弃
		之后，它和其中的所有内容，将会被 Go 语言的垃圾回收器标记并回收掉

知识扩展
问题 1：strings.Builder类型在使用上有约束吗？
	约束：在已被真正使用后就不可再被复制（拼接或扩容，会改变其内容容器的状态）
		由于其内容不是完全不可变的，所以需要使用方自行解决操作冲突和并发安全问题
		否则，拷贝原变量后，只要在任何副本上调用上述方法就都会引发 panic
		复制方式，包括但不限于在函数间传递值、通过通道传递值、把值赋予变量等
	约束的好处
		正是由于已使用的Builder值不能再被复制，所以肯定不会出现多个Builder值中的内容容器共用一个底层字节数组的情况
		这样也就避免了多个同源的Builder值在拼接内容时可能产生的冲突问题
	复制指针值
		无论什么时候，我们都可以通过任何方式复制这样的指针值，这样的指针值指向的都会是同一个Builder值
		但是会产生问题：操作冲突和并发安全
	操作冲突和并发安全问题
		如果Builder值被多方同时操作，那么其中的内容就很可能会产生混乱
		所以，我们在通过传递其指针值共享Builder值的时候，一定要确保各方对它的使用是正确、有序的，并且是并发安全的
		而最彻底的解决方案是，绝不共享Builder值以及它的指针值
		即分别声明、分开使用、互不干涉
	不过，对于处在零值状态的Builder值，复制不会有任何问题
		第一次使用时，才分配底层切片
问题 2：为什么说strings.Reader类型的值可以高效地读取字符串？
	strings.Reader类型是为了高效读取字符串而存在的，高效主要体现在它对字符串的读取机制上
	读取机制
		它封装了很多用于在string值上读取内容的最佳实践
		在读取的过程中，Reader值会保存已读取的字节的计数
			已读计数也代表着下一次读取的起始索引位置
			Reader值正是依靠这样一个计数，以及针对字符串值的切片表达式，从而实现快速读取
			此外，这个已读计数也是读取回退和位置设定时的重要依据
		计算已读计数
			虽然它属于Reader值的内部结构，但我们还是可以通过该值的Len方法和Size把它计算出来的
			readingIndex := reader.Size() - int64(reader.Len())
	API：demo77.go
		Reader值拥有的大部分用于读取的方法都会及时地更新已读计数
			比如，ReadByte方法会在读取成功后将这个计数的值加1
			又比如，ReadRune方法在读取成功之后，会把被读取的字符所占用的字节数作为计数的增量
		ReadAt
			它既不会依据已读计数进行读取，也不会在读取后更新它
			正因为如此，这个方法可以自由地读取其所属的Reader值中的任何内容
		Seek
			也会更新该值的已读计数
			实际上，这个Seek方法的主要作用正是设定下一次读取的起始索引位置
			把常量io.SeekCurrent的值作为第二个参数值传给该方法，那么它还会依据当前的已读计数，以及第一个参数offset的值来计算新的计数值

strings包其他常用函数
	`Count`、`IndexRune`、`Map`、`Replace`、`SplitN`、`Trim`，等

思考
	*strings.Builder和*strings.Reader都分别实现了哪些接口？这样做有什么好处吗？
A
	strings.Builder类型实现了 3 个接口，分别是：fmt.Stringer、io.Writer和io.ByteWriter
	strings.Reader类型则实现了 8 个接口，即：io.Reader、io.ReaderAt、io.ByteReader、io.RuneReader、io.Seeker、io.ByteScanner、io.RuneScanner和io.WriterTo
	实现的接口越多，它们的用途就越广。它们会适用于那些要求参数的类型为这些接口类型的地方
*/

func StringsBuilder() {
	//strings.Builder{}
	s := "abc"
	pointer := (*[]byte)(unsafe.Pointer(&s))
	fmt.Println(*pointer)
	//fmt.Printf("%p\n", &(*pointer)[1])
	//(*pointer)[1] = 101
	//fmt.Println(s)
	//a := (*string)(unsafe.Pointer(pointer))
	//fmt.Println(*a)

	var sb strings.Builder
	sb.Grow(20)
	fmt.Println(sb.Len(), sb.Cap())
	sb.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9})
	sb.Grow(5)  // 仍然是 20
	sb.Grow(13) // 扩容至 53
	fmt.Println(sb.Len(), sb.Cap())
	// 这样的使用方式是并不合法的，因为这里的Builder值是副本而不是原值
	//sbCopy := sb
	//sbCopy.WriteByte('a') // panic: strings: illegal use of non-zero Builder copied by value

	var ssb strings.Builder
	s1, s2 := ssb, ssb // 三个变量有不同的值
	s1.WriteByte('a')  // 第一次使用才分配内容容器
	s2.WriteByte('b')
	ssb.WriteByte('c')
	fmt.Println(s1, s2, ssb)
}
