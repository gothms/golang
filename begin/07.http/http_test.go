package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

/*
Default Router

	api:http.Handler.ServeHTTP()
	路径 net/http/server.go
	func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
		handler := sh.srv.Handler
		if handler == nil {
			handler = DefaultServeMux	// 使用缺省的Router
		}
		if !sh.srv.DisableGeneralOptionsHandler && req.RequestURI == "*" && req.Method == "OPTIONS" {
			handler = globalOptionsHandler{}
		}
		...
		handler.ServeHTTP(rw, req)
	}

路由规则

	1.URL分为两种，末尾是 / 表示一个子树，后面可以跟其他子路径；末尾不是 / 表示一个叶子，固定的路径
		以 / 结尾的 URL 可以匹配它的任何子路径，比如 /images 会匹配 /images/cute-cat.jpg

		/time/a.go：只能匹配 /time/a.go
		/time/：匹配 /time 和 /time/ + 任意
	2.采用最长匹配原则，如果有多个匹配，一定采用匹配路径最长的那个进行处理
	3.如果没有找到任何匹配项，会返回 404 错误
*/
func TestName(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello HTTP!")
	})
	//http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
	//http.HandleFunc("/time/", func(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/time/a.go", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		timeStr := fmt.Sprintf(`"time":%s`, t)
		w.Write([]byte(timeStr))
	})
	http.ListenAndServe(":8080", nil)
}
