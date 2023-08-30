package practice

/*
io包中的接口和工具 （下）

知识扩展
问题：io包中的接口都有哪些？它们之间都有着怎样的关系？
	简单接口
		没有嵌入其他接口并且只定义了一个方法的接口叫做简单接口
		在io包中，这样的接口一共有 11 个
	核心接口
		核心接口有着众多的扩展接口和实现类型
		io包中的核心接口只有 3 个，它们是：io.Reader、io.Writer和io.Closer
	分类
		把io包中的简单接口分为四大类，分别针对于四种操作，即：读取、写入、关闭和读写位置设定
		前三种操作属于基本的 I/O 操作
		实际上，在io包中，与写入操作有关的接口都与读取操作的相关接口有着一定的对应关系
	读取
		核心接口io.Reader
			它在io包中有 5 个扩展接口，并有 6 个实现类型
		io.ByteReader和io.RuneReader
			它们分别定义了一个读取方法，即：ReadByte和ReadRune
			与io.Reader接口中Read方法不同的是，这两个读取方法分别只能够读取下一个单一的字节和 Unicode 字符
			实现类型
				数据类型strings.Reader和bytes.Buffer都是io.ByteReader和io.RuneReader的实现类型
			上层接口
				都实现了io.ByteScanner接口和io.RuneScanner接口
				io.ByteScanner接口内嵌了简单接口io.ByteReader，并定义了额外的UnreadByte方法
					如此一来，它就抽象出了一个能够读取和读回退单个字节的功能集
				io.RuneScanner内嵌了简单接口io.RuneReader，并定义了额外的UnreadRune方法
					它抽象的是可以读取和读回退单个 Unicode 字符的功能集
		io.ReaderAt接口
			只定义了一个方法ReadAt，ReadAt是一个纯粹的只读方法
			它只去读取其所属值中包含的字节，而不对这个值进行任何的改动，比如，它绝对不能去修改已读计数的值
			这也是io.ReaderAt接口与其实现类型之间最重要的一个约定
			因此，如果仅仅并发地调用某一个值的ReadAt方法，那么安全性应该是可以得到保障的
		io.ReaderFrom接口
			定义了一个名叫ReadFrom的写入方法
			方法会接受一个io.Reader类型的参数值，并会从该参数值中读出数据，并写入到其所属值中
			io.CopyN 函数
				在复制数据的时候会先检测其参数src的值是否实现了io.WriterTo接口。如果是，那么它就直接利用该值的WriteTo方法，把其中的数据拷贝给参数dst代表的值
				类似的，这个函数还会检测dst的值是否实现了io.ReaderFrom接口。如果是，那么它就会利用这个值的ReadFrom方法，直接从src那里把数据拷贝进该值
				io.CopyBuffer函数也是如此，它们在内部做数据复制的时候用的都是同一套代码
			io.ReaderFrom接口与io.WriterTo接口对应得很规整
	写入
		核心接口io.Writer
			基于它的扩展接口
			io.ReadWriter、io.ReadWriteCloser、io.ReadWriteSeeker、io.WriteCloser和io.WriteSeeker
		io.ReadWriter接口
			*io.pipe是io.ReadWriter接口的实现类型
		io.ReadWriteCloser
			在io包中并没有io.ReadWriteCloser接口的实现，它的实现类型主要集中在net包中
		io.ByteWriter和io.WriterAt
			io包中也没有它们的实现类型
			*os.File类型不但是io.WriterAt接口的实现类型，还同时实现了io.ReadWriteCloser接口和io.ReadWriteSeeker接口
			该类型支持的 I/O 操作非常的丰富
	关闭
		io.Closer接口
			非常通用，它的扩展接口和实现类型都不少
		实现类型
			io包中只有io.PipeReader和io.PipeWriter
	读写位置设定
		io.Seeker接口
			仅仅定义了一个方法Seek
			该方法主要用于寻找并设定下一次读取或写入时的起始索引位置
			strings.Reader类型是它的实现类型
		io.ReadSeeker和io.ReadWriteSeeker
			基于io.Seeker的扩展接口
		io.WriteSeeker
			基于io.Writer和io.Seeker的扩展接口
		io.Seeker接口的其他实现类型
			两个指针类型strings.Reader和io.SectionReader
			它们也是io.ReaderAt接口的实现类型

总结
	简单分类
		根据接口定义的方法的数量以及是否有接口嵌入，把io包中的接口分为了简单接口和扩展接口
		根据这些简单接口的扩展接口和实现类型的数量级，把它们分为了核心接口和非核心接口
	io包核心接口
		io.Reader、io.Writer和io.Closer
		这些核心接口在 Go 语言标准库中的实现类型都在 200 个以上
	根据针对的 I/O 操作的不同，把简单接口分为了四大类
		针对的操作分别是：读取、写入、关闭和读写位置设定
		其中，前三种操作属于基本的 I/O 操作
	程序实体的功用和机理
		数据段读取器io.SectionReader
		作为同步内存管道核心实现的io.pipe类型
		以及用于数据拷贝的io.CopyN函数
		...
	概括
		io包中的简单接口共有 11 个
			读取操作相关的接口有 5 个
			写入操作相关的接口有 4 个
			与关闭操作有关的接口只有 1 个
			另外还有一个读写位置设定相关的接口
		io包还包含了 9 个基于这些简单接口的扩展接口

思考
	io包中的同步内存管道的运作机制是什么？
*/
