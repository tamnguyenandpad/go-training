package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day1/pool"
)

func main() {
	p := pool.NewWorkerPool(10, 100)
	p.Start(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for i := 0; i < 200; i++ {
		go func() {
			job := pool.Job{ID: fmt.Sprintf("job-%d", i), Handler: func(ctx context.Context) pool.Result {
				var state int
				if i%2 == 0 {
					state = 1
				} else {
					state = 0
				}
				return pool.Result{
					JobID: fmt.Sprintf("job-%d", i),
					State: state,
				}
			}}
			if err := p.Submit(job); err != nil {
				log.Printf("Failed to submit job: %v", err)
			}
		}()
	}

	// wait for stop signal to release the pool
	<-c
	// total goroutines before release
	log.Println("Total Goroutine before release the pool: ", runtime.NumGoroutine())
	p.Release()
	log.Printf("Jobs success:%d - failed:%d", p.TotalSucceed, p.TotalFailed)
	// Expect total 2 goroutines: main goroutine and the goroutine that wait for stop signal
	log.Println("Total Goroutine after release the pool: ", runtime.NumGoroutine())
}
