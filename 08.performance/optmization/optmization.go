package optmization

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

/*
性能调优过程

	S->设定优化目标->分析系统瓶颈点->优化瓶颈点->E
		   ↑_______________________丨

常见分析指标

	Wall Time：运行的绝对时间(包括阻塞、等待外部响应...)
	CPU Time：CPU消耗时间
	Block Time：
	Memory allocation：内存分配
	GC times/time spent：GC次数和耗时

示例

	fori VS forr
		benchmark 测试证明，两种写法性能上基本没有差别
	easyjson VS API
		easyjson 更快
	bytes.Buffer VS +
		这里测试，性能差不多
	make([]T,0,n) VS make([]T,n,n)
		make([]T,n,n) 高很多，为什么？
		性能差异：
			make([]T,0,n) VS make([]T,n,n) >> bytes.Buffer VS +

		precessRequestEasyJson：两种写法，性能差 20 倍
		为什么

测试：实际开发中，单个 Test 可以简化 prof 文件的大小，便于分析

	go test -bench="." --cpuprofile=optmization.prof
	go tool pprof optmization.prof
		top
		top -cum
		list BenchmarkProcessRequestEasyJson
		exit
	go-torch optmization.prof
*/
func createRequest() string {
	payLoad := make([]int, 100, 100)
	for i := 1; i < 100; i++ {
		payLoad[i] = i
	}
	req := Request{"demo_transaction", payLoad}
	v, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	return string(v)
}
func processRequest(reqs []string) []string {
	reps := make([]string, 0, len(reqs))
	for _, req := range reqs {
		reqObj := &Request{}
		reqObj.UnmarshalJSON([]byte(req))
		//	json.Unmarshal([]byte(req), reqObj)
		var buf strings.Builder
		for _, e := range reqObj.PayLoad {
			buf.WriteString(strconv.Itoa(e))
			buf.WriteString(",")
		}
		repObj := &Response{reqObj.TransactionID, buf.String()}
		repJson, err := repObj.MarshalJSON()
		//repJson, err := json.Marshal(&repObj)
		if err != nil {
			panic(err)
		}
		reps = append(reps, string(repJson))
	}
	return reps
}

func precessRequestEasyJson(reqs []string) []string {
	//reps := make([]string, 0, len(reqs))
	reps := make([]string, len(reqs), len(reqs))
	//var buf bytes.Buffer
	//for _, req := range reqs {
	for i, req := range reqs {
		reqObj := &Request{}
		reqObj.UnmarshalJSON([]byte(req))
		var buf bytes.Buffer
		for _, v := range reqObj.PayLoad {
			buf.WriteString(strconv.Itoa(v))
			buf.WriteRune(',')
		}
		repObj := &Response{reqObj.TransactionID, buf.String()}
		repJson, err := repObj.MarshalJSON()
		if err != nil {
			panic(err)
		}
		reqs[i] = string(repJson)
		//reps = append(reps, string(repJson))
		//buf.Reset()	// 复用 buf，并没有提高性能
	}
	return reps
}

func processRequestAPI(reqs []string) []string {
	n := len(reqs)
	reps := make([]string, n)
	//var sb strings.Builder
	for i := 0; i < n; i++ {
		reqObj := &Request{}
		json.Unmarshal([]byte(reqs[i]), reqObj)
		var s string
		for _, v := range reqObj.PayLoad {
			s += strconv.Itoa(v) + ","
		}
		repObj := &Response{reqObj.TransactionID, s}
		repJson, err := json.Marshal(repObj)
		if err != nil {
			panic(err)
		}
		reqs[i] = string(repJson)
	}
	return reps
}
func processRequestAPI_1(reqs []string) []string {
	n := len(reqs)
	reps := make([]string, n)
	//var sb strings.Builder
	for i := 0; i < n; i++ {
		reqObj := &Request{}
		reqObj.UnmarshalJSON([]byte(reqs[i]))
		var s string
		for _, v := range reqObj.PayLoad {
			s += strconv.Itoa(v) + ","
		}
		repObj := &Response{reqObj.TransactionID, s}
		repJson, err := repObj.MarshalJSON()
		if err != nil {
			panic(err)
		}
		reqs[i] = string(repJson)
	}
	return reps
}

func processRequestAPI_Range(reqs []string) []string {
	reps := make([]string, len(reqs))
	//var sb strings.Builder
	for i, req := range reqs {
		reqObj := &Request{}
		json.Unmarshal([]byte(req), reqObj)
		var s string
		for _, v := range reqObj.PayLoad {
			s += strconv.Itoa(v) + ","
		}
		repObj := &Response{reqObj.TransactionID, s}
		repJson, err := json.Marshal(repObj)
		if err != nil {
			panic(err)
		}
		reqs[i] = string(repJson)
	}
	return reps
}
