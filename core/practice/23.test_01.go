package practice

/*
测试的基本规则和流程 （上）

测试函数
	测试函数往往用于描述和保障某个程序实体的某方面功能
	比如，该功能在正常情况下会因什么样的输入，产生什么样的输出，又比如，该功能会在什么情况下报错或表现异常，等
可以为 Go 程序编写三类测试
	功能测试（test）、基准测试（benchmark，也称性能测试），以及示例测试（example）
	示例测试严格来讲也是一种功能测试，只不过它更关注程序打印出来的内容
测试源码文件
	一般情况下，一个测试源码文件只会针对于某个命令源码文件，或库源码文件做测试
		所以我们总会（并且应该）把它们放在同一个代码包内
	命名
		测试源码文件的主名称应该以被测源码文件的主名称为前导，并且必须以“_test”为后缀
	每个测试源码文件都必须至少包含一个测试函数
		从语法上讲，每个测试源码文件中，都可以包含用来做任何一类测试的测试函数，即使把这三类测试函数都塞进去也没有问题
		只要把控好测试函数的分组和数量就可以了
	组织
		可以依据这些测试函数针对的不同程序实体，把它们分成不同的逻辑组，并且，利用注释以及帮助类的变量或函数来做分割
		还可以依据被测源码文件中程序实体的先后顺序，来安排测试源码文件中测试函数的顺序
	不仅仅对测试源码文件的名称，对于测试函数的名称和签名，Go 语言也是有明文规定的

问题：Go 语言对测试函数的名称和签名都有哪些规定？
	对于功能测试函数来说，其名称必须以Test为前缀，并且参数列表中只应有一个*testing.T类型的参数声明
	对于性能测试函数来说，其名称必须以Benchmark为前缀，并且唯一参数的类型必须是*testing.B类型的
	对于示例测试函数来说，其名称必须以Example为前缀，但对函数的参数列表没有强制规定
问题解析
	go test
		只有测试源码文件的名称对了，测试函数的名称和签名也对了，当我们运行go test命令的时候，其中的测试代码才有可能被运行
	go test命令执行的主要测试流程是什么？
		go test命令在开始运行时，会先做一些准备工作
			比如，确定内部需要用到的命令，检查我们指定的代码包或源码文件的有效性，以及判断我们给予的标记是否合法，等
		通常情况下的主要测试流程
			在准备工作顺利完成之后，go test命令就会针对每个被测代码包，依次地进行构建、执行包中符合要求的测试函数，清理临时文件，打印测试结果
			“串行”：
				对于每个被测代码包，go test命令会串行地执行测试流程中的每个步骤
				但是，为了加快测试速度，它通常会并发地对多个被测代码包进行功能测试
				只不过，在最后打印测试结果的时候，它会依照我们给定的顺序逐个进行，这会让我们感觉到它是在完全串行地执行测试流程
			性能测试
				由于并发的测试会让性能测试的结果存在偏差，所以性能测试一般都是串行进行的
				更具体地说，只有在所有构建步骤都做完之后，go test命令才会真正地开始进行性能测试
				并且，下一个代码包性能测试的进行，总会等到上一个代码包性能测试的结果打印完成才会开始，而且性能测试函数的执行也都会是串行的
			所以，即使是简单的性能测试，执行起来也会比功能测试慢

总结
	中小型的公司，他们往往完全依靠软件质量保障团队，甚至真正的用户去帮他们测试
	在这些情况下，软件错误或缺陷的发现、反馈和修复的周期通常会很长，成本也会很大，也许还会造成很不好的影响
	Go 语言
		是一门很重视程序测试的编程语言，它不但自带了testing包，还有专用于程序测试的命令go test

思考
	你还知道或用过testing.T类型和testing.B类型的哪些方法？它们都是做什么用的？
*/
