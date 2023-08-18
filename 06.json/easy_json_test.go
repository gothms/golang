package json

import (
	"encoding/json"
	"testing"
)

/*
EasyJSON 采用代码生成而非反射
	性能数一数二
		更高效
		内存开销更小

安装：https://github.com/mailru/easyjson.git

	$ go get -u github.com/mailru/easyjson/...
	$ go install github.com/mailru/easyjson

使用：$ easyjson -all <file>.go

	<file>.go 为结构体的go文件
	针对结构而生成了代码，在 pkg 下并包含了实现(未使用反射)：
		Marshal
		Unmarshal
	使用
		struct{}.UnmarshalJSON()
		struct{}.MarshalJSON()
	测试
		go test -bench="." -benchmem
*/

var jsons = `{
	"basic_info":{
		"name":"Mike",
		"age":18
	},
	"job_info":{
		"skills":["Java","Go","C"]
	}
}`

func TestEasyJson(t *testing.T) {
	e := new(Employee)
	e.UnmarshalJSON([]byte(jsons))
	t.Log(e)
	if v, err := e.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		t.Log(string(v))
	}
}
func BenchmarkAPIJson(b *testing.B) {
	b.ResetTimer()
	e := new(Employee)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsons), e)
		if err != nil {
			b.Error(err)
		}
		if _, err = json.Marshal(e); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}
func BenchmarkEasyJson(b *testing.B) {
	b.ResetTimer()
	e := new(Employee)
	for i := 0; i < b.N; i++ {
		err := e.UnmarshalJSON([]byte(jsons))
		if err != nil {
			b.Error(err)
		}
		if _, err = e.MarshalJSON(); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}
