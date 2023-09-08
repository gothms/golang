package main

import (
	"fmt"
	"math"
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
}
