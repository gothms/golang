package api

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

/*
bug
	test 执行结果是：{“a”:6673221165400540000}
	原始数据是：{“a”:6673221165400540161}
	依赖：jsoniter "github.com/json-iterator/go"
解决
	solution
	使用了 func (*Decoder) UseNumber 方法告诉反序列化 json 的数字类型的时候，不要直接转换成 float64
	而是转换成 json.Number 类型
	json.Number 本质是字符串，反序列化的时候将 json 的数值先转成 json.Number，其实是一种延迟处理的手段
	待后续逻辑需要时候，再把 json.Number 转成 float64 或者 int64
*/

func test() {
	s := "{\"a\":6673221165400540161}"

	d := make(map[string]interface{})
	err := jsoniter.Unmarshal([]byte(s), &d)
	if err != nil {
		panic(err)
	}

	s2, err := jsoniter.Marshal(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s2))
}
func solution() {
	s := "{\"a\":6673221165400540161}"
	decoder := jsoniter.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	d := make(map[string]interface{})
	err := decoder.Decode(&d)
	if err != nil {
		panic(err)
	}

	s2, err := jsoniter.Marshal(d)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s2))
}
