package http

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"testing"
)

/*
httprouter
	高性能 Router
	github：https://github.com/julienschmidt/httprouter.git
	$ go get -u github.com/julienschmidt/httprouter

	Restful 比较强调采用 http 原有的基本动作：Get Post Put Delete ...
		httprouter 支持了这点
原理：prefix tree or radix tree(基数树)
	通过压缩 Trie 树，实现高效的路径匹配
gin
	路由模块：httprouter
	json：jsoniter
	github：https://github.com/gin-gonic/gin
	社区middlewares：https://github.com/gin-gonic/contrib
相关测试
	https://github.com/julienschmidt/go-http-routing-benchmark
*/

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome!\n")
}
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s...\n", ps.ByName("name"))
}
func TestHttpRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello) // ps.ByName("name")
	log.Fatal(http.ListenAndServe(":8080", router))
}
