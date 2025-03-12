package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/cmd/server"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/config"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal/pkg/worker"
)

type testInstanceHelper struct {
	pool worker.WorkerPool

	// response
	httpResponseRecorder *httptest.ResponseRecorder
}

func initDependencies(t *testing.T, cfg *config.Config) (http.Handler, worker.WorkerPool) {
	t.Helper()
	// Start the server
	c := config.Config{
		PoolSize:              2,
		MaxJobs:               2,
		WorkerPoolNonBlocking: true,
		HTTPPort:              8081,
		BannedNames:           map[string]struct{}{},
	}
	if cfg != nil {
		c = *cfg
	}

	// init worker pool
	pool := worker.NewPool(worker.Config{
		PoolSize:    c.PoolSize,
		MaxJobs:     c.MaxJobs,
		NonBlocking: c.WorkerPoolNonBlocking,
	})

	// init services
	service := internal.NewService(pool, c.BannedNames)

	s := server.NewServer(c, pool, service)
	pool.Start(context.Background())
	httpHandler := server.NewHTTPHandler(s)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Stop(ctx)
	})
	return httpHandler, pool
}

func DoHTTPRequestWithConfig(t *testing.T, re *http.Request, cfg *config.Config) *testInstanceHelper {
	handler, pool := initDependencies(t, cfg)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, re)
	return &testInstanceHelper{
		pool:                 pool,
		httpResponseRecorder: rr,
	}
}
