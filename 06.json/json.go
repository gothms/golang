package json

/*
api：简单而少量解析json，推荐使用api
	简单
	易于使用
	不适用于高性能环境
原理：参考 reflect_test.go 万能程序
	通过反射 reflect 实现
	通过 FieldTag 标识对应的 json 值
easyjson VS jsoniter
	https://www.libhunt.com/compare-easyjson-vs-json-iterator--go
*/

type BasicInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type JobInfo struct {
	Skills []string `json:"skills"`
}
type Employee struct {
	BasicInfo BasicInfo `json:"basic_info"`
	JobInfo   JobInfo   `json:"job_info"`
}
