package modes

import "errors"

// IntSet 业务逻辑
type IntSet struct {
	data map[int]bool
	undo Undo
}

func NewIntSet() IntSet {
	return IntSet{data: make(map[int]bool)}
}
func (set *IntSet) Undo() error {
	return set.undo.Undo()
}
func (set *IntSet) Contains(v int) bool {
	return set.data[v]
}
func (set *IntSet) Add(v int) {
	if !set.Contains(v) {
		set.data[v] = true
		set.undo.Add(func() {
			set.Delete(v)
		})
	} else {
		set.undo.Add(nil)
	}
}
func (set *IntSet) Delete(v int) {
	if set.Contains(v) {
		delete(set.data, v)
		set.undo.Add(func() {
			set.Add(v)
		})
	} else {
		set.undo.Add(nil)
	}
}

// Undo 控制逻辑：IntSet 依赖 Undo（一个没有参数的函数数组 协议）
type Undo []func()

func (undo *Undo) Add(f func()) {
	*undo = append(*undo, f)
}
func (undo *Undo) Undo() error {
	if len(*undo) == 0 {
		return errors.New("no functions to undo")
	}
	i := len(*undo) - 1
	if f := (*undo)[i]; f != nil {
		f()
		(*undo)[i] = nil
	}
	*undo = (*undo)[:i]
	return nil
}
