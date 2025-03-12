package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/cmd/server"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/config"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal"
	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal/pkg/worker"
)

func main() {
	cfg := loadConfig()

	// init worker pool
	pool := worker.NewPool(worker.Config{
		PoolSize:    cfg.PoolSize,
		MaxJobs:     cfg.MaxJobs,
		NonBlocking: cfg.WorkerPoolNonBlocking,
	})

	// init async greeter
	service := internal.NewService(pool, cfg.BannedNames)

	// new server
	server := server.NewServer(cfg, pool, service)

	// start server
	server.Start()

	// init signal channel
	sig := make(chan os.Signal, 1)
	// notify the signal channel when receiving interrupt signal
	// syscall.SIGTERM is a signal sent to a process to request its termination
	// os.Interrupt is a signal sent to a process by its controlling terminal when a user wishes to interrupt the process
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// wait for interrupt signal to shutdown the server and release the worker pool
	<-sig

	// shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Stop(ctx)
}

func loadConfig() config.Config {
	// Load the configuration from os environment
	c := config.Config{
		PoolSize:              2,
		MaxJobs:               2,
		WorkerPoolNonBlocking: false,
		BannedNames:           map[string]struct{}{},
		HTTPPort:              8080,
	}
	if poolSize := os.Getenv("POOL_SIZE"); poolSize != "" {
		pz, err := strconv.Atoi(poolSize)
		if err != nil {
			log.Fatalf("Invalid pool size: %v", err)
		}
		c.PoolSize = pz
	}
	if maxJobs := os.Getenv("MAX_JOBS"); maxJobs != "" {
		mj, err := strconv.Atoi(maxJobs)
		if err != nil {
			log.Fatalf("Invalid max jobs: %v", err)
		}
		c.MaxJobs = mj
	}
	if nonBlocking := os.Getenv("WORKER_POOL_NON_BLOCKING"); nonBlocking != "" {
		nb, err := strconv.ParseBool(nonBlocking)
		if err != nil {
			log.Fatalf("Invalid worker pool non-blocking: %v", err)
		}
		c.WorkerPoolNonBlocking = nb
	}
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		hp, err := strconv.Atoi(httpPort)
		if err != nil {
			log.Fatalf("Invalid http port: %v", err)
		}
		c.HTTPPort = hp
	}
	if bannedNames := os.Getenv("BANNED_NAMES"); bannedNames != "" {
		banned := map[string]struct{}{}
		for _, name := range strings.Split(bannedNames, ",") {
			banned[name] = struct{}{}
		}
		c.BannedNames = banned
	}

	return c
}
