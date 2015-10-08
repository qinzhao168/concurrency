package concurrency

import (
	"errors"
	"golang.org/x/net/context"
	"log"
	"sync"
	"testing"
	"time"
)

func mockFuncA(ctx context.Context, req interface{}) (interface{}, error) {
	log.Println("Tasker in mockFuncA, req=", req)
	return "mockFuncA", nil
}

func mockFuncB(ctx context.Context, req interface{}) (interface{}, error) {
	log.Println("Tasker in mockFuncB, req=", req)
	time.Sleep(200 * time.Millisecond)
	return "mockFuncB", nil
}

func mockFuncC(ctx context.Context, req interface{}) (interface{}, error) {
	log.Println("Tasker in mockFuncC, req=", req)
	time.Sleep(600 * time.Millisecond)
	return "mockFuncC", errors.New("timeout")
}

func mockFuncD(ctx context.Context, req interface{}) (interface{}, error) {
	log.Println("Tasker in mockFuncD, req=", req)
	panic("boom!!")
	return "mockFuncD", errors.New("panic")
}

func TestTaskFunc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var wg sync.WaitGroup

	taskerA := TaskFunc(mockFuncA)
	taskerB := TaskFunc(mockFuncB)
	taskerC := TaskFunc(mockFuncC)
	taskerD := TaskFunc(mockFuncD)

	executor1 := NewExecutor(ctx, 400)
	req1 := "req in executor 1"

	wg.Add(1)
	future1 := executor1.Submit(req1, taskerA, taskerA, taskerB)
	for index, r := range future1.Get() {
		log.Printf("index=%d,ret=%v", index, r)
	}
	wg.Done()

	executor2 := NewExecutor(ctx, 400)
	req2 := "req in executor 2"
	wg.Add(1)
	future2 := executor2.Submit(req2, taskerA, taskerB, taskerC)
	for index, r := range future2.Get() {
		log.Printf("index=%d,ret=%v", index, r)
	}
	wg.Done()

	executor3 := NewExecutor(ctx, 400)
	req3 := "req in executor 3"
	wg.Add(1)
	future3 := executor3.Submit(req3, taskerA, taskerB, taskerC, taskerD)
	for index, r := range future3.Get() {
		log.Printf("index=%d,ret=%v", index, r)
	}
	wg.Done()
	time.Sleep(100 * time.Millisecond)
	wg.Wait()
	log.Println("Finish test TaskFunc")
}
