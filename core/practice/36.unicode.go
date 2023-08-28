package practice

import "fmt"

/*
unicode与字符编码

参考
	05.scope.go & 06.assert.go
	Go 语言中的标识符可以包含任何 Unicode 编码可以表示的字母字符
	虽然我们可以直接把一个整数值转换为一个string类型的值，但是被转换的整数值应该可以代表一个有效的 Unicode 代码点
	否则转换的结果就将会是"�"，即：一个仅由高亮的问号组成的字符串值
	另外，当一个string类型的值被转换为[]rune类型值的时候，其中的字符串会被拆分成一个一个的 Unicode 字符
Unicode & UTF-8
	Go 语言采用的字符编码方案从属于 Unicode 编码规范，即 Go 语言的代码正是由 Unicode 字符组成的
	Go 语言的所有源代码，都必须按照 Unicode 编码规范中的 UTF-8 编码格式进行编码
	即 Go 语言的源码文件必须使用 UTF-8 编码格式进行存储
	如果源码文件中出现了非 UTF-8 编码的字符，那么在构建、安装以及运行的时候，go 命令就会报告错误“illegal UTF-8 encoding”
ASCII 编码
	简介
		ASCII 是英文“American Standard Code for Information Interchange”的缩写，中文译为美国信息交换标准代码
		它是由美国国家标准学会（ANSI）制定的单字节字符编码方案，可用于基于文本的数据交换
		它最初是美国的国家标准，后又被国际标准化组织（ISO）定为国际标准，称为 ISO 646 标准，并适用于所有的拉丁文字字母
	ASCII 编码集
		ASCII 编码方案使用单个字节（byte）的二进制数来编码一个字符
		标准的 ASCII 编码用一个字节的最高比特（bit）位作为奇偶校验位，而扩展的 ASCII 编码则将此位也用于表示字符
		ASCII 编码支持的可打印字符和控制字符的集合也被叫做 ASCII 编码集
	Unicode 编码规范
		Unicode 编码规范，实际上是另一个更加通用的、针对书面字符和文本的字符编码标准
		它为世界上现存的所有自然语言中的每一个字符，都设定了一个唯一的二进制编码
		它定义了不同自然语言的文本数据在国际间交换的统一方式，并为全球化软件创建了一个重要的基础
	Unicode 与 ASCII
		Unicode 编码规范以 ASCII 编码集为出发点，并突破了 ASCII 只能对拉丁字母进行编码的限制
		它不但提供了可以对世界上超过百万的字符进行编码的能力，还支持所有已知的转义序列和控制代码
	代码空间、代码点
		在计算机系统的内部，抽象的字符会被编码为整数。这些整数的范围被称为代码空间
		在代码空间之内，每一个特定的整数都被称为一个代码点
		一个受支持的抽象字符会被映射并分配给某个特定的代码点，反过来讲，一个代码点总是可以被看成一个被编码的字符
	U+
		Unicode 编码规范通常使用十六进制表示法来表示 Unicode 代码点的整数值，并使用“U+”作为前缀
		比如，英文字母字符“a”的 Unicode 代码点是 U+0061
		在 Unicode 编码规范中，一个字符能且只能由与它对应的那个代码点表示
	版本
		Unicode 编码规范现在的最新版本是 11.0，并会于 2019 年 3 月发布 12.0 版本
		而 Go 语言从 1.10 版本开始，已经对 Unicode 的 10.0 版本提供了全面的支持
		对于绝大多数的应用场景来说，这已经完全够用了
	三种不同的编码格式
		Unicode 编码规范提供了三种不同的编码格式，即：UTF-8、UTF-16 和 UTF-32
			其中的 UTF 是 UCS Transformation Format 的缩写
			而 UCS 又是 Universal Character Set 的缩写，但也可以代表 Unicode Character Set
			所以，UTF 也可以被翻译为 Unicode 转换格式。它代表的是字符与字节序列之间的转换方式
		“-”右边的整数的含义是，以多少个比特位作为一个编码单元
			以 UTF-8 为例，它会以 8 个比特，也就是一个字节，作为一个编码单元
			并且，它与标准的 ASCII 编码是完全兼容的
			也就是说，在 [0x00, 0x7F] 的范围内，这两种编码表示的字符都是相同的。这也是 UTF-8 编码格式的一个巨大优势
	UTF-8
		UTF-8 是一种可变宽的编码方案。换句话说，它会用一个或多个字节的二进制数来表示某个字符，最多使用四个字节
		比如，对于一个英文字符，它仅用一个字节的二进制数就可以表示，而对于一个中文字符，它需要使用三个字节才能够表示
		不论怎样，一个受支持的字符总是可以由 UTF-8 编码为一个字节序列
问题：一个string类型的值在底层是怎样被表达的？
	在底层，一个string类型的值是由一系列相对应的 Unicode 代码点的 UTF-8 编码值来表达的
问题解析
	在 Go 语言中，一个string类型的值既可以被拆分为一个包含多个字符的序列，也可以被拆分为一个包含多个字节的序列
		字符序列：可以由一个以rune为元素类型的切片来表示
		字节序列：可以由一个以byte为元素类型的切片代表
	rune
		rune是 Go 语言特有的一个基本数据类型，它的一个值就代表一个字符，即：一个 Unicode 字符
		比如，'G'、'o'、'爱'、'好'、'者'代表的就都是一个 Unicode 字符
		UTF-8 编码方案会把一个 Unicode 字符编码为一个长度在 [1, 4] 范围内的字节序列
		所以，一个rune类型的值也可以由一个或多个字节来代表
	type rune = int32
		一个rune类型的值会由四个字节宽度的空间来存储。它的存储空间总是能够存下一个 UTF-8 编码值
		一个rune类型的值在底层其实就是一个 UTF-8 编码值
		rune 是便于我们人类理解的外部展现，UTF-8 是便于计算机系统理解的内在表达
	示例：RuneAndUTF8() 测试一
		runes(hex)：字符序列
		bytes(hex)：字节序列
		一个string类型的值会由若干个 Unicode 字符组成，每个 Unicode 字符都可以由一个rune类型的值来承载
	小结
		这些字符在底层都会被转换为 UTF-8 编码值，而这些 UTF-8 编码值又会以字节序列的形式表达和存储
		因此，一个string类型的值在底层就是一个能够表达若干个 UTF-8 编码值的字节序列

知识扩展
问题 1：使用带有range子句的for语句遍历字符串值的时候应该注意什么？
	遍历
		带有range子句的for语句会先把被遍历的字符串值拆成一个字节序列
		然后再试图找出这个字节序列中包含的每一个 UTF-8 编码值，或者说每一个 Unicode 字符
	示例：RuneAndUTF8() 测试二
		for语句可以逐一地迭代出字符串值里的每个 Unicode 字符
		但是，相邻的 Unicode 字符的索引值并不一定是连续的。这取决于前一个 Unicode 字符是否为单字节字符
		如果我们想得到其中某个 Unicode 字符对应的 UTF-8 编码值的宽度，就可以用下一个字符的索引值减去当前字符的索引值
	参考
		api -> string_test.go

思考
	判断一个 Unicode 字符是否为单字节字符通常有几种方式？
*/

func RuneAndUTF8() {
	// 测试一
	str := "Go 爱好者 "
	fmt.Printf("The string: %q\n", str)                 // The string: "Go 爱好者 "
	fmt.Printf("  => runes(char): %q\n", []rune(str))   // => runes(char): ['G' 'o' ' ' '爱' '好' '者' ' ']
	fmt.Printf("  => runes(hex): %x\n", []rune(str))    // => runes(hex): [47 6f 20 7231 597d 8005 20]
	fmt.Printf("  => bytes(hex): [% x]\n", []byte(str)) // => bytes(hex): [47 6f 20 e7 88 b1 e5 a5 bd e8 80 85 20]
	fmt.Printf("  => bytes(b): [%b]\n", []byte(str))    // => bytes(b): [[1000111 1101111 100000 11100111 10001000 10110001 11100101 10100101 10111101 11101000 10000000 10000101 100000]]

	// 测试二
	//0: 'G' [47]
	//1: 'o' [6f]
	//2: ' ' [20]
	//3: '爱' [e7 88 b1]
	//6: '好' [e5 a5 bd]
	//9: '者' [e8 80 85]
	//12: ' ' [20]
	for i, c := range str {
		fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c)))
	}
}
