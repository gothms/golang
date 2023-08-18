package performance

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

/*
通过 文件 方式输出 Profile

	适合短时间批量运行/细粒度的程序的分析

通过 HTTP 方式输出 Profile

	1.简单，适用于持续性运行的应用：服务端程序
	2.在应⽤用程序中导入 import _ "net/http/pprof"，并启动 http server 即可
		重要 import _ "net/http/pprof"

访问方式：http://127.0.0.1:8080/fb

	1.http://<host>:<port>/debug/pprof/
		http://127.0.0.1:8080/debug/pprof/
		点击 profile 会把 profile 文件下载到本地
	2.go tool pprof http://<host>:<port>/debug/pprof/profile?seconds=10 （默认值为30秒）
		$ go tool pprof http://127.0.0.1:8080/debug/pprof/profile?seconds=10
	3.go-torch -seconds 10 http://<host>:<port>/debug/pprof/profile
		$ go-torch http://127.0.0.1:8080/debug/pprof/profile
		保存在了上次的目录
*/
func TestHttpProf(t *testing.T) {
	http.HandleFunc("/", index)
	http.HandleFunc("/fb", createFBS) // http://127.0.0.1:8080/fb
	t.Fatal(http.ListenAndServe(":8080", nil))
}
func GetFibonacciSerie(n int) []int {
	ret := make([]int, 2, n)
	ret[0], ret[1] = 1, 1
	for i := 2; i < n; i++ {
		ret = append(ret, ret[i-2]+ret[i-1])
	}
	return ret
}
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}
func createFBS(w http.ResponseWriter, r *http.Request) {
	var fbs []int
	for i := 0; i < 1000000; i++ {
		fbs = GetFibonacciSerie(50)
	}
	w.Write([]byte(fmt.Sprintf("%v", fbs)))
}
