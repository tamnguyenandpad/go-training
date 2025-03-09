package main

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3, 10)
	pool.Start(context.Background())

	var mockJobHandlerSuccess = func(ctx context.Context) Result {
		return Result{JobID: "mock-job", State: 1}
	}
	var mockJobHandlerFailed = func(ctx context.Context) Result {
		return Result{JobID: "mock-job", State: 0}
	}

	// Submit 10 jobs - expected 5 success and 5 failed
	for i := 0; i < 10; i++ {
		job := Job{ID: fmt.Sprintf("job-%d", i), Handler: mockJobHandlerSuccess}
		if i%2 == 0 {
			job.Handler = mockJobHandlerFailed
		}
		if err := pool.Submit(job); err != nil {
			t.Errorf("Failed to submit job: %v", err)
		}
	}

	// Release the pool and closed
	pool.Release()

	// Check results
	if pool.totalSucceed != 5 {
		t.Errorf("Expected 10 successful jobs, got %d", pool.totalSucceed)
	}
	if pool.totalFailed != 5 {
		t.Errorf("Expected 0 failed jobs, got %d", pool.totalFailed)
	}

	// verify no goroutine leak
	goleak.VerifyNone(t)
}

func TestWorkerPoolNonBlocking(t *testing.T) {
	// Create a worker pool with 3 workers and a job queue of size 5
	pool := NewWorkerPool(3, 5, WithNonBlocking)
	pool.Start(context.Background())

	wait := make(chan struct{})
	var blockingHandler = func(ctx context.Context) Result {
		<-wait
		return Result{JobID: "mock-job", State: 1}
	}

	var handler = func(ctx context.Context) Result {
		return Result{JobID: "mock-job", State: 0}
	}

	// Submit 10 jobs and the job 1 -> 5 takes longer time to process
	// The job 6 -> 10 will be rejected
	// We expect the pool to process 5 jobs successfully and reject 5 jobs
	// We simulate the job 1 -> 5 takes longer time to process by using a blocking handler
	// until the job 6 -> 10 is submitted
	for i := 0; i < 10; i++ {
		job := Job{ID: fmt.Sprintf("job-%d", i), Handler: handler}
		if i < 5 {
			job.Handler = blockingHandler
		}
		err := pool.Submit(job)
		if i < 5 && err != nil {
			t.Errorf("Failed to submit job: %v", err)
		}
		if i >= 5 && err == nil {
			t.Errorf("Expected job queue to be full, but job was submitted successfully")
		}
	}

	// Close the wait channel to unblock the blocking handlers
	close(wait)

	// Release the pool ans closed
	pool.Release()

	// Check results
	if pool.totalSucceed != 5 {
		t.Errorf("Expected 5 successful jobs, got %d", pool.totalSucceed)
	}
	if pool.totalFailed != 0 {
		t.Errorf("Expected 0 failed jobs, got %d", pool.totalFailed)
	}

	// verify no goroutine leak
	goleak.VerifyNone(t)
}
