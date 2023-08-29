package advance

import "fmt"

/*
if语句、for语句和switch语句

问题：使用携带range子句的for语句时需要注意哪些细节？
问题解析
	range表达式的结果值可以是数组、数组的指针、切片、字符串、字典或者允许接收操作的通道，并且结果值只能有一个
	注意
		range表达式只会在for语句开始执行时被求值一次，无论后边会有多少次迭代
		range表达式的求值结果会被复制，也就是说，被迭代的对象是range表达式结果值的副本而不是原值
	示例：RangeArray()
		数组：值类型
		切片：引用类型

知识扩展
问题 1：switch语句中的switch表达式和case表达式之间有着怎样的联系？
	类型相同
		switch 表达式的结果类型，必须与各个 case 表达式的结果类型相同
	无类型常量自动转换
		如果case表达式中子表达式的结果值是无类型的常量，那么它的类型会被自动地转换为switch表达式的结果类型
		如果自动转换没能成功，那么switch语句照样通不过编译
	接口类型
		如果这些表达式的结果类型有某个接口类型，那么一定要小心检查它们的动态值是否都具有可比性（或者说是否允许判等操作）
		如果答案是否定的，虽然不会造成编译错误，但是后果会更加严重：引发 panic（也就是运行时恐慌）
问题 2：switch语句对它的case表达式有哪些约束？
	switch语句不允许case表达式中的子表达式结果值存在相等的情况
		不论这些结果值相等的子表达式，是否存在于不同的case表达式中，都会是这样的结果
		不过这只是对于由字面量直接表示的子表达式而言的
	绕过case子表达式结果值相等的限制
		示例：SwitchCase()
	类型判断的switch语句
		byte 和 uint8 是同类型，编译不通过
	最上边的case子句中的子表达式总是会被最先求值，在判等的时候顺序也是这样
		如果某些子表达式的结果值有重复并且它们与switch表达式的结果值相等，那么位置靠上的case子句总会被选中

总结
	range表达式的结果值是会被复制的，实际迭代时并不会使用原值
	至于会影响到什么，那就要看这个结果值的类型是值类型还是引用类型了

思考
	1.在类型switch语句中，我们怎样对被判断类型的那个值做相应的类型转换？
	2.在if语句中，初始化子句声明的变量的作用域是什么？

补充
	在 switch x.(type) 语句中，无法使用 fallthrough
*/

// RangeArray range 表达式的求值结果会被复制
func RangeArray() {
	numbers2 := [...]int{1, 2, 3, 4, 5, 6} // [7 3 5 7 9 11]
	//numbers2 := []int{1, 2, 3, 4, 5, 6} // [22 3 6 10 15 21]
	maxIndex2 := len(numbers2) - 1
	for i, e := range numbers2 { // 被迭代的数组与numbers2已经是毫不相关的两个数组了
		if i == maxIndex2 {
			numbers2[0] += e
		} else {
			numbers2[i+1] += e
		}
	}
	fmt.Println(numbers2)

	//arr := [...]int{1, 2, 3, 4, 5, 6}
	//for i, v := range arr {
	//	if i&1 == 0 {
	//		arr[i+1] += 10 // 注意越界
	//	}
	//	fmt.Println(v)
	//}
	//fmt.Println(arr)
}

// SwitchCase 绕过case子表达式结果值相等的限制
func SwitchCase() {
	value5 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value5[4] {
	case value5[0], value5[1], value5[2]:
		fmt.Println("0 or 1 or 2")
	//case value5[4], value5[5], value5[6]: // 选择在前面的
	//	fmt.Println("4 or 5 or 6")
	case value5[2], value5[3], value5[4]: // 选择在前面的
		fmt.Println("2 or 3 or 4")
	case value5[4], value5[5], value5[6]:
		fmt.Println("4 or 5 or 6")
	}
}
