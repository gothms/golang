package test

import (
	"strconv"
	"strings"
	"testing"
)

const numbers = 100

func BenchmarkStringBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.Grow(numbers)
		for j := 0; j < numbers; j++ {
			sb.WriteString(strconv.Itoa(j))
		}
		_ = sb.String()
	}
	b.StopTimer()
}
func BenchmarkStringAdd(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var s string
		for j := 0; j < numbers; j++ {
			s += strconv.Itoa(j)
		}
	}
	b.StopTimer()
}
