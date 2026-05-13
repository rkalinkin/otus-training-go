package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error { //nolint:gocognit
	if len(tasks) == 0 {
		return nil
	}

	shouldCheckErrors := m > 0

	var errCount atomic.Int64
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCh := make(chan Task, n)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)

		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case taskCh <- task:
			}
		}
	}()

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskCh:
					if !ok {
						return
					}

					err := task()
					if err != nil && shouldCheckErrors {
						newCount := errCount.Add(1)
						if newCount >= int64(m) {
							cancel()
							return
						}
					}
				}
			}
		}()
	}

	wg.Wait()

	if shouldCheckErrors && errCount.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
