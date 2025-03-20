package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/cmd/app/server"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/config"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/pkg/worker_pool"
	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/service"
)

// TODO: can not deal with godotenv.Load
func TestLoadConfig(t *testing.T) {
	// Mock config
	if err := os.Setenv("NUM_WORKERS", "10"); err != nil {
		t.Fatalf("Failed to set NUM_WORKERS: %v", err)
	}
	if err := os.Setenv("QUEUE_CAP", "20"); err != nil {
		t.Fatalf("Failed to set QUEUE_CAP: %v", err)
	}
	if err := os.Setenv("IS_NON_BLOCKING", "true"); err != nil {
		t.Fatalf("Failed to set IS_NON_BLOCKING: %v", err)
	}
	if err := os.Setenv("HTTP_PORT", "8081"); err != nil {
		t.Fatalf("Failed to set HTTP_PORT: %v", err)
	}
	if err := os.Setenv("BANNED_NAMES", "name1,name2"); err != nil {
		t.Fatalf("Failed to set BANNED_NAMES: %v", err)
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

	expectedConfig := config.Config{
		NumWorkers:    numWorkers,
		QueueCap:      queueCap,
		IsNonBlocking: isNonBlocking,
		HTTPPort:      httpPort,
		BannedNames:   banned,
	}

	cfg := loadConfig()
	assert.Equal(t, expectedConfig, cfg)
}

func TestServerStartStop(t *testing.T) {
	cfg := loadConfig()

	wp := worker_pool.NewWorkerPool(worker_pool.Config{
		NumWorkers:    cfg.NumWorkers,
		QueueCap:      cfg.QueueCap,
		IsNonBlocking: cfg.IsNonBlocking,
	})

	service := service.NewGreetingService(wp, cfg.BannedNames)

	server := server.NewServer(cfg, wp, service)

	go server.Start()

	time.Sleep(1 * time.Second) // Give the server time to start

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Stop(ctx)
}
