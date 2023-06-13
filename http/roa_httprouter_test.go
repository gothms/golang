package http

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"testing"
)

/*
Restful: Resource Oriented Architecture

	面向资源的架构
	书籍 ROA：RESTful Web Services
*/
type Employee struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var employeeDB map[string]*Employee

func init() {
	employeeDB = map[string]*Employee{}
	employeeDB["Mike"] = &Employee{"e_01", "Mike", 19}
	employeeDB["Rose"] = &Employee{"e_02", "Rose", 19}
	employeeDB["Lee"] = &Employee{"e_03", "Lee", 16}
}
func ROAIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome!\n")
}
func GetEmployeeByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qName := ps.ByName("name")
	var (
		ok       bool
		info     *Employee
		infoJson []byte
		err      error
	)
	if info, ok = employeeDB[qName]; !ok {
		w.Write([]byte(`{"error": "Not Found"}`))
		return
	}
	if infoJson, err = json.Marshal(info); err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": %s`, err)))
		return
	}
	w.Write(infoJson)
}
func TestROAHttpRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/employees/:name", GetEmployeeByName) // ps.ByName("name")
	log.Fatal(http.ListenAndServe(":8080", router))
}
