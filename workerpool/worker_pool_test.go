package workerpool_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/fgmaia/goworkerpool/workerpool"
)

const (
	jobsCount   = 10
	workerCount = 2
)

func TestWorkerPool(t *testing.T) {
	workerPool := workerpool.NewWorkerPool(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go workerPool.GenerateFrom(testJobs())

	go workerPool.Run(ctx)

	for {
		select {
		case r, ok := <-workerPool.Results():
			if !ok {
				continue
			}

			i, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			val := r.Value.(int)
			if val != int(i)*2 {
				t.Fatalf("wrong value %v; expected %v", val, int(i)*2)
			}
		case <-workerPool.GetDone():
			return
		}
	}
}

func TestWorkerPool_TimeOut(t *testing.T) {
	workerPool := workerpool.NewWorkerPool(workerCount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond*10)
	defer cancel()

	go workerPool.Run(ctx)

	for {
		select {
		case r := <-workerPool.Results():
			if r.Err != nil && r.Err != context.DeadlineExceeded {
				t.Fatalf("expected error: %v; got: %v", context.DeadlineExceeded, r.Err)
			}
		case <-workerPool.GetDone():
			return
		}
	}
}

func TestWorkerPool_Cancel(t *testing.T) {
	workerPool := workerpool.NewWorkerPool(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())

	go workerPool.Run(ctx)
	cancel()

	for {
		select {
		case r := <-workerPool.Results():
			if r.Err != nil && r.Err != context.Canceled {
				t.Fatalf("expected error: %v; got: %v", context.Canceled, r.Err)
			}
		case <-workerPool.GetDone():
			return
		}
	}
}

func testJobs() []workerpool.Job {
	jobs := make([]workerpool.Job, jobsCount)
	for i := 0; i < jobsCount; i++ {
		jobs[i] = workerpool.Job{
			Descriptor: workerpool.JobDescriptor{
				ID:       workerpool.JobID(fmt.Sprintf("%v", i)),
				JType:    "anyType",
				Metadata: nil,
			},
			ExecFn: execFn,
			Args:   i,
		}
	}
	return jobs
}
