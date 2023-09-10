package main

import (
	"fmt"
	"math"
	"math/bits"
)

func main() {
	f := math.NaN()
	fmt.Println(f == f)

	type keyType struct{ v float64 }
	key := keyType{math.NaN()}
	h := make(map[keyType]int, 2)
	h[key] = 123
	v := h[key]
	fmt.Println(v)

	var state int32 = 1
	fmt.Println(state << 30)
	fmt.Printf("%b, %d\n", state<<30, bits.Len32(uint32(state<<30)))
	fmt.Println(state<<31 - 1)
	fmt.Println((state<<31 - 1) >> 3)
	fmt.Println(state << 32)
}
