package advance

import (
	"fmt"
	"reflect"
)

/*
接口类型的合理运用

接口类型
	接口类型与其他数据类型不同，它是没法被值化的，或者说是没法被实例化的
	更具体地说，我们既不能通过调用new函数或make函数创建出一个接口类型的值，也无法用字面量来表示一个接口类型的值
	对于某一个接口类型来说，如果没有任何数据类型可以作为它的实现，那么该接口的值就不可能存在
	接口类型声明中的方法所代表的就是该接口的方法集合。一个接口的方法集合就是它的全部特征
接口的实现类型
	对于任何数据类型，只要它的方法集合中完全包含了一个接口的全部特征（即全部的方法），那么它就一定是这个接口的实现类型
	这是一种无侵入式的接口实现方式，专有名词叫“Duck typing”
判定一个数据类型的某一个方法实现的就是某个接口类型中的某个方法
	两个充分必要条件，一个是“两个方法的签名需要完全一致”，另一个是“两个方法的名称要一模一样”
示例
	接口 Pet，实现 Dog
		dog := Dog{"little pig"}
		var pet Pet = &dog
	&dog 是变量 pet 的实际值（或动态值）
	*Dog 是 pet 变量的实际类型（或动态类型）
	对于变量pet来讲，它的静态类型就是Pet，并且永远是Pet，但是它的动态类型却会随着我们赋给它的动态值而变化
	在我们给一个接口类型的变量赋予实际的值之前，它的动态类型是不存在的

问题：当我们为一个接口变量赋值时会发生什么？
	测试一：InterfaceTest
	pet变量的字段name的值依然是"little pig"
问题解析
	赋值规则
		使用一个变量给另外一个变量赋值，那么真正赋给后者的，并不是前者持有的那个值，而是该值的一个副本
		测试二
	接口赋值
		零值
			接口类型本身是无法被值化的。在我们赋予它实际的值之前，它的值一定会是nil，这也是它的零值
			一旦它被赋予了某个实现类型的值，它的值就不再是nil了
		当我们给一个接口变量赋值的时候，该变量的动态类型会与它的动态值一起被存储在一个专用的数据结构中
		严格来讲，这样一个变量（如 dog）的值其实是这个专用数据结构的一个实例，而不是我们赋给该变量的那个实际的值
		无论是从它们存储的内容，还是存储的结构上，与“赋值规则”是不同的
		在 Go 语言的runtime包中，这个专用的数据结构叫做 iface
	iface
		iface的实例会包含两个指针，一个是指向类型信息的指针，另一个是指向动态值的指针
		总之，接口变量被赋予动态值的时候，存储的是包含了这个动态值的副本的一个结构更加复杂的值

知识扩展
问题 1：接口变量的值在什么情况下才真正为nil？
	示例：InterfaceTestNil
		虽然被包装（iface 的“值” = dog2）的动态值是nil，但是pet的值却不会是nil，因为这个动态值只是pet值的一部分而已
		这时的pet的动态类型就存在了，是*Dog
		这时 Go 语言会识别出赋予pet的值是一个*Dog类型的nil，此时 pet != 字面量nil
		然后，Go 语言就会用一个iface的实例包装它，包装后的产物肯定就不是nil了
	nil
		在 Go 语言中，我们把由字面量nil表示的值叫做无类型的nil。这是真正的nil，因为它的类型也是nil的
	panic
		此时调用 pet.Name()，会空指针
		panic: value method golang/core/advance.Dog.Name called using nil *Dog pointer
问题 2：怎样实现接口之间的组合？
	接口的组合
		只要组合的接口之间有同名的方法就会产生冲突，从而无法通过编译，即使同名方法的签名彼此不同也会是如此
		因此，接口的组合根本不可能导致“屏蔽”现象的出现
	较小的接口
		Go 语言团队鼓励我们声明体量较小的接口，并建议我们通过这种接口间的组合来扩展程序、增加程序的灵活性
		因为相比于包含很多方法的大接口而言，小接口可以更加专注地表达某一种能力或某一类特征，同时也更容易被组合在一起
	Go 语言标准库代码包io中的ReadWriteCloser接口和ReadWriter接口就是这样的例子，它们都是由若干个小接口组合而成的
		以io.ReadWriteCloser接口为例，它是由io.Reader、io.Writer和io.Closer这三个接口组成的
		这三个接口都只包含了一个方法，是典型的小接口
		它们中的每一个都只代表了一种能力，分别是读出、写入和关闭
		即使我们只实现了io.Reader和io.Writer，那么也等同于实现了io.ReadWriter接口，因为后者就是前两个接口组成的
		这几个io包中的接口共同组成了一个接口矩阵。它们既相互关联又独立存在
	善用接口组合和小接口可以让你的程序框架更加稳定和灵活

总结
	当我们给接口变量赋值时，接口变量会持有被赋予值的副本，而不是它本身
	接口变量的值并不等同于这个可被称为动态值的副本。它会包含两个指针，一个指针指向动态值，一个指针指向类型信息
	除非我们只声明而不初始化，或者显式地赋给它nil，否则接口变量的值就不会为nil

思考
	如果我们把一个值为nil的某个实现类型的变量赋给了接口变量，那么在这个接口变量上仍然可以调用该接口的方法吗？
	如果可以，有哪些注意事项？如果不可以，原因是什么？

补充
	接口变量使用的数据结构是iface，那引用类型使用的数据结构又是什么呢？比如slice
	查看 runtime 包中，有个叫slice的结构体类型
*/

type Pet interface {
	Name() string
	Category() string
}
type Dog struct{ name string }

func (d Dog) Name() string     { return d.name }
func (d Dog) Category() string { return "dog" }

// setName 如果是值方法（结构体接收器），那么 4 个打印都是 dog wang!
func (d *Dog) setName(name string) { d.name = name }

// InterfaceTest 值
func InterfaceTest() {
	dog := Dog{"wang!"}
	// 测试一
	var pet Pet = dog
	dog.setName("miao~")
	fmt.Println("dog: ", dog.Category(), dog.Name()) // dog:  dog miao~
	fmt.Println("pet: ", pet.Category(), pet.Name()) // pet:  dog wang!

	// 测试二
	dog1 := dog
	dog.setName("mmm?")
	fmt.Println("dog: ", dog.Category(), dog.Name())    // dog:  dog mmm?
	fmt.Println("dog1: ", dog1.Category(), dog1.Name()) // dog1:  dog miao~

	// 测试三
	var pet1 Pet = &dog
	dog.setName("m...iao...")
	fmt.Println("dog: ", dog.Category(), dog.Name())    // dog:  dog m...iao...
	fmt.Println("pet1: ", pet1.Category(), pet1.Name()) // pet1:  dog m...iao...
}

// InterfaceTestNil 测试接口 nil
func InterfaceTestNil() {
	// 测试四
	var dog1 *Dog
	fmt.Println(dog1) // <nil>
	dog2 := dog1
	fmt.Println(dog2) // <nil>
	var pet Pet = dog2
	if pet == nil {
		fmt.Println("The pet is nil.")
	} else {
		fmt.Println("The pet is not nil.") // The pet is not nil.
	}
	fmt.Printf("The type of pet is %T.\n", pet)                          // The type of pet is *advance.Dog.
	fmt.Printf("The type of pet is %s.\n", reflect.TypeOf(pet).String()) // The type of pet is *advance.Dog.
	fmt.Printf("The type of second dog is %T.\n", dog2)                  // The type of second dog is *advance.Dog.
	//pet.Category()
	//fmt.Println(pet.Name())                                              // panic: value method golang/core/advance.Dog.Name called using nil *Dog pointer

	// 测试五
	wrap := func(dog *Dog) Pet {
		if dog == nil {
			return nil
		}
		return dog
	}
	pet = wrap(dog2)
	if pet == nil {
		fmt.Println("The pet is nil.") // The pet is nil.
	} else {
		fmt.Println("The pet is not nil.")
	}
}
