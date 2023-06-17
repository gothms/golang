package micro_kernel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// var _ Collector = (*Agent)(nil)
var WrongStateError = errors.New("can not take the operation in the current state")

type State int

const (
	Running State = iota
	Waiting
)

type CollectorError struct {
	CollectorErrors []error
}

func (ce CollectorError) Error() string {
	var strs []string
	for _, err := range ce.CollectorErrors {
		strs = append(strs, err.Error())
	}
	return strings.Join(strs, ";")
}

type EventReceiver interface {
	OnEvent(evt Event)
}
type Event struct {
	Source  string
	Content string
}
type Collector interface {
	Init(evtRcv EventReceiver) error
	Start(ctx context.Context) error
	Stop() error
	Destroy() error
}
type Agent struct {
	collectors map[string]Collector
	evtBuf     chan Event
	cancel     context.CancelFunc
	ctx        context.Context
	state      State
}

func NewAgent(sizeEvtBuf int) *Agent {
	agt := Agent{
		collectors: map[string]Collector{},
		evtBuf:     make(chan Event, sizeEvtBuf),
		state:      Waiting,
	}
	return &agt
}
func (agt *Agent) RegisterCollector(name string, collector Collector) error {
	if agt.state != Waiting {
		return WrongStateError
	}
	agt.collectors[name] = collector
	return collector.Init(agt)
}
func (agt *Agent) EventProcessGoroutine() {
	fmt.Println("EventProcessGoroutine")
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ { // 每收到 10 个，输出一次
			select {
			case evtSeg[i] = <-agt.evtBuf:
				fmt.Println(evtSeg)
			case <-agt.ctx.Done():
				return
			}
		}
		fmt.Println("??", evtSeg)
	}
}
func (agt *Agent) startCollectors() error {
	var (
		err  error
		errs CollectorError
		mut  sync.Mutex
	)
	for name, collector := range agt.collectors {
		go func(name string, collector2 Collector, ctx context.Context) {
			//defer mut.Unlock()
			defer func() { mut.Unlock() }()
			err = collector.Start(ctx)
			mut.Lock()
			if err != nil {
				errs.CollectorErrors = append(errs.CollectorErrors,
					errors.New(name+":"+err.Error()))
			}
		}(name, collector, agt.ctx)
	}
	return errs
}
func (agt *Agent) stopCollectors() error {
	var (
		err  error
		errs CollectorError
	)
	for name, collector := range agt.collectors {
		if err = collector.Stop(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors,
				errors.New(name+":"+err.Error()))
		}
	}
	return errs
}
func (agt *Agent) destroyCollectors() error {
	var (
		err  error
		errs CollectorError
	)
	for name, collector := range agt.collectors {
		if err = collector.Destroy(); err != nil {
			errs.CollectorErrors = append(errs.CollectorErrors,
				errors.New(name+":"+err.Error()))
		}
	}
	return errs
}

//func (agt *Agent) Init(evtRcv EventReceiver) error {
//
//}

func (agt *Agent) Start() error {
	//fmt.Println("start", agt.state)
	if agt.state != Waiting {
		return WrongStateError
	}
	agt.state = Running
	agt.ctx, agt.cancel = context.WithCancel(context.Background())
	go agt.EventProcessGoroutine()
	return agt.startCollectors()
}

func (agt *Agent) Stop() error {
	if agt.state != Running {
		return WrongStateError
	}
	agt.state = Waiting
	agt.cancel()
	return agt.stopCollectors()
}

func (agt *Agent) Destroy() error {
	if agt.state != Waiting {
		return WrongStateError
	}
	return agt.destroyCollectors()
}
func (agt *Agent) OnEvent(evt Event) {
	agt.evtBuf <- evt
}
