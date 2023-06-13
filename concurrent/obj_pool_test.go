package concurrent

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type ReusableObj struct {
}
type ObjPool struct {
	bufChan chan *ReusableObj
}

func NewObjPool(num int) *ObjPool {
	ch := make(chan *ReusableObj, num)
	for i := 0; i < num; i++ {
		ch <- new(ReusableObj)
	}
	return &ObjPool{ch}
}
func (p *ObjPool) GetObj(timeout time.Duration) (*ReusableObj, error) {
	select {
	case ret := <-p.bufChan:
		return ret, nil
	case <-time.After(timeout):
		return nil, errors.New("time out")
	}
}
func (p *ObjPool) ReleaseObj(obj *ReusableObj) error {
	select {
	case p.bufChan <- obj:
		return nil
	default:
		return errors.New("overflow")
	}
}
func TestObjPool(t *testing.T) {
	pool := NewObjPool(10)
	//if err := pool.ReleaseObj(&ReusableObj{}); err != nil {
	//	t.Error(err)
	//}	// overflow
	for i := 0; i < 11; i++ {
		if o, err := pool.GetObj(time.Second * 1); err != nil {
			t.Error(err)
		} else {
			fmt.Printf("%d,%T\n", i, o)
			if err = pool.ReleaseObj(o); err != nil {
				t.Error(err)
			}
		}
	}
	fmt.Println("Done")
}
