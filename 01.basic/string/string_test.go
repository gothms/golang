package string

import (
	"strconv"
	"strings"
	"testing"
)

func TestStringFn(t *testing.T) {
	s := "A,B,C"
	parts := strings.Split(s, ",")
	for _, p := range parts {
		t.Log(p)
	}
	t.Log(strings.Join(parts, "-"))
}
func TestConv(t *testing.T) {
	s := strconv.Itoa(10)
	t.Logf("%T", s)
	if v, err := strconv.Atoi(s); err == nil {
		t.Logf("%[1]T %[1]d", v+3)
	}
}
