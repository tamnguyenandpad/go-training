package worker_pool

type Job struct {
	ID      string
	Handler JobHandler
}

type JobHandler func() Result