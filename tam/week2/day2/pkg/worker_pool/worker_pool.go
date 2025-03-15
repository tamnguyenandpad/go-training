package worker_pool

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type WorkerPool struct {
	mutex      sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	jobWG      sync.WaitGroup
	resultWG   sync.WaitGroup

	numWorkers    int
	jobQueue      chan Job
	resultQueue   chan Result
	isRuning      bool
	isNonBlocking bool
	totalSucceed  int
	totalFailed   int
}

type PoolOpt func(p *WorkerPool)

func WithNonBlocking(p *WorkerPool) {
	p.isNonBlocking = true
}

func NewWorkerPool(numWorkers int, queueCap int, opts ...PoolOpt) *WorkerPool {
	wp := &WorkerPool{
		mutex:    sync.Mutex{},
		jobWG:    sync.WaitGroup{},
		resultWG: sync.WaitGroup{},

		numWorkers:  numWorkers,
		jobQueue:    make(chan Job, queueCap),
		resultQueue: make(chan Result, queueCap),
	}

	for _, opt := range opts {
		opt(wp)
	}

	return wp
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if wp.isRuning {
		log.Println("Worker Pool is running ...")
		return
	}
	wp.isRuning = true
	wp.ctx, wp.cancelFunc = context.WithCancel(ctx)
	wp.jobWG.Add(wp.numWorkers)
	// spawn worker goroutine
	for i := 0; i < wp.numWorkers; i++ {
		go worker(wp.ctx, i, wp.jobQueue, wp.resultQueue, &wp.jobWG)
	}

	// aggregate job's result
	wp.resultWG.Add(1)
	go func() {
		defer wp.resultWG.Done()
		for result := range wp.resultQueue {
			if result.State == 1 {
				wp.totalSucceed++
			} else {
				wp.totalFailed++
			}
		}
	}()
}

func (wp *WorkerPool) Release() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	if !wp.isRuning {
		log.Println("Worker Pool is not running ...")
		return
	}
	// close the Jobs channel to prevent dispatcher send jobs
	close(wp.jobQueue)
	// wait for all workers to finish processing the rest of jobs
	wp.jobWG.Wait()
	// call context.CancelFunc to notify all workers to stop
	wp.cancelFunc()
	// close the result channel to stop the goroutine aggregate the job's result
	close(wp.resultQueue)
	// wait for the goroutine aggregate the job result is done
	wp.resultWG.Wait()
	wp.isRuning = false
}

func (wp *WorkerPool) Submit(job Job) error {
	if !wp.isRuning {
		return fmt.Errorf("pool is closed")
	}
	if wp.isNonBlocking {
		select {
		case wp.jobQueue <- job:
			// If the job channel has space, the job is sent and the function returns nil.
			return nil
		default:
			// If the job channel is full, the default case is executed and an error is returned.
			return fmt.Errorf("job queue is full")
		}
	} else {
		// blocking if the jobs channel is full
		wp.jobQueue <- job
		return nil
	}
}

func (wp *WorkerPool) Results() (totalSucceed, totalFailed int) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	return wp.totalSucceed, wp.totalFailed
}

func worker(ctx context.Context, jobQueueIndex int, jobQueue <-chan Job, resultQueue chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobQueue:
			if !ok {
				return
			}
			// Process the job
			log.Printf("Worker %d processed job %s", jobQueueIndex, job.ID)
			resultQueue <- job.Handler()
		}
	}
}
