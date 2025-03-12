package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/config"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal/pkg/worker"
)

type Server struct {
	config  config.Config
	pool    worker.WorkerPool
	service internal.AsyncGreeter

	httpServer *http.Server
}

func NewServer(cfg config.Config, pool worker.WorkerPool, service internal.AsyncGreeter) *Server {
	s := &Server{
		config:  cfg,
		pool:    pool,
		service: service,
	}

	// init http server
	mux := NewHTTPHandler(s)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.HTTPPort),
		Handler: mux,
	}
	s.httpServer = httpServer

	return s
}

func (s *Server) Start() {
	// start worker pool
	log.Printf("Worker pool: starting with %d workers - %d max jobs", s.config.PoolSize, s.config.MaxJobs)
	s.pool.Start(context.Background())

	// start HTTP server
	go func() {
		log.Printf("http.server: running on :%d", s.config.HTTPPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server: %v", err)
	}
	log.Printf("http.server: shutdown")

	// stop the worker pool
	s.pool.Release()
	totalJobSuccess, totalJobFailed := s.pool.Results()
	log.Printf("Worker pool: stopped")
	log.Printf("Jobs success:%d - failed:%d", totalJobSuccess, totalJobFailed)
}

func NewHTTPHandler(s *Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		names, ok := r.URL.Query()["name"]
		if !ok {
			http.Error(w, "missing name query parameter", http.StatusBadRequest)
			return
		}
		if err := s.service.Greeting(r.Context(), names); err != nil {
			type joinErr interface {
				Error() string
				Unwrap() []error
			}
			uw, ok := err.(joinErr)
			if !ok {
				http.Error(w, "failed to process greeting", http.StatusInternalServerError)
				return
			}
			w.Write(fmt.Appendf([]byte{}, `{"errors": %s}`, uw.Error()))
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("Greeting job submitted successfully!")
		w.Write([]byte(`{"message": "OK"}`))
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
