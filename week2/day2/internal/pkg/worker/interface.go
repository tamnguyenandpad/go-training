package worker

import (
	"context"
)

type WorkerPool interface {
	Submit(job Job) error
	Start(ctx context.Context)
	Release()
	Results() (totalSucceed int, totalFailed int)
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
