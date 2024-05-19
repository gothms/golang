package json

import (
	"encoding/json"
	"testing"
)

var jsonStr = `{
	"basic_info":{
		"name":"Mike",
		"age":18
	},
	"job_info":{
		"skills":["Java","Go","C"]
	}
}`

func TestJson(t *testing.T) {
	e := new(Employee)
	err := json.Unmarshal([]byte(jsonStr), e)
	if err != nil {
		t.Error(err)
	}
	t.Log(*e)
	if v, err := json.Marshal(e); e != nil {
		t.Log(string(v))
	} else {
		t.Error(err)
	}
}
