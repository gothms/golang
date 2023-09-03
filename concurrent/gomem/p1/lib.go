package p1

import (
	"fmt"
	gomem "golang/concurrent/gomem"
	"golang/concurrent/gomem/p2"
)

var V1_p1 = gomem.Trace("init v1_p1", p2.V1_p2)
var V2_p1 = gomem.Trace("init v2_p1", p2.V2_p2)

func init() {
	fmt.Println("init func in p1")
}
