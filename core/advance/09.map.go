package advance

import "fmt"

/*
字典的操作和约束

键-元素对
	Go 语言的字典类型其实是一个哈希表（hash table）的特定实现
	在这个实现中，键和元素的最大不同在于，前者的类型是受限的，而后者却可以是任意类型的
查找
	由于键 - 元素对总是被捆绑在一起存储的，所以一旦找到了键，就一定能找到对应的元素值
	哈希表就会把相应的元素值作为结果返回
	只要这个键 - 元素对存在于哈希表中就一定会被查找到，因为哈希表增、改、删键 - 元素对时侯的映射过程
	映射过程的第一步就是把键值转换为哈希值。在 Go 语言的字典中，每一个键值都是由它的哈希值代表的
	字典不会独立存储任何键的值，但会独立存储它们的哈希值

问题
	字典的键类型不能是哪些类型？
	Go 语言字典的键类型不可以是函数类型、字典类型和切片类型
问题解析
	键值类型
		Go 语言规范规定，在键类型的值之间必须可以施加操作符==和!=，即键类型的值必须要支持判等操作
		由于函数类型、字典类型和切片类型的值并不支持判等操作，所以字典的键类型不能是这些类型
	panic
		如果键的类型是接口类型的，那么键值的实际类型也不能是上述三种类型，否则在程序运行过程中会引发 panic（即运行时恐慌）
		所以最好不要把字典的键类型设定为任何接口类型。如果非要这么做，请一定确保代码在可控的范围之内
		如果键的类型是数组类型，那么还要确保该类型的元素类型不是函数类型、字典类型或切片类型
			由于类型[1][]string的元素类型是[]string，所以它就不能作为字典类型的键类型
		如果键的类型是结构体类型，那么还要保证其中字段的类型的合法性
			比如：m := map[[1][2][3][]string]int
			Invalid map key type: comparison operators == and != must be fully defined for the key type
	为什么键类型的值必须支持判等操作？Go 语言一旦定位到了某一个哈希桶，那么就会试图在这个桶中查找键值
		首先，每个哈希桶都会把自己包含的所有键的哈希值存起来
			Go 语言会用被查找键的哈希值与这些哈希值逐个对比，看看是否有相等的
			如果一个相等的都没有，那么就说明这个桶中没有要查找的键值，这时 Go 语言就会立刻返回结果了
		如果有相等的，那就再用键值本身去对比一次
			不同值的哈希值是可能相同的，即哈希碰撞
		所以，即使哈希值一样，键值也不一定一样
			只有键的哈希值和键值都相等，才能说明查找到了匹配的键 - 元素对
		如果键类型的值之间无法判断相等，那么此时这个映射的过程就没办法继续下去了

知识扩展
问题 1：应该优先考虑哪些类型作为字典的键类型？
	耗时操作
		只从性能的角度看，把键值转换为哈希值”以及“把要查找的键值与哈希桶中的键值做对比”，两个重要且比较耗时的操作
		求哈希和判等操作的速度越快，对应的类型就越适合作为键类型
	对于所有的基本类型、指针类型，以及数组类型、结构体类型和接口类型，Go 语言都有一套算法与之对应
		这套算法中就包含了哈希和判等
	类型的宽度
		类型的宽度是指它的单个值需要占用的字节数
		比如，bool、int8和uint8类型的一个值需要占用的字节数都是1，因此这些类型的宽度就都是1
	以求哈希的操作为例，宽度越小的类型速度通常越快
		对于布尔类型、整数类型、浮点数类型、复数类型和指针类型来说都是如此
		对于字符串类型，由于它的宽度是不定的，所以要看它的值的具体长度，长度越短求哈希越快
		对数组类型的值求哈希实际上是依次求得它的每个元素的哈希值并进行合并，所以速度就取决于它的元素类型以及它的长度
		对结构体类型的值求哈希实际上就是对它的所有字段值求哈希并进行合并，所以关键在于它的各个字段的类型以及字段的数量
		而对于接口类型，具体的哈希算法，则由值的实际类型决定
	不建议使用这些高级数据类型作为字典的键类型
		不仅仅是因为对它们的值求哈希，以及判等的速度较慢，更是因为在它们的值中存在变数
		对一个数组来说，我可以任意改变其中的元素值，但在变化前后，它却代表了两个不同的键值
		对于结构体类型的值情况可能会好一些，因为如果我可以控制其中各字段的访问权限的话，就可以阻止外界修改它了
		把接口类型作为字典的键类型最危险，某个键值不支持判等操作会抛出一个 panic
	建议
		优先选用数值类型和指针类型
		通常情况下类型的宽度越小越好。如果非要选择字符串类型的话，最好对键值的长度进行额外的约束
		不通常的情况，Go 语言有时会对字典的增、删、改、查操作做一些优化
			比如，在字典的键类型为字符串类型的情况下
			又比如，在字典的键类型为宽度为4或8的整数类型的情况下
问题 2：在值为nil的字典上执行读操作会成功吗，那写操作呢？
	由于字典是引用类型，所以当我们仅声明而不初始化一个字典类型的变量的时候，它的值会是nil
	当我们试图在一个值为nil的字典中添加键 - 元素对的时候，Go 语言的运行时系统就会立即抛出一个 panic

总结
	runtime map.go
	sync map.go
	concurrent_map
		https://github.com/easierway/concurrent_map
	参考：08.performance -> lock -> lock_test.go

思考
	在同一时间段内但在不同的 goroutine（或者说 go 程）中对同一个值进行操作是否是安全的
	安全是指，该值不会因这些操作而产生混乱，或其它不可预知的问题
	即字典类型的值是并发安全的吗？
	如果不是，那么在我们只在字典上添加或删除键 - 元素对的情况下，依然不安全吗？
A
	非原子操作需要加锁， map并发读写需要加锁，map操作不是并发安全的
	判断一个操作是否是原子的可以使用 go run race 命令做数据的竞争检测

补充：CGO_ENABLED=1
	官方文档：https://pkg.go.dev/cmd/cgo
	go run -race map.go
	报错
		go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
	管理员权限开启 cmd
		set CGO_ENABLED=1
	go run -race map.go
	报错
		cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
	解决
		https://blog.51cto.com/u_15274085/2918704
		下载地址：https://sourceforge.net/projects/mingw-w64/files/mingw-w64/
	ps：目前只能这样ugly的解决
	Windows
		CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go build main.go
*/

// MapTest -race
func MapTest() {
	//runtime.BlockProfile()
	//sync.Map{}
	var badMap2 = map[interface{}]int{
		"1": 1,
		// 把接口类型作为字典的键类型最危险，某个键值不支持判等操作会抛出一个 panic
		//[]int{2}: 2, // panic: runtime error: hash of unhashable type []int [recovered]
		3: 3,
	}
	fmt.Println(badMap2)

	// 在一个值为nil的字典中添加键 - 元素对，panic
	var m map[int]int
	m[1] = 2 // panic: assignment to entry in nil map [recovered]
}

/*
==================
WARNING: DATA RACE
Write at 0x00c00001e1b0 by goroutine 10:
  runtime.mapassign_fast64()
      E:/Go/src/runtime/map_fast64.go:93 +0x0
  main.write()
      e:/gothmslee/golang/main/map.go:29 +0x48
  main.main.func2()
      e:/gothmslee/golang/main/map.go:15 +0x39

Previous read at 0x00c00001e1b0 by goroutine 7:
  runtime.mapaccess1_fast64()
      E:/Go/src/runtime/map_fast64.go:13 +0x0
  main.read()
      e:/gothmslee/golang/main/map.go:22 +0x3a
  main.main.func1()
      e:/gothmslee/golang/main/map.go:13 +0x39

Goroutine 10 (running) created at:
  main.main()
      e:/gothmslee/golang/main/map.go:15 +0x15a

Goroutine 7 (running) created at:
  main.main()
      e:/gothmslee/golang/main/map.go:13 +0xe6
==================
==================
WARNING: DATA RACE
Write at 0x00c000020388 by goroutine 10:
  main.write()
      e:/gothmslee/golang/main/map.go:29 +0x54
  main.main.func2()
      e:/gothmslee/golang/main/map.go:15 +0x39

Previous read at 0x00c000020388 by goroutine 7:
  main.read()
      e:/gothmslee/golang/main/map.go:22 +0x44
  main.main.func1()
      e:/gothmslee/golang/main/map.go:13 +0x39

Goroutine 10 (running) created at:
  main.main()
      e:/gothmslee/golang/main/map.go:15 +0x15a

Goroutine 7 (running) created at:
  main.main()
      e:/gothmslee/golang/main/map.go:13 +0xe6
==================
fatal error: concurrent map read and map write

goroutine 6 [running]:
main.read(0x0?)
        e:/gothmslee/golang/main/map.go:22 +0x3b
created by main.main
        e:/gothmslee/golang/main/map.go:13 +0xe7

goroutine 1 [sleep]:
time.Sleep(0x6fc23ac00)
        E:/Go/src/runtime/time.go:195 +0x13a
main.main()
        e:/gothmslee/golang/main/map.go:16 +0x16a

goroutine 18 [sleep]:
time.Sleep(0x1)
        E:/Go/src/runtime/time.go:195 +0x13a
main.write(0x0?)
        e:/gothmslee/golang/main/map.go:30 +0x33
created by main.main
        e:/gothmslee/golang/main/map.go:15 +0x15b
exit status 2
*/
