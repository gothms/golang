package modes

/*
Go 编程模式：Map-Reduce

Map、Reduce、Filter
	函数式编程中非常重要的 Map、Reduce、Filter 三种操作
	可以让我们轻松灵活地进行一些数据处理，毕竟，我们的程序大多数情况下都在倒腾数据
	尤其是对于一些需要统计的业务场景来说，Map、Reduce、Filter 是非常通用的玩法
	关键：函数作为参数，泛型，反射

基本示例
	Map
	Reduce
	Filter
业务示例
	分离 控制逻辑 和 业务逻辑
		Map、Reduce、Filter 只是一种控制逻辑，真正的业务逻辑是以传给它们的数据和函数来定义的
		这是一个很经典的“业务逻辑”和“控制逻辑”分离解耦的编程模式
	示例代码如下

泛型 Map-Reduce
	用 interface{} + reflect来完成
	interface{} 可以理解为 C 中的 void*、Java 中的 Object ，reflect是 Go 的反射机制包，作用是在运行时检查类型
		func Map(data interface{}, fn interface{}) []interface{} {
			vfn := reflect.ValueOf(fn)
			vdata := reflect.ValueOf(data)
			result := make([]interface{}, vdata.Len())

			for i := 0; i < vdata.Len(); i++ {
				result[i] = vfn.Call([]reflect.Value{vdata.Index(i)})[0].Interface()
			}
			return result
		}
		首先，我们通过 reflect.ValueOf() 获得 interface{} 的值，其中一个是数据 vdata，另一个是函数 vfn
		然后，通过 vfn.Call() 方法调用函数，通过 []refelct.Value{vdata.Index(i)}获得数据
	panic
		因为反射是运行时的事，所以，如果类型出问题的话，就会有运行时的错误
		代码可以很轻松地编译通过，但是在运行时却出问题了，而且还是 panic 错误
健壮版的 Generic Map
	示例代码
		func Transform(slice, function interface{}) interface{} {
		  return transform(slice, function, false)
		}

		func TransformInPlace(slice, function interface{}) interface{} {
		  return transform(slice, function, true)
		}

		func transform(slice, function interface{}, inPlace bool) interface{} {

		  //check the `slice` type is Slice
		  sliceInType := reflect.ValueOf(slice)
		  if sliceInType.Kind() != reflect.Slice {
			panic("transform: not slice")
		  }

		  //check the function signature
		  fn := reflect.ValueOf(function)
		  elemType := sliceInType.Type().Elem()
		  if !verifyFuncSignature(fn, elemType, nil) {
			panic("trasform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
		  }

		  sliceOutType := sliceInType
		  if !inPlace {
			sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
		  }
		  for i := 0; i < sliceInType.Len(); i++ {
			sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
		  }
		  return sliceOutType.Interface()

		}

		func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {

		  //Check it is a funciton
		  if fn.Kind() != reflect.Func {
			return false
		  }
		  // NumIn() - returns a function type's input parameter count.
		  // NumOut() - returns a function type's output parameter count.
		  if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
			return false
		  }
		  // In() - returns the type of a function type's i'th input parameter.
		  for i := 0; i < len(types)-1; i++ {
			if fn.Type().In(i) != types[i] {
			  return false
			}
		  }
		  // Out() - returns the type of a function type's i'th output parameter.
		  outType := types[len(types)-1]
		  if outType != nil && fn.Type().Out(0) != outType {
			return false
		  }
		  return true
		}
	代码中的几个要点
		代码中没有使用 Map 函数，因为和数据结构有含义冲突的问题，所以使用Transform，这个来源于 C++ STL 库中的命名
		有两个版本的函数，一个是返回一个全新的数组 Transform()，一个是“就地完成” TransformInPlace()
		在主函数中，用 Kind() 方法检查了数据类型是不是 Slice，函数类型是不是 Func
		检查函数的参数和返回类型是通过 verifyFuncSignature() 来完成的：NumIn()用来检查函数的“入参”；NumOut() ：用来检查函数的“返回值”
		如果需要新生成一个 Slice，会使用 reflect.MakeSlice() 来完成
健壮版的 Generic Reduce
	func Reduce(slice, pairFunc, zero interface{}) interface{} {
	  sliceInType := reflect.ValueOf(slice)
	  if sliceInType.Kind() != reflect.Slice {
		panic("reduce: wrong type, not slice")
	  }

	  len := sliceInType.Len()
	  if len == 0 {
		return zero
	  } else if len == 1 {
		return sliceInType.Index(0)
	  }

	  elemType := sliceInType.Type().Elem()
	  fn := reflect.ValueOf(pairFunc)
	  if !verifyFuncSignature(fn, elemType, elemType, elemType) {
		t := elemType.String()
		panic("reduce: function must be of type func(" + t + ", " + t + ") " + t)
	  }

	  var ins [2]reflect.Value
	  ins[0] = sliceInType.Index(0)
	  ins[1] = sliceInType.Index(1)
	  out := fn.Call(ins[:])[0]

	  for i := 2; i < len; i++ {
		ins[0] = out
		ins[1] = sliceInType.Index(i)
		out = fn.Call(ins[:])[0]
	  }
	  return out.Interface()
	}
健壮版的 Generic Filter
	func Filter(slice, function interface{}) interface{} {
	  result, _ := filter(slice, function, false)
	  return result
	}

	func FilterInPlace(slicePtr, function interface{}) {
	  in := reflect.ValueOf(slicePtr)
	  if in.Kind() != reflect.Ptr {
		panic("FilterInPlace: wrong type, " +
		  "not a pointer to slice")
	  }
	  _, n := filter(in.Elem().Interface(), function, true)
	  in.Elem().SetLen(n)
	}

	var boolType = reflect.ValueOf(true).Type()

	func filter(slice, function interface{}, inPlace bool) (interface{}, int) {

	  sliceInType := reflect.ValueOf(slice)
	  if sliceInType.Kind() != reflect.Slice {
		panic("filter: wrong type, not a slice")
	  }

	  fn := reflect.ValueOf(function)
	  elemType := sliceInType.Type().Elem()
	  if !verifyFuncSignature(fn, elemType, boolType) {
		panic("filter: function must be of type func(" + elemType.String() + ") bool")
	  }

	  var which []int
	  for i := 0; i < sliceInType.Len(); i++ {
		if fn.Call([]reflect.Value{sliceInType.Index(i)})[0].Bool() {
		  which = append(which, i)
		}
	  }

	  out := sliceInType

	  if !inPlace {
		out = reflect.MakeSlice(sliceInType.Type(), len(which), len(which))
	  }
	  for i := range which {
		out.Index(i).Set(sliceInType.Index(which[i]))
	  }

	  return out.Interface(), len(which)
	}

未尽事宜
	用反射来实现泛型，代码的性能会很差，不能用在需要高性能的地方
	代码大量地参考了 Rob Pike 的版本：https://github.com/robpike/filter
	什么时候在标准库中支持 Map、Reduce？Rob Pike 说到：我一次都没有用过，我还是喜欢用“For 循环”，我觉得你最好也跟我一起用 “For 循环”
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
