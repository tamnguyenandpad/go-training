package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal/pkg/worker"
)

const GreetingJobName = "Greeting"

type GreetingJob struct {
	id          string
	name        string
	bannedNames map[string]struct{}
}

func (j GreetingJob) ID() string {
	return j.id
}

func (j GreetingJob) Name() string {
	return GreetingJobName
}

func (j GreetingJob) Handler() func() worker.Result {
	return func() worker.Result {
		time.Sleep(1 * time.Second)
		state := 1
		if _, ok := j.bannedNames[j.Name()]; ok {
			state = 0
			log.Println("I'm sorry, I can't greet you!")
		} else {
			log.Printf("Hi %s, I'm a bot for greeting!\n", j.name)
		}

		return GreetingResult{
			GreetingJobID: j.id,
			State:         state,
		}
	}
}

type GreetingResult struct {
	GreetingJobID string
	State         int
}

func (r GreetingResult) JobID() string {
	return r.GreetingJobID
}

func (r GreetingResult) Success() bool {
	return r.State == 1
}

type greetingJobProducer struct {
	bannedNames map[string]struct{}
}

func (p greetingJobProducer) NewJob(name string) GreetingJob {
	// I don't check the banned names here, because I want to simulate the job handler failed case.
	return GreetingJob{
		id:          uuid.NewString(),
		name:        name,
		bannedNames: p.bannedNames,
	}
}

type AsyncGreeter interface {
	Greeting(ctx context.Context, names []string) error
}

type asyncGreetingService struct {
	pool        worker.WorkerPool
	jobProducer greetingJobProducer
}

func NewService(pool worker.WorkerPool, bannedNames map[string]struct{}) *asyncGreetingService {
	return &asyncGreetingService{
		pool: pool,
		jobProducer: greetingJobProducer{
			bannedNames: bannedNames,
		},
	}
}

func (s *asyncGreetingService) Greeting(ctx context.Context, names []string) error {
	var errs error
	var mux sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(names))
	for _, name := range names {
		go func() {
			defer wg.Done()
			job := s.jobProducer.NewJob(name)
			if err := s.pool.Submit(job); err != nil {
				mux.Lock()
				defer mux.Unlock()
				errs = errors.Join(errs, fmt.Errorf("failed to submit job for name %s: %v", name, err))
			}
		}()
	}
	wg.Wait()
	return errs
}
