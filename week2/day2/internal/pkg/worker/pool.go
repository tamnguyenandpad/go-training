package worker

import (
	w "github.com/tuannguyenandpadcojp/go-training/week2/day2/pkg/worker"
)

type Config struct {
	PoolSize    int
	MaxJobs     int
	NonBlocking bool
}

type pool struct {
	*w.Pool
}

func NewPool(cfg Config) WorkerPool {
	opts := []w.PoolOpt{}
	if cfg.NonBlocking {
		opts = append(opts, w.WithNonBlocking)
	}

	return &pool{
		Pool: w.NewWorkerPool(cfg.PoolSize, cfg.MaxJobs, opts...),
	}
}

func (p *pool) Submit(job Job) error {
	h := func() w.Result {
		r := job.Handler()()
		state := 0
		if r.Success() {
			state = 1
		}
		return w.Result{
			JobID: r.JobID(),
			State: state,
		}
	}
	j := w.Job{
		ID:      job.ID(),
		Handler: h,
	}
	return p.Pool.Submit(j)
}
