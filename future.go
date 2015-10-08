package concurrency

import (
	"golang.org/x/net/context"
)

type TaskResult struct {
	Ret interface{}
	Err error
}

type Future struct {
	cancel context.CancelFunc
	resQ   chan TaskResult
}

func (p *Future) Get() []TaskResult {
	var result []TaskResult
	defer p.cancel()
	for r := range p.resQ {
		result = append(result, r)
	}
	return result
}
