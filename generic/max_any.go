package generic

import "golang.org/x/exp/constraints"

func MaxInt(sl []int) int {
	if len(sl) == 0 {
		panic("slice is empty")
	}

	max := sl[0]
	for _, v := range sl[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func MaxAny(sl []any) any {
	if len(sl) == 0 {
		panic("slice is empty")
	}

	max := sl[0]
	for _, v := range sl[1:] {
		switch v.(type) {
		case int:
			if v.(int) > max.(int) {
				max = v
			}
		case string:
			if v.(string) > max.(string) {
				max = v
			}
		case float64:
			if v.(float64) > max.(float64) {
				max = v
			}
		}
	}
	return max
}
func MaxGenerics[T constraints.Ordered](sl []T) T {
	if len(sl) == 0 {
		panic("slice is empty")
	}
	ans := sl[0]
	for _, v := range sl[1:] {
		//if v > ans {
		//	ans = v
		//}
		ans = max(ans, v)
	}
	return ans
}
