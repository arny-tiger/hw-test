package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var Lock sync.Mutex

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	errCounter := 0

	wg.Add(n)
	for i := 0; i < n; i++ {
		go runWorker(tasksCh, stopCh, &errCounter, &wg)
	}

	var err error
	for _, task := range tasks {
		tasksCh <- task
		Lock.Lock()
		isLimitExceeded := errCounter >= m
		Lock.Unlock()
		if isLimitExceeded {
			err = ErrErrorsLimitExceeded
			break
		}
	}

	close(stopCh)
	close(tasksCh)
	wg.Wait()
	return err
}

func runWorker(tasksCh <-chan Task, stopCh <-chan struct{}, errCounter *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopCh:
			return
		case task, ok := <-tasksCh:
			if !ok {
				return
			}
			err := task()
			if err != nil {
				Lock.Lock()
				*errCounter++
				Lock.Unlock()
			}
		}
	}
}
