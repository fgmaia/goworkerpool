package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type IWorkerPool interface {
	Run(ctx context.Context)
	Results() <-chan Result
	GenerateFrom(jobsBulk []Job)
	GetDone() chan struct{}
}

type workerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	done         chan struct{}
}

func NewWorkerPool(workersCount int) IWorkerPool {
	return &workerPool{
		workersCount: workersCount,
		jobs:         make(chan Job, workersCount),
		results:      make(chan Result, workersCount),
		done:         make(chan struct{}),
	}
}

func (wp *workerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		// fan out worker goroutines
		//reading from jobs channel and
		//pushing calcs into results channel
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.done)
	close(wp.results)
}

func (wp *workerPool) Results() <-chan Result {
	return wp.results
}

func (wp *workerPool) GenerateFrom(jobsBulk []Job) {
	for i := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}

func (wp *workerPool) GetDone() chan struct{} {
	return wp.done
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			results <- job.Execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
