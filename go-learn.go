package golang

/*
加餐｜我“私藏”的那些优质且权威的Go语言学习资料

学习编程语言并没有捷径，脑勤 + 手勤才是正确的学习之路
	如今随着互联网的高速发展，现在很多同学学习编程语言，已经从技术书籍转向了各种屏幕，以专栏或视频实战课为主，技术书籍等参考资料为辅的学习方式已经成为主流
	当然，和传统的、以编程类书籍为主的学习方式相比，谈不上哪种方式更好，只不过它更适合如今快节奏的生活工作状态，更适合碎片化学习占主流的学习形态罢了
	但在编程语言的学习过程中，技术书籍等参考资料依旧是不可或缺的，优秀的参考资料是编程语言学习过程的催化剂，拥有正确的、权威的参考资料可以让你减少反复查找资料所浪费的时间与精力，少走弯路

Go 技术书籍
	和 C（1972 年）、C++（1983）、Java（1995）、Python（1991）等编程语言在市面上的书籍数量相比，Go 流行于市面（尤其是中国大陆地区）上的图书要少很多
		首先，我觉得主要原因还是 Go 语言太年轻了
		其次，Go 以品类代名词的身份占据的“领域”还很少
			提到 Web，人们想到的是 Java Spring
			提到深度学习、机器学习、人工智能，人们想到的是 Python
			提到游戏，人们想到的是 C++
			提到前端，人们想到的是 JavaScript
			这些语言在这些垂直领域早早以杀手级框架入场，使得它们成为了这一领域的“品类代名词”
			Go 虽然通过努力覆盖了一些领域并占据优势地位，比如云原生、API、微服务、区块链，等等，但还不能说已经成为了这些领域的“品类代名词”，因此被垂直领域书籍关联的机会也不像上面那几门语言那么多
		最后是翻译的时间问题
			相对于国内，国外关于 Go 语言的作品要多不少，但引进国外图书资料需要时机以及时间（毕竟要找译者翻译）
	国内 Go 发展
		Go 在国内真正开始快速流行起来，大致是在 2015 年第一届 GopherChina 大会（2015 年 4 月）之后，当时的 Go 版本是 1.4
		同一年下半年发布的 Go 1.5 版本实现了 Go 的自举，并让 GC 延迟大幅下降，让 Go 在国内彻底流行开来
		不过，2020 年开始，国内作者出版的 Go 语言相关书籍已经逐渐多了起来
		2022 年加入泛型的 Go 1.18 版本发布后，相信会有更多 Gopher 加入 Go 技术书籍的写作行列，在未来 3 年，国内 Go 语言技术书籍也会迎来一波高峰
	第五名：《The Way To Go》- Go 语言百科全书
		https://book.douban.com/subject/10558892/
		这本书成书于 2012 年 3 月，恰逢 Go 1.0 版本刚刚发布，当时作者承诺书中代码都可以在 Go 1.0 版本上编译通过并运行
		这本书分为 4 个部分：
			为什么学习 Go 以及 Go 环境安装入门
			Go 语言核心语法
			Go 高级用法（I/O 读写、错误处理、单元测试、并发编程、socket 与 web 编程等)
			Go 应用（常见陷阱、语言应用模式、从性能考量的代码编写建议、现实中的 Go 应用等）
		每部分的每个章节都很精彩，而且这本书也是我目前见到的、最全面详实的、讲解 Go 语言的书籍了，可以说是 Gopher 们的第一本“Go 百科全书”
		Gopher 无闻 在 GitHub 上发起了这本书的中译版项目，如果你感兴趣的话，可以去 GitHub 上看或下载阅读
			https://github.com/Unknwon/the-way-to-go_ZH_CN
	第四名：《Go 101》- Go 语言参考手册
		一本在国外人气和关注度比在国内高的中国人编写的英文书，当然它也是有中文版的
			https://go101.org/article/101.html
		这本书大致可以分为三个部分：
			Go 语法基础
			Go 类型系统与运行时实现
			以专题（topic）形式阐述的 Go 特性、技巧与实践模式
		除了第一部分算 101 范畴，其余两个部分都是 Go 语言的高级话题，也是我们要精通 Go 语言必须要掌握的“知识点”
		并且，作者结合 Go 语言规范，对每个知识点的阐述都细致入微，也结合大量示例进行辅助说明
		《Go 101》这本书，就可以理解为 Go 语言的标准解读或参考手册
			C 和 C++ 语言在市面上都有一些由语言作者或标准规范委员会成员编写的 Annotated 或 Rationale 书籍（语言参考手册或标准解读）
		Go 101 这本书是开源电子书，这就让这本书相对于其他纸板书有着另外一个优势：与时俱进
			https://github.com/go101/go101
			在作者的不断努力下，这本书的知识点更新基本保持与 Go 的演化同步，目前书的内容已经覆盖了最新的 Go 1.17 版本
		这本书的作者是国内资深工程师老貘，他花费三年时间“呕心沥血”完成这本书并且免费奉献给 Go 社区
			近期老貘的两本新书《Go 编程优化 101》和《Go 细节大全 101》也将问世，想必也是不可多得的优秀作品
			https://gfw.tapirgames.com/
	第三名：《Go 语言学习笔记》- Go 源码剖析与实现原理探索
		《Go 语言学习笔记》是一本在国内影响力和关注度都很高的作品
			https://book.douban.com/subject/26832468/
			https://github.com/qyuhen/
			一来，它的作者雨痕老师是国内资深工程师，也是 2015 年第一届 GopherChina 大会讲师
			二来，这部作品的前期版本是以开源电子书的形式分享给国内 Go 社区的
			三来，作者在 Go 源码剖析方面可谓之条理清晰，细致入微
		这本书整体上分为两大部分：
			Go 语言详解：以短平快、“堆干货”的风格对 Go 语言语法做了说明，能用示例说明的，绝不用文字做过多修饰
			Go 源码剖析：这是这本书的精华，也是最受 Gopher 们关注的部分
				这部分对 Go 运行时神秘的内存分配、垃圾回收、并发调度、channel 和 defer 的实现原理、sync.Pool 的实现原理都做了细致的源码剖析与原理总结
		随着 Go 语言的演化，它的语言和运行时实现一直在不断变化，但 Go 1.5 版本的实现是后续版本的基础，所以这本书对它的剖析非常值得每位 Gopher 阅读
	第二名：《Go 语言实战》- 实战系列经典之作，紧扣 Go 语言的精华
		Manning 出版社出版的“实战系列（xx in action）”一直是程序员心中高质量和经典的代名词
			在出版 Go 语言实战系列书籍方面，这家出版社也是丝毫不敢怠慢，邀请了 Go 社区知名的三名明星级作者联合撰写
			威廉·肯尼迪 (William Kennedy) ，知名 Go 培训师，培训机构 Ardan Labs 的联合创始人，“Ultimate Go”培训的策划实施者
			布赖恩·克特森 (Brian Ketelsen) ，世界上最知名的 Go 技术大会 GopherCon 大会的联合发起人和组织者，GopherAcademy创立者，现微软 Azure 工程师
				GopherAcademy创立者：https://gopheracademy.com/
			埃里克·圣马丁 (Erik St.Martin) ，世界上最知名的 Go 技术大会 GopherCon 大会的联合发起人和组织者
		《Go 语言实战》这本书并不是大部头，而是薄薄的一本（中文版才 200 多页），所以你不要期望从本书得到百科全书一样的阅读感
			https://book.douban.com/subject/27015617/
			而且，这本书的作者们显然也没有想把它写成面面俱到的作品，而是直击要点，也就是挑出 Go 语言和其他语言相比与众不同的特点进行着重讲解
			这些特点构成了这本书的结构框架：
				入门：快速上手搭建、编写、运行一个 Go 程序
				语法：数组（作为一个类型而存在）、切片和 map
				Go 类型系统的与众不同：方法、接口、嵌入类型
				Go 的拿手好戏：并发及并发模式
				标准库常用包：log、marshal/unmarshal、io（Reader 和 Writer）
				原生支持的测试
			读完这本书，你就掌握了 Go 语言的精髓之处，这也迎合了多数 Gopher 的内心需求
			而且，这本书中文版译者李兆海也是 Go 圈子里的资深 Gopher，翻译质量上乘
	第一名：《Go 程序设计语言》- 人手一本的 Go 语言“圣经”
		如果说由Brian W. Kernighan和Dennis M. Ritchie联合编写的《The C Programming Language》（也称 K&R C）是 C 程序员（甚至是所有程序员）心目中的“圣经”的话
			那么同样由 Brian W. Kernighan(K) 参与编写的《The Go Programming Language》（也称tgpl）就是 Go 程序员心目中的“圣经”
			https://book.douban.com/subject/1882483/
			https://book.douban.com/subject/26337545/
		这本书模仿并致敬“The C Programming Language”的经典结构，从一个"hello, world"示例开始带领大家开启 Go 语言之旅
			第二章程序结构是 Go 语言这个“游乐园”的向导图
			了解它之后，我们就会迫不及待地奔向各个“景点”细致参观
			Go 语言规范中的所有“景点”在这本书中都覆盖到了，并且由浅入深、循序渐进：
				从基础数据类型到复合数据类型
				从函数、方法到接口
				从创新的并发 Goroutine 到传统的基于共享变量的并发
				从包、工具链到测试
				从反射到低级编程（unsafe 包）
			作者行文十分精炼，字字珠玑，这与《The C Programming Language》的风格保持了高度一致
			而且，书中的示例在浅显易懂的同时，又极具实用性，还突出 Go 语言的特点（比如并发 web 爬虫、并发非阻塞缓存等）
		能得到 Brian W. Kernighan 老爷子青睐的编程语言只有 C 和 Go，这也是 Go 的幸运
			这本书出版于 2015 年 10 月 26 日，也是既当年中旬 Go 1.5 这个里程碑版本发布后，Go 社区的又一重大历史事件
			并且 Brian W. Kernighan 老爷子的影响力让更多程序员加入到 Go 阵营，这也或多或少促成了 Go 成为下一个年度，也就是 2016 年年度 TIOBE 最佳编程语言
			这本书的另一名作者 Alan A. A. Donovan 也并非等闲之辈，他是 Go 核心开发团队的成员，专注于 Go 工具链方面的开发
		这本书的中文版由七牛云团队翻译，总体质量也是不错的
			https://book.douban.com/subject/27044219/

其他形式的参考资料
	Go 官方文档
		最新稳定发布版的文档：https://go.dev/doc/
		项目主线分支（master）上最新开发版本的文档：https://tip.golang.org/
		同时 Go 还将整个 Go 项目文档都加入到了 Go 发行版中，这样开发人员在本地安装 Go 的同时也拥有了一份完整的 Go 项目文档
		不久前，Go 团队已经将原 Go 官方站点 golang.org 重定向到最新开发的 go.dev 网站上
		Go 官方文档中：
			Go 语言规范：https://go.dev/ref/spec
			Go module 参考文档：https://go.dev/ref/mod
			Go 命令参考手册：https://go.dev/doc/cmd
			Effective Go：https://go.dev/doc/effective_go
			Go 标准库包参考手册：https://pkg.go.dev/std
			以及 Go 常见问答：https://go.dev/doc/faq
			等
			都是每个 Gopher 必看的内容
	Go 相关博客
		Go 语言官博，Go 核心团队关于 Go 语言的权威发布渠道
			https://go.dev/blog/
		Go 语言之父 Rob Pike 的个人博客
			https://commandcenter.blogspot.com/
		Go 核心团队技术负责人 Russ Cox 的个人博客
			https://research.swtch.com/
		Go 核心开发者 Josh Bleecher Snyder 的个人博客
			https://commaok.xyz/
		Go 核心团队前成员 Jaana Dogan 的个人博客
			https://rakyll.org/
		Go 鼓吹者 Dave Cheney 的个人博客
			https://dave.cheney.net/
		Go 语言培训机构 Ardan Labs 的博客
			https://www.ardanlabs.com/blog/
		GoCN 社区GoCN 社区
			https://gocn.vip/
		Go 语言百科全书：由欧长坤维护的 Go 语言百科全书网站
			https://golang.design/
	Go 播客
		使用播客这种形式作编程语言类相关内容传播的资料并不多，能持续进行下去的就更少了
		changelog 这个技术类播客平台下的Go Time 频道
		这个频道有几个 Go 社区知名的 Gopher 主持，目前已经播出了 200 多期，每期的嘉宾也都是 Go 社区的重量级人物，其中也不乏像 Go 语言之父这样的大神参与
			https://changelog.com/gotime
	Go 技术演讲
		建议以各大洲举办的 GopherCon 技术大会为主，这些已经基本可以涵盖每年 Go 语言领域的最新发展
		Go 官方的技术演讲归档，这个文档我强烈建议你按时间顺序看一下，通过这些 Go 核心团队的演讲资料，我们可以清晰地了解 Go 的演化历程
			https://go.dev/talks/
		GopherCon 技术大会，这是 Go 语言领域规模最大的技术盛会，也是 Go 官方技术大会
			https://www.youtube.com/c/GopherAcademy/playlists
		GopherCon Europe 技术大会
			https://www.youtube.com/c/GopherConEurope/playlists
		GopherConUK 技术大会
			https://www.youtube.com/c/GopherConUK/playlists
		GoLab 技术大会
			https://www.youtube.com/channel/UCMEvzoHTIdZI7IM8LoRbLsQ/playlists
		Go Devroom@FOSDEM
			https://www.youtube.com/user/fosdemtalks/playlists
		GopherChina 技术大会，这是中国大陆地区规模最大的 Go 语言技术大会，由 GoCN 社区主办
			https://space.bilibili.com/436361287
	Go 日报 / 周刊邮件列表
		通过邮件订阅 Go 语言类日报或周刊，我们也可以获得关于 Go 语言与 Go 社区最新鲜的信息。对于国内的 Gopher 们来说，订阅下面两个邮件列表就足够了：
		Go 语言爱好者周刊，由 Go 语言中文网维护
			https://studygolang.com/go/weekly
		Gopher 日报，由我本人维护的 Gopher 日报项目，创立于 2019 年 9 月
			https://github.com/bigwhite/gopherdaily
	其他
		Go 语言项目的官方 issue 列表
			通过这个 issue 列表，我们可以实时看到 Go 项目演进状态，及时看到 Go 社区提交的各种 bug
			同时，我们通过挖掘该列表，还可以了解某个 Go 特性的来龙去脉，这对深入理解 Go 特性很有帮助
			https://github.com/golang/go/issues
		Go 项目的代码 review 站点
			通过阅读 Go 核心团队 review 代码的过程与评审意见，我们可以看到 Go 核心团队是如何使用 Go 进行编码的
			能够学习到很多 Go 编码原则以及地道的 Go 语言惯用法，对更深入地理解 Go 语言设计哲学，形成 Go 语言编程思维有很大帮助
			https://go-review.googlesource.com/q/status:open+-is:wip
*/
