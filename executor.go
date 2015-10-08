package concurrency

import (
	"fmt"
	"golang.org/x/net/context"
	"runtime"
	"sync"
	"time"
)

type Executor struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewExecutor(ctx context.Context, timeoutMs int) *Executor {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(int64(timeoutMs))*time.Millisecond)
	executor := &Executor{
		ctx:    ctx,
		cancel: cancel,
	}
	return executor
}

func (p *Executor) Submit(req interface{}, tasks ...Tasker) *Future {
	var wg sync.WaitGroup

	resQ := make(chan TaskResult, len(tasks))
	for _, task := range tasks {
		wg.Add(1)
		go func(t Tasker) {
			select {
			case resQ <- func() (result TaskResult) {
				defer func() {
					if err := recover(); err != nil {
						result = TaskResult{nil, ErrStack(err)}
					}
				}()
				rsp, err := t.Call(p.ctx, req)
				result = TaskResult{rsp, err}
				return
			}():
			case <-p.ctx.Done():
				resQ <- TaskResult{nil, p.ctx.Err()}
			}
			wg.Done()
		}(task)
	}
	go func() {
		wg.Wait()
		close(resQ)
	}()
	return &Future{
		cancel: p.cancel,
		resQ:   resQ,
	}
}

func ErrStack(err interface{}) error {
	stack := make([]byte, 4096)
	size := runtime.Stack(stack, false)
	return fmt.Errorf("recover panic %s\nstack size is %d bytes\ntrace is %s", err, size, stack)
}
