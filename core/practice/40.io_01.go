package practice

/*
io包中的接口和工具 （上）

strings.Builder类型主要用于构建字符串，实现的接口
	io.Writer、io.ByteWriter、io.StringWriter、fmt.Stringer
strings.Reader类型主要用于读取字符串，实现的接口
	io.Reader
	io.ReaderAt
	io.ByteReader
	io.RuneReader
	io.Seeker
	io.ByteScanner：io.ByteReader的扩展接口
	io.RuneScanner：io.RuneReader的扩展接口
	io.WriterTo
bytes.Buffer类型是集读、写功能的数据类型，实现的接口
	读取相关：
		io.Reader
		io.ByteReader
		io.RuneReader
		io.ByteScanner
		io.RuneScanner
		io.WriterTo
	写入相关：
		io.Writer
		io.ByteWriter
		io.StringWriter
		io.ReaderFrom
	导出相关：
		fmt.Stringer
实现了这么多的接口，其目的究竟是什么呢？
	为了提高不同程序实体之间的互操作性
	举例
		在io包中，几个用于拷贝数据的函数，io.Copy、io.CopyBuffer和io.CopyN
			虽然这几个函数在功能上都略有差别，但是它们都首先会接受两个参数
			即：用于代表数据目的地、io.Writer类型的参数dst，以及用于代表数据来源的、io.Reader类型的参数src
			这些函数的功能大致上都是把数据从src拷贝到dst
			第一个参数值，只要实现了io.Writer接口即可，第二个参数值，只要实现了io.Reader接口就行
		当参数被传到io.CopyN函数时，就已经分别被包装成了io.Reader类型和io.Writer类型的值
			为了优化的目的，io.CopyN函数中的代码会对参数值进行再包装，也会检测这些参数值是否还实现了别的接口
			甚至还会去探求某个参数值被包装后的实际类型，是否为某个特殊的类型
	面向接口编程
		io.CopyN函数通过面向接口编程，极大地拓展了它的适用范围和应用场景
		换个角度看，正因为strings.Reader类型和strings.Builder类型都实现了不少接口，所以它们的值才能够被使用在更广阔的场景中
		如此一来，Go 语言的各种库中，能够操作它们的函数和数据类型明显多了很多
	io包中的接口对于 Go 语言的标准库和很多第三方库而言，都起着举足轻重的作用
		io.Reader和io.Writer这两个最核心的接口，它们是很多接口的扩展对象和设计源泉
		单从 Go 语言的标准库中统计，实现了它们的数据类型都（各自）有上百个，而引用它们的代码更是都（各自）有 400 多处
		很多数据类型实现了io.Reader接口，是因为它们提供了从某处读取数据的功能
		类似的，许多能够把数据写入某处的数据类型，也都会去实现io.Writer接口
		其实，有不少类型的设计初衷都是实现这两个核心接口的某个或某些扩展接口，以提供比单纯的字节序列读取或写入更加丰富的功能
	接口设计
		在 Go 语言中，对接口的扩展是通过接口类型之间的嵌入来实现的，这也常被叫做接口的组合
		我在讲接口的时候也提到过，Go 语言提倡使用小接口加接口组合的方式，来扩展程序的行为以及增加程序的灵活性
		io代码包恰恰就可以作为这样的一个标杆，它可以成为我们运用这种技巧时的一个参考标准

问题：在io包中，io.Reader的扩展接口和实现类型都有哪些？它们分别都有什么功用？
	扩展接口
		1. io.ReadWriter：此接口既是io.Reader的扩展接口，也是io.Writer的扩展接口
			该接口定义了一组行为，包含且仅包含了基本的字节序列读取方法Read，和字节序列写入方法Write
		2. io.ReadCloser：此接口除了包含基本的字节序列读取方法之外，还拥有一个基本的关闭方法Close
			Close一般用于关闭数据读写的通路。这个接口其实是io.Reader接口和io.Closer接口的组合
		3. io.ReadWriteCloser：很明显，此接口是io.Reader、io.Writer和io.Closer这三个接口的组合
		4. io.ReadSeeker：此接口的特点是拥有一个用于寻找读写位置的基本方法Seek
			该方法可以根据给定的偏移量基于数据的起始位置、末尾位置，或者当前读写位置去寻找新的读写位置
			这个新的读写位置用于表明下一次读或写时的起始索引。Seek是io.Seeker接口唯一拥有的方法
		5. io.ReadWriteSeeker：显然，此接口是另一个三合一的扩展接口，它是io.Reader、io.Writer和io.Seeker的组合
	实现类型
		1. *io.LimitedReader：此类型的基本类型会包装io.Reader类型的值，并提供一个额外的受限读取的功能
			所谓的受限读取指的是，此类型的读取方法Read返回的总数据量会受到限制，无论该方法被调用多少次
			这个限制由该类型的字段N指明，单位是字节
		2. *io.SectionReader：此类型的基本类型可以包装io.ReaderAt类型的值，并且会限制它的Read方法，只能够读取原始数据中的某一个部分（或某一段）
			这个数据段的起始位置和末尾位置，需要在它被初始化的时候就指明，并且之后无法变更
			该类型值的行为与切片有些类似，它只会对外暴露在其窗口之中的那些数据
		3. *io.teeReader：此类型是一个包级私有的数据类型，也是io.TeeReader函数结果值的实际类型。这个函数接受两个参数r和w，类型分别是io.Reader和io.Writer
			其结果值的Read方法会把r中的数据经过作为方法参数的字节切片p写入到w
			可以说，这个值就是r和w之间的数据桥梁，而那个参数p就是这座桥上的数据搬运者
		4. io.multiReader：此类型也是一个包级私有的数据类型。类似的，io包中有一个名为MultiReader的函数，它可以接受若干个io.Reader类型的参数值，并返回一个实际类型为io.multiReader的结果值
			当这个结果值的Read方法被调用时，它会顺序地从前面那些io.Reader类型的参数值中读取数据
			因此，我们也可以称之为多对象读取器
		5. io.pipe：此类型为一个包级私有的数据类型，它比上述类型都要复杂得多。它不但实现了io.Reader接口，而且还实现了io.Writer接口
			实际上，io.PipeReader类型和io.PipeWriter类型拥有的所有指针方法都是以它为基础的。这些方法都只是代理了io.pipe类型值所拥有的某一个方法而已
			又因为io.Pipe函数会返回这两个类型的指针值并分别把它们作为其生成的同步内存管道的两端，所以可以说，*io.pipe类型就是io包提供的同步内存管道的核心实现
		6. io.PipeReader：此类型可以被视为io.pipe类型的代理类型。它代理了后者的一部分功能，并基于后者实现了io.ReadCloser接口
			同时，它还定义了同步内存管道的读取端
问题解析
	io 代码包是 Go 语言标准库中所有 I/O 相关 API 的根基
		通过io.Reader接口，我们应该能够梳理出基于它的类型树，并知晓其中每一个类型的功用
		io.Reader可谓是io包乃至是整个 Go 语言标准库中的核心接口，所以我们可以从它那里牵扯出很多扩展接口和实现类型
		在很多时候，我们可以根据实际需求将它们搭配起来使用
		例如，对施加在原始数据之上的（由Read方法提供的）读取功能进行多层次的包装（比如受限读取和多对象读取等），以满足较为复杂的读取需求
	示例
		demo82.go

思考
	你用过哪些io包中的接口和工具呢，又有哪些收获和感受呢
*/
