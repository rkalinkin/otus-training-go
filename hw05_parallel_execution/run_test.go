package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunWithoutSleep(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("concurrency test without sleep", func(t *testing.T) {
		tasksCount := 10
		workersCount := 5
		maxErrorsCount := 1

		var startedTasks int32

		startCh := make(chan struct{})
		tasks := make([]Task, 0, tasksCount)

		for range tasksCount {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&startedTasks, 1)
				<-startCh
				return nil
			})
		}

		done := make(chan error)
		go func() {
			done <- Run(tasks, workersCount, maxErrorsCount)
		}()

		time.Sleep(50 * time.Millisecond)
		started := atomic.LoadInt32(&startedTasks)
		require.GreaterOrEqual(t, started, int32(2),
			"tasks should be started concurrently by multiple workers, got %d", started)

		close(startCh)

		err := <-done
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), atomic.LoadInt32(&startedTasks), "all tasks should be started")
	})
}

func TestRunWithZeroMaxErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("zero max errors should ignore errors", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := range tasksCount {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "all tasks should be executed even with errors when m=0")
	})

	t.Run("negative max errors should ignore errors", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := range tasksCount {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := -1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, int32(tasksCount), runTasksCount, "all tasks should be executed even with errors when m<0")
	})
}
