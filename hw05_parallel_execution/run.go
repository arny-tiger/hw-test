package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrCounter struct {
	mu      sync.Mutex
	counter int
}

func (c *ErrCounter) increment() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	errCounter := ErrCounter{sync.Mutex{}, 0}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go runWorker(tasksCh, stopCh, &errCounter, &wg)
	}

	var err error
	for _, task := range tasks {
		tasksCh <- task
		if errCounter.counter >= m {
			err = ErrErrorsLimitExceeded
			break
		}
	}

	close(stopCh)
	close(tasksCh)
	wg.Wait()
	return err
}

func runWorker(tasksCh <-chan Task, stopCh <-chan struct{}, errCounter *ErrCounter, wg *sync.WaitGroup) {
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
				errCounter.increment()
			}
		}
	}
}
