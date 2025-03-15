package worker_pool

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/pkg/worker_pool"
)

type WorkerPool interface {
	Release()
	Start(ctx context.Context)
	Submit(job Job) error
	Results() (totalSucceed, totalFailed int)
}

type Job interface {
	ID() string
	Name() string
	Handler() func() Result
}

type Result interface {
	JobID() string
	Success() bool
}

type Config struct {
	NumWorkers    int
	QueueCap      int
	IsNonBlocking bool
}

type workerPool struct {
	*worker_pool.WorkerPool
}

func NewWorkerPool(cfg Config) WorkerPool {
	opts := []worker_pool.PoolOpt{}
	if cfg.IsNonBlocking {
		opts = append(opts, worker_pool.WithNonBlocking)
	}
	return &workerPool{
		WorkerPool: worker_pool.NewWorkerPool(cfg.NumWorkers, cfg.QueueCap, opts...),
	}
}

func (wp *workerPool) Submit(job Job) error {
	h := func() worker_pool.Result {
		r := job.Handler()()
		state := 0
		if r.Success() {
			state = 1
		}
		return worker_pool.Result{
			JobID: r.JobID(),
			State: state,
		}
	}
	j := worker_pool.Job{
		ID:      job.ID(),
		Handler: h,
	}
	return wp.WorkerPool.Submit(j)
}
