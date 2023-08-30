package practice

import (
	"bytes"
	"fmt"
)

/*
bytes包与字节串操作（上）

strings vs bytes
	strings包和bytes包可以说是一对孪生兄弟，它们在 API 方面非常的相似
	单从它们提供的函数的数量和功能上讲，差别微乎其微
	strings包主要面向的是 Unicode 字符和经过 UTF-8 编码的字符串，而bytes包面对的则主要是字节和字节切片
strings.Builder vs bytes.Buffer
	bytes.Buffer类型的用途主要是作为字节序列的缓冲区
	与strings.Builder类型一样，bytes.Buffer也是开箱即用的
	但不同的是，strings.Builder只能拼接和导出字符串
	而bytes.Buffer不但可以拼接、截断其中的字节序列，以各种形式导出其中的内容，还可以顺序地读取其中的子序列
	作为一个缓冲区，bytes.Buffer是集读、写功能于一身的数据类型
bytes.Buffer 底层
	bytes.Buffer类型同样是使用字节切片作为内容容器的
	并且，与strings.Reader类型类似，bytes.Buffer有一个int类型的字段，用于代表已读字节的计数，可以简称为已读计数
	Buffer 的已读计数无法通过bytes.Buffer提供的方法计算出来
	与strings.Reader类型的Len方法一样，buffer1的Len方法返回的也是内容容器中未被读取部分的长度，而不是其中已存内容的总长度
示例：Len
	BufferSlice()
	Buffer值的长度是未读内容的长度，而不是已存内容的总长度
	它与在当前值之上的读操作和写操作都有关系，并会随着这两种操作的进行而改变，它可能会变得更小，也可能会变得更大
	而Buffer值的容量指的是它的内容容器的容量，它只与在当前值之上的写操作有关，并会随着内容的写入而不断增长
已读计数
	由于strings.Reader还有一个Size方法可以给出内容长度的值，所以我们用内容长度减去未读部分的长度，就可以很方便地得到它的已读计数
	然而，bytes.Buffer类型却没有这样一个方法，它只有Cap方法。可是Cap方法提供的是内容容器的容量，也不是内容长度
	因此，没有了现成的计算公式，只要遇到稍微复杂些的情况，我们就很难估算出Buffer值的已读计数

问题：bytes.Buffer类型的值记录的已读计数，在其中起到了怎样的作用？
	bytes.Buffer中的已读计数的大致功用如下
	1. 读取内容时，相应方法会依据已读计数找到未读部分，并在读取后更新计数
	2. 写入内容时，如需扩容，相应方法会根据已读计数实现扩容策略
	3. 截断内容时，相应方法截掉的是已读计数代表索引之后的未读部分
	4. 读回退时，相应方法需要用已读计数记录回退点
	5. 重置内容时，相应方法会把已读计数置为0
	6. 导出内容时，相应方法只会导出已读计数代表的索引之后的未读部分
	7. 获取长度时，相应方法会依据已读计数和内容容器的长度，计算未读部分的长度并返回
问题解析
	读取内容：包括了所有名称以Read开头的方法，以及Next方法和WriteTo方法
		相应方法会先根据已读计数，判断一下内容容器中是否还有未读的内容。如果有，那么它就会从已读计数代表的索引处开始读取
		在读取完成后，它还会及时地更新已读计数
		相应方法
	写入内容：包括了所有名称以Write开头的方法，以及ReadFrom方法
		绝大多数的相应方法都会先检查当前的内容容器，是否有足够的容量容纳新的内容。如果没有，那么它们就会对内容容器进行扩容
		在扩容的时候，在必要时会依据已读计数找到未读部分，并把其中的内容拷贝到扩容后内容容器的头部位置
		然后，方法将会把已读计数的值置为0，以表示下一次读取需要从内容容器的第一个字节开始
	截断内容：方法Truncate
		它会接受一个int类型的参数，这个参数的值代表了：在截断时需要保留头部（未读部分的头部）的多少个字节
		头部的起始索引正是由已读计数的值表示的。因此，已读计数的值再加上参数值后得到的和，就是内容容器新的总长度
	读回退：UnreadByte和UnreadRune方法
		UnreadByte和UnreadRune方法，分别用于回退一个字节和回退一个 Unicode 字符
		调用它们一般都是为了退回在上一次被读取内容末尾的那个分隔符，或者为重新读取前一个字节或字符做准备
		退回的前提是，在调用它们之前的那一个操作必须是“读取”，并且是成功的读取，否则这些方法就只能忽略后续操作并返回一个非nil的错误值
		UnreadByte方法，把已读计数的值减1
		UnreadRune方法需要从已读计数中减去上一次被读取的 Unicode 字符所占用的字节数
			由bytes.Buffer的 lastRead readOp 字段负责存储
			它在这里的有效取值范围是 [1, 4]，只有ReadRune方法才会把这个字段的值设定在此范围之内
		即，只有紧接在调用ReadRune方法之后，对UnreadRune方法的调用才能够成功完成
	Bytes方法和String方法
		bytes.Buffer的Len方法返回的是内容容器中未读部分的长度，而不是其中已存内容的总长度
		Bytes方法和String方法与Len方法是保持一致的，即只会去访问未读部分中的内容，并返回相应的结果值
	小结
		在已读计数代表的索引之前的那些内容，永远都是已经被读过的，它们几乎没有机会再次被读取
		不过，这些已读内容所在的内存空间可能会被存入新的内容。这一般都是由于重置或者扩充内容容器导致的
		这时，已读计数一定会被置为0，从而再次指向内容容器中的第一个字节

总结
	虽然我们无法直接计算出这个已读计数，但是由于它在Buffer值中起到的作用非常关键，所以我们很有必要去理解它
	无论是读取、写入、截断、导出还是重置，已读计数都是功能实现中的重要一环
*/

func BufferSlice() {
	var buffer1 bytes.Buffer
	contents := "Simple byte buffer for marshaling data."
	fmt.Printf("Writing contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())   // 39
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap()) // 64

	p1 := make([]byte, 7)
	n, _ := buffer1.Read(p1)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len()) // 32
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
}
