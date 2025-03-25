package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/cmd/app/server"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/config"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/pkg/worker_pool"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/service"
)

type testInstanceHelper struct {
	wp worker_pool.WorkerPool

	// response
	httpResponseRecorder *httptest.ResponseRecorder
}

func initDependencies(t *testing.T, cfg *config.Config) (http.Handler, worker_pool.WorkerPool) {
	t.Helper()
	// Start the server
	c := config.Config{
		NumWorkers:    2,
		QueueCap:      2,
		IsNonBlocking: true,
		HTTPPort:      8081,
		BannedNames:   map[string]struct{}{},
	}
	if cfg != nil {
		c = *cfg
	}

	// init worker pool
	wp := worker_pool.NewWorkerPool(worker_pool.Config{
		NumWorkers:    c.NumWorkers,
		QueueCap:      c.QueueCap,
		IsNonBlocking: c.IsNonBlocking,
	})

	// init services
	service := service.NewGreetingService(wp, c.BannedNames)

	s := server.NewServer(c, wp, service)
	wp.Start(context.Background())
	httpHandler := server.NewHTTPHandler(s)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Stop(ctx)
	})
	return httpHandler, wp
}

func DoHTTPRequestWithConfig(t *testing.T, re *http.Request, cfg *config.Config) *testInstanceHelper {
	handler, wp := initDependencies(t, cfg)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, re)
	return &testInstanceHelper{
		wp:                   wp,
		httpResponseRecorder: rr,
	}
}
