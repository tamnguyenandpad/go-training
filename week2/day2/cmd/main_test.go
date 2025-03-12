package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/config"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name             string
		prepareVariables func()
		want             config.Config
	}{
		{
			name:             "default config",
			prepareVariables: func() {},
			want: config.Config{
				PoolSize:              2,
				MaxJobs:               2,
				WorkerPoolNonBlocking: false,
				HTTPPort:              8080,
			},
		},
		{
			name: "custom config",
			prepareVariables: func() {
				_ = os.Setenv("POOL_SIZE", "10")
				_ = os.Setenv("MAX_JOBS", "20")
				_ = os.Setenv("WORKER_POOL_NON_BLOCKING", "true")
				_ = os.Setenv("HTTP_PORT", "8081")
			},
			want: config.Config{
				PoolSize:              10,
				MaxJobs:               20,
				WorkerPoolNonBlocking: true,
				HTTPPort:              8081,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepareVariables != nil {
				tt.prepareVariables()
			}
			if got := loadConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
