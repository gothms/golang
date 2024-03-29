package advance

import (
	"container/list"
	"fmt"
)

/*
container包中的那些容器

链表
	Go 语言的链表实现在标准库的container/list代码包中
	有两个公开的程序实体——List和Element
	List 实现了一个双向链表，而 Element 则代表了链表中元素的结构

问题
	可以把自己生成的Element类型值传给链表吗？
	API：
		MoveBefore方法和MoveAfter方法：分别用于把给定的元素移动到另一个元素的前面和后面
		MoveToFront方法和MoveToBack方法，分别用于把给定的元素移动到链表的最前端和最后端
	A
		不会接受，这些方法将不会对链表做出任何改动
		因为我们自己生成的Element值并不在链表中，所以也就谈不上“在链表中移动元素”
		更何况链表不允许我们把自己生成的Element值插入其中
问题解析
	为什么不会接受？
		在List包含的方法中，用于插入新元素的那些方法都只接受interface{}类型的值
		这些方法在内部会使用Element值，包装接收到的新元素
		这样做正是为了避免直接使用我们自己生成的元素，主要原因是避免链表的内部关联，遭到外界破坏，这对于链表本身以及我们这些使用者来说都是有益的
	other API：
		Front和Back方法分别用于获取链表中最前端和最后端的元素
		InsertBefore和InsertAfter方法分别用于在指定的元素之前和之后插入新元素
		PushFront和PushBack方法则分别用于在链表的最前端和最后端插入新元素
	安全“接口”
		这些方法都会把一个Element值的指针作为结果返回，它们就是链表留给我们的安全“接口”
		拿到这些内部元素的指针，我们就可以去调用前面提到的用于移动元素的方法了

知识扩展
1. 问题：为什么链表可以做到开箱即用？
	List和Element都是结构体类型
	结构体类型特点
		就是它们的零值都会是拥有特定结构，但是没有任何定制化内容的值，相当于一个空壳
		值中的字段也都会被分别赋予各自类型的零值
		广义来讲，所谓的零值就是只做了声明，但还未做初始化的变量被给予的缺省值。每个类型的零值都会依据该类型的特性而被设定
		示例
			经过语句var a [2]int声明的变量a的值，将会是一个包含了两个0的整数数组
			经过语句var s []int声明的变量s的值将会是一个[]int类型的、值为nil的切片
	“开箱即用”
		声明变量 var l list.List
			这个零值将会是一个长度为0的链表
			这个链表持有的根元素也将会是一个空壳，其中只会包含缺省的内容
		这样的链表我们可以直接拿来使用，这被称为“开箱即用”
			Go 语言标准库中很多结构体类型的程序实体都做到了开箱即用
		这也是在编写可供别人使用的代码包（或者说程序库）时，我们推荐遵循的最佳实践之一
		关键在于它的“延迟初始化”机制
	“延迟初始化”机制
		可以理解为把初始化操作延后，仅在实际需要的时候才进行
		延迟初始化的优点在于“延后”，它可以分散初始化操作带来的计算量和存储空间消耗
	切片：延迟初始化
		如果我们需要集中声明非常多的大容量切片的话，那么那时的 CPU 和内存空间的使用量肯定都会一个激增，并且只有设法让其中的切片及其底层数组被回收，内存使用量才会有所降低
		如果数组是可以被延迟初始化的，那么计算量和存储空间的压力就可以被分散到实际使用它们的时候
		这些数组被实际使用的时间越分散，延迟初始化带来的优势就会越明显
		Go 语言的切片就起到了延迟初始化其底层数组的作用
	延迟初始化的缺点
		延迟初始化的缺点恰恰也在于“延后”
		如链表：
		在调用链表的每个方法的时候，它们都需要先去判断链表是否已经被初始化，那这也会是一个计算量上的浪费
		在这些方法被非常频繁地调用的情况下，这种浪费的影响就开始显现了，程序的性能将会降低
	链表：延迟初始化
		链表实现中，一些方法是无需对是否初始化做判断的
		Front方法和Back方法，一旦发现链表的长度为0，直接返回nil就好了
		用于删除元素、移动元素，以及一些用于插入元素的方法中，只要判断一下传入的元素中指向所属链表的指针，是否与当前链表的指针相等就可以了
			如果不相等，就一定说明传入的元素不是这个链表中的，后续的操作就不用做了
			反之，就一定说明这个链表已经被初始化了
		原因在于，链表的PushFront方法、PushBack方法、PushBackList方法以及PushFrontList方法总会先判断链表的状态，并在必要时进行初始化，这就是延迟初始化
			在向一个空的链表中添加新元素的时候，肯定会调用这四个方法中的一个，这时新元素中指向所属链表的指针，一定会被设定为当前链表的指针
			所以，指针相等是链表已经初始化的充分必要条件
		List利用了自身以及Element在结构上的特点，巧妙地平衡了延迟初始化的优缺点，使得链表可以开箱即用，并且在性能上可以达到最优
2. 问题：Ring与List的区别在哪儿？
	Ring
		container/ring包中的Ring类型实现的是一个循环链表，也就是我们俗称的环
	List本质
		List在内部就是一个循环链表
		它的根元素永远不会持有任何实际的元素值，而该元素的存在就是为了连接这个循环链表的首尾两端
		List的零值是一个只包含了根元素，但不包含任何实际元素值的空链表
	Ring vs List
		1. Ring类型的数据结构仅由它自身即可代表，而List类型则需要由它以及Element类型联合表示
			这是表示方式上的不同，也是结构复杂度上的不同
		2. 一个Ring类型的值严格来讲，只代表了其所属的循环链表中的一个元素，而一个List类型的值则代表了一个完整的链表
			这是表示维度上的不同
		3. 在创建并初始化一个Ring值的时候，我们可以指定它包含的元素的数量，但是对于一个List值来说却不能这样做（也没有必要这样做）
			循环链表一旦被创建，其长度是不可变的
			这是两个代码包中的New函数在功能上的不同，也是两个类型在初始化值方面的第一个不同
		4. 仅通过var r ring.Ring语句声明的r将会是一个长度为1的循环链表，而List类型的零值则是一个长度为0的链表
			List中的根元素不会持有实际元素值，因此计算长度时不会包含它
			这是两个类型在初始化值方面的第二个不同
		5. Ring值的Len方法的算法复杂度是 O(N) 的，而List值的Len方法的算法复杂度则是 O(1) 的
			这是两者在性能方面最显而易见的差别
		其他的不同基本上都是方法方面的了。比如，循环链表也有用于插入、移动或删除元素的方法，不过用起来都显得更抽象一些

总结
	List
		List这个结构体类型有两个字段，一个是Element类型的字段root，另一个是int类型的字段len
		前者代表的就是那个根元素，而后者用于存储链表的长度
		它们都是包级私有的，也就是说使用者无法查看和修改它们
		字段root和len都会被赋予相应的零值
		len的零值是0，正好可以表明该链表还未包含任何元素。由于root是Element类型的，所以它的零值就是该类型的空壳，用字面量表示的话就是Element{}
	Element
		Element类型包含了几个包级私有的字段，分别用于存储前一个元素、后一个元素以及所属链表的指针值
		一个名叫Value的公开的字段，该字段的作用就是持有元素的实际值，它是interface{}类型的
		在Element类型的零值中，这些字段的值都会是nil

思考
	1.container/ring包中的循环链表的适用场景都有哪些？
	2.你使用过container/heap包中的堆吗？它的适用场景又有哪些呢？

补充
	切片 & 数组
		切片本身有着占用内存少和创建便捷等特点，但它的本质上还是数组
		切片的一大好处是可以让我们通过窗口快速地定位并获取，或者修改底层数组中的元素
	内存泄漏/溢出
		删除切片中的元素时，要注意空出的元素槽位的“清空”，否则很可能会造成内存泄漏
		另一方面，在切片被频繁“扩容”的情况下，新的底层数组会不断产生，这时内存分配的量以及元素复制的次数可能就很可观了，这肯定会对程序的性能产生负面的影响
		尤其是当我们没有一个合理、有效的”缩容“策略的时候，旧的底层数组无法被回收，新的底层数组中也会有大量无用的元素槽位
		过度的内存浪费不但会降低程序的性能，还可能会使内存溢出并导致程序崩溃
	切片 & 链表
		一个链表所占用的内存空间，往往要比包含相同元素的数组所占内存大得多
		这是由于链表的元素并不是连续存储的，所以相邻的元素之间需要互相保存对方的指针
		不但如此，每个元素还要存有它所属链表的指针
		有了这些关联，链表的结构反倒更简单了。它只持有头部元素（或称为根元素）基本上就可以了
		当然了，为了防止不必要的遍历和计算，链表的长度记录在内也是必须的

Q
	为什么源码中，PushBackList 和 PushFrontList 要一个个添加元素，而不一次性添加整个链？
*/

// List 测试链表
func List() {
	l := list.List{}
	l.PushFront(3)
	fmt.Println(l, l.Front())
}

// LazyInit 测试延迟初始化
func LazyInit() {
	var arr []int
	arr = append(arr, 0)
	fmt.Println(arr)
}
