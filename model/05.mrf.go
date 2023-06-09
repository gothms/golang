package model

/*
	函数式编程：https://time.geekbang.com/column/article/330232?utm_source=pc_cp&utm_term=pc_interstitial_1346
	关键：函数作为参数，泛型，反射

	1.Map
	2.Reduce
	3.Filter
	4.分离 控制逻辑 和 业务逻辑
		Map、Reduce、Filter 只是一种控制逻辑，真正的业务逻辑是传给它们的数据和函数定义的
	5.泛型 Map-Reduce
		interface{}：Object
		反射：
	6.用反射来实现泛型，代码的性能会很差，不能用在需要高性能的地方
		参考：https://github.com/robpike/filter
*/
type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

// EmployeeCountIf Filter + Reduce 语义
func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}
	return count
}

// EmployeeFilterIn Filter 语义
func EmployeeFilterIn(list []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i, _ := range list {
		if fn(&list[i]) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

// EmployeeSumIf Filter + Reduce 语义
func EmployeeSumIf(list []Employee, fn func(e *Employee) int) int {
	var sum = 0
	for i, _ := range list {
		sum += fn(&list[i])
	}
	return sum
}
