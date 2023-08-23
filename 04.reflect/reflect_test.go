package unsafe_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

/*
reflect.TypeOf vs reflect.ValueOf
	1.reflect.TypeOf 返回类型 (reflect.Type)
	2.reflect.ValueOf 返回值 (reflect.Value)
	3.可以从 reflect.Value 获得类型
	4.通过 kind 来判断类型（27种）
		const (
			Invalid Kind = iota
			Bool
			Int
			Int8
			...
		)

利用反射编写灵活的代码
	按名字访问结构的成员：e 为结构体
		reflect.ValueOf(e).FieldByName("Name")
	按名字访问结构的方法：e 为结构体指针
		reflect.ValueOf(e).MethodByName("UpdateAge").
			Call([]reflect.Value{reflect.ValueOf(1)})

struct Tag
	key value 结构
	如：
		Name string `json:"name"`
		Name string `format:"normal"`
区别：
	reflect.ValueOf：FieldByName
	reflect.TypeOf：FieldByName
万能程序
	优点
		与配置相关的程序，提高灵活性
		复用性高，万能程序
	缺点
		可读性差
		debug困难
		性能大大降低
	Elem()：获取指针指向的值
*/

// 万能程序：反序列化

func TestFillField(t *testing.T) {
	set := map[string]interface{}{"Name": "Mike", "Age": 18}
	e := Employee{}
	if err := fillField(&e, set); err != nil {
		t.Fatal(err)
	}
	t.Log(e)
	c := new(Customer)
	if err := fillField(c, set); err != nil {
		t.Fatal(err)
	}
	t.Log(*c)
}

func fillField(s interface{}, set map[string]interface{}) error {
	if set == nil {
		return errors.New("set is nil")
	}
	//if reflect.TypeOf(s).Kind() != reflect.Ptr {	// Ptr 是老版本的名称
	if reflect.TypeOf(s).Kind() != reflect.Pointer { // 不是指针
		if reflect.TypeOf(s).Elem().Kind() != reflect.Struct { // 不是结构体
			return errors.New(`the first param should be a pointer to the struct type`)
		}
	}
	var (
		field reflect.StructField
		ok    bool
	)

	for k, v := range set {
		if field, ok = reflect.ValueOf(s).Elem().Type().FieldByName(k); !ok {
			continue
		} // 没有 k 这个字段
		if field.Type == reflect.TypeOf(v) { // 字段值的类型相同
			vstr := reflect.ValueOf(s)
			vstr = vstr.Elem()                          // 根据指针获取结构
			vstr.FieldByName(k).Set(reflect.ValueOf(v)) // 设置值
		}
	}
	return nil
}

// 反射编写灵活代码
type Employee struct {
	EId  string
	Name string `format:"normal"`
	Age  int
}

func (e *Employee) UpdateAge(age int) {
	e.Age = age
}

type Customer struct {
	CId  string
	Name string
	Age  int
}

func TestInvokeByName(t *testing.T) {
	e := Employee{"1", "Lee", 25}
	t.Logf("Name: value(%[1]v),type(%[1]T)", reflect.ValueOf(e).FieldByName("Name"))
	if name, ok := reflect.TypeOf(e).FieldByName("Name"); ok {
		t.Log("TypeOf:name", name)                   // TypeOf:name {Name  string format:"normal" 16 [1] false}
		t.Log("Tag->format", name.Tag.Get("format")) // struct Tag: Tag->format normal
	} else {
		t.Error(`failed to get "Name" field`)
	}
	// 反射修改 Age
	reflect.ValueOf(&e).MethodByName("UpdateAge").
		Call([]reflect.Value{reflect.ValueOf(1)})
	t.Log("Updated Age:", e)
}

// reflect.TypeOf vs reflect.ValueOf
func TestTypeAndValue(t *testing.T) {
	var f int64 = 10
	t.Log(reflect.TypeOf(f), reflect.ValueOf(f))
	t.Log(reflect.ValueOf(f).Type())
}
func TestBasicType(t *testing.T) {
	var f float64 = 3.2
	CheckType(f)
	CheckType(&f) // *float64
}
func CheckType(v interface{}) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		fmt.Println("Float")
	case reflect.Int, reflect.Int32, reflect.Int64:
		fmt.Println("Int")
	default:
		fmt.Println("Unknown", t)
	}
}
