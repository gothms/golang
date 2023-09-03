package p2

import (
	"fmt"
	gomem "golang/concurrent/gomem"
	"golang/concurrent/gomem/p3"
)

var V1_p2 = gomem.Trace("init v1_p2", 2)
var V2_p2 = gomem.Trace("init v2_p2", p3.V2_p3)

func init() {
	fmt.Println("init func in p2")
	V1_p2 = 200
}
