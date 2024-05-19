package generic

import (
	"golang.org/x/exp/constraints"
	"sort"
)

type plusable interface {
	~int | ~string
}

func Add[T plusable](a, b T) T {
	return a + b
}

func AddInt(a, b int) int {
	return a + b
}
func AddString(a, b string) string {
	return a + b
}

// 定义支持排序的泛型切片
type sortableSlice[T constraints.Ordered] []T

func (s sortableSlice[T]) Len() int           { return len(s) } // 让泛型切片实现 sort.Interface
func (s sortableSlice[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s sortableSlice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// SortGeneric 定义一个泛型排序函数
func SortGeneric[T constraints.Ordered](sl sortableSlice[T]) {
	sort.Sort(sl)
}
