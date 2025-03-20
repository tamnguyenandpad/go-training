package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day1/pool"
	"go.uber.org/goleak"
)

func TestWorkerPool(t *testing.T) {
	p := pool.NewWorkerPool(1, 10)
	p.Start(context.Background())

	var mockJobHandlerSuccess = func(ctx context.Context) pool.Result {
		return pool.Result{JobID: "mock-job", State: 1}
	}
	var mockJobHandlerFailed = func(ctx context.Context) pool.Result {
		return pool.Result{JobID: "mock-job", State: 0}
	}

	// Submit 10 jobs - expected 5 success and 5 failed
	for i := 0; i < 10; i++ {
		job := pool.Job{ID: fmt.Sprintf("job-%d", i), Handler: mockJobHandlerSuccess}
		if i%2 == 0 {
			job.Handler = mockJobHandlerFailed
		}
		if err := p.Submit(job); err != nil {
			t.Errorf("Failed to submit job: %v", err)
		}
	}

	// Release the pool and closed
	p.Release()

	// Check results
	if p.TotalSucceed != 5 {
		t.Errorf("Expected 10 successful jobs, got %d", p.TotalSucceed)
	}
	if p.TotalFailed != 5 {
		t.Errorf("Expected 0 failed jobs, got %d", p.TotalFailed)
	}

	// verify no goroutine leak
	goleak.VerifyNone(t)
}

func TestWorkerPoolNonBlocking(t *testing.T) {
	// Create a worker pool with 3 workers and a job queue of size 5
	p := pool.NewWorkerPool(3, 5, pool.WithNonBlocking)
	p.Start(context.Background())

	wait := make(chan struct{})
	var handler = func(ctx context.Context) pool.Result {
		<-wait
		return pool.Result{JobID: "mock-job", State: 1}
	}

	submitErrs := make([]error, 0, 10)
	for i := 0; i < 10; i++ {
		job := pool.Job{ID: fmt.Sprintf("job-%d", i), Handler: handler}
		if err := p.Submit(job); err != nil {
			submitErrs = append(submitErrs, err)
		}
	}

	// Close the wait channel to unblock the blocking handlers
	close(wait)

	// Release the pool and closed
	p.Release()

	// Submit job progress from worst to best for the pool which has 3 workers and a job queue of size 5
	// - 3 workers process 3 jobs, 5 jobs are in the queue, 2 jobs are rejected => 8 jobs are submitted and handled at final
	// - 2 workers process 2 jobs, 5 jobs are in the queue, 3 jobs are rejected => 7 jobs are submitted and handled at final
	// - 1 worker process 1 job, 5 jobs are in the queue, 4 jobs are rejected => 6 jobs are submitted and handled at final
	// - 0 worker process 0 job, 5 jobs are in the queue, 5 jobs are rejected => 5 jobs are submitted and handled at final

	if len(submitErrs) > 5 {
		t.Errorf("Expected submit errs <= 5, got %d", len(submitErrs))
	}

	if p.TotalSucceed < 5 || p.TotalSucceed > 8 {
		t.Errorf("Expected 5 to 8 successful jobs, got %d", p.TotalSucceed)
	}

	goleak.VerifyNone(t)
}
