package micro_kernel

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

var _ Collector = (*DemoCollector)(nil)

type DemoCollector struct {
	evtRcv   EventReceiver
	agtCtx   context.Context
	stopChan chan struct{}
	name     string
	content  string
}

func NewCollect(name string, content string) *DemoCollector {
	return &DemoCollector{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}
func (d *DemoCollector) Init(evtRcv EventReceiver) error {
	fmt.Println("initialize collector", d.name)
	d.evtRcv = evtRcv
	return nil
}

func (d *DemoCollector) Start(ctx context.Context) error {
	fmt.Println("start collector", d.name)
	for {
		select {
		case <-ctx.Done():
			d.stopChan <- struct{}{}
			break
		default:
			time.Sleep(time.Millisecond * 50)
			d.evtRcv.OnEvent(Event{d.name, d.content})
		}
	}
}

func (d *DemoCollector) Stop() error {
	fmt.Println("stop collector", d.name)
	select {
	case <-d.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("failed to stop for timeout")
	}
}

func (d *DemoCollector) Destroy() error {
	fmt.Println("destroy collector", d.name)
	return nil
}
func TestAgent(t *testing.T) {
	agt := NewAgent(100)
	c1 := NewCollect("c1", "1")
	c2 := NewCollect("c2", "2")
	agt.RegisterCollector(c1.name, c1)
	agt.RegisterCollector(c2.name, c2)
	agt.Start()
	fmt.Println(agt.Start())
	time.Sleep(time.Second * 1)
	agt.Stop()
	agt.Destroy()
}
