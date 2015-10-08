package concurrency

import (
	"golang.org/x/net/context"
)

type Tasker interface {
	Call(context.Context, interface{}) (interface{}, error)
}

type TaskFunc func(context.Context, interface{}) (interface{}, error)

func (p TaskFunc) Call(ctx context.Context, req interface{}) (interface{}, error) {
	return p(ctx, req)
}
