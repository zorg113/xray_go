package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrErrorsNoWorkers     = errors.New("errors limit exceeded")
)

type Result struct {
	err error
}

type Task func() error

func (r Result) Err() error {
	return r.err
}

// Push tasks 
func pushTasks(taskCh chan<- Task, tasks []Task, doneCh <-chan struct{}) {
	defer close(taskCh)
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-doneCh:
			return
		}
	}
}

// Run worker 
func runWorker(taskCh <-chan Task, resCh chan<- Result, doneCh <-chan struct{}) {
	for {
		task, ok := <-taskCh
		if !ok {
			return
		}
		err := task()
		select {
		case resCh <- Result{err: err}:
		case <-doneCh:
			return
		}
	}
}

// Run workers 
func runWorkers(count int, taskCh <-chan Task, resCh chan<- Result, doneCh <-chan struct{}) {
	defer close(resCh)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			runWorker(taskCh, resCh, doneCh)
		}()
	}
	wg.Wait()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrErrorsNoWorkers
	}

	ignoreErrors := m < 1

	taskCh := make(chan Task)
	resCh := make(chan Result)
	doneCh := make(chan struct{})

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		runWorkers(n, taskCh, resCh, doneCh)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		pushTasks(taskCh, tasks, doneCh)
	}()

	errorCount := 0
	var err error
	for {
		result, ok := <-resCh
		if !ok {
			break
		}
		if ignoreErrors {
			continue
		}
		if result.Err() != nil {
			errorCount++
		}
		if errorCount >= m {
			err = ErrErrorsLimitExceeded
			close(doneCh)
			break
		}
	}
	waitGroup.Wait()
	return err
}
