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

	"github.com/joho/godotenv"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/cmd/app/server"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/config"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/pkg/worker_pool"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/service"
)

func main() {
	cfg := loadConfig()

	wp := worker_pool.NewWorkerPool(worker_pool.Config{
		NumWorkers:    cfg.NumWorkers,
		QueueCap:      cfg.QueueCap,
		IsNonBlocking: cfg.IsNonBlocking,
	})

	service := service.NewGreetingService(wp, cfg.BannedNames)

	server := server.NewServer(cfg, wp, service)

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
	// cd to cmd/app first
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	numWorkers, _ := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	queueCap, _ := strconv.Atoi(os.Getenv("QUEUE_CAP"))
	isNonBlocking, _ := strconv.ParseBool(os.Getenv("IS_NON_BLOCKING"))
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))

	banned := map[string]struct{}{}
	if bannedNames := os.Getenv("BANNED_NAMES"); bannedNames != "" {
		for _, name := range strings.Split(bannedNames, ",") {
			banned[name] = struct{}{}
		}
	}

	return config.Config{
		NumWorkers:    numWorkers,
		QueueCap:      queueCap,
		IsNonBlocking: isNonBlocking,
		HTTPPort:      httpPort,
		BannedNames:   banned,
	}
}
