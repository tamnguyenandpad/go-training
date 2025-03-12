package integration

import (
	"net/http"
	"testing"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/config"
)

func Test_integration_Greeting(t *testing.T) {
	type args struct {
		names      []string
		httpMethod string
		cfg        *config.Config
	}

	tests := []struct {
		name   string
		args   args
		expect func(t *testing.T, out *testInstanceHelper)
	}{
		{
			name: "Sent request with a GET method",
			args: args{
				httpMethod: http.MethodGet,
				names:      []string{"Alice", "Bob"},
				cfg:        nil,
			},
			expect: func(t *testing.T, out *testInstanceHelper) {
				if out.httpResponseRecorder.Code != http.StatusMethodNotAllowed {
					t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, out.httpResponseRecorder.Code)
				}
				totalSuccess, totalFailed := out.pool.Results()
				if totalSuccess != 0 || totalFailed != 0 {
					t.Errorf("expected 0 success, 0 failed, got %d success, %d failed", totalSuccess, totalFailed)
				}
			},
		},
		{
			name: "Greeting 2 members with default config success",
			args: args{
				httpMethod: http.MethodPost,
				names:      []string{"Alice", "Bob"},
				cfg:        nil,
			},
			expect: func(t *testing.T, out *testInstanceHelper) {
				if out.httpResponseRecorder.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, out.httpResponseRecorder.Code)
				}
				totalSuccess, totalFailed := out.pool.Results()
				if totalSuccess != 2 || totalFailed != 0 {
					t.Errorf("expected 2 success, 0 failed, got %d success, %d failed", totalSuccess, totalFailed)
				}
			},
		},
		{
			name: "Greeting 5 members with non-blocking worker pool partial success",
			args: args{
				httpMethod: http.MethodPost,
				names:      []string{"Alice", "Bob", "Charlie", "David", "Eve"},
				cfg: &config.Config{
					PoolSize:              1,
					MaxJobs:               2,
					WorkerPoolNonBlocking: true,
				},
			},
			expect: func(t *testing.T, out *testInstanceHelper) {
				if out.httpResponseRecorder.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, out.httpResponseRecorder.Code)
				}
				totalSuccess, totalFailed := out.pool.Results()
				if totalSuccess < 2 || totalSuccess >= 5 || totalFailed != 0 {
					t.Errorf("expected 1 < total success < 5 and failed = 0, got %d success, %d failed", totalSuccess, totalFailed)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the request
			req, err := http.NewRequest(tt.args.httpMethod, "/greeting", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			// build the query params
			q := req.URL.Query()
			for _, name := range tt.args.names {
				q.Add("name", name)
			}
			req.URL.RawQuery = q.Encode()

			testInstance := DoHTTPRequestWithConfig(t, req, tt.args.cfg)
			testInstance.pool.Release()
			tt.expect(t, testInstance)
		})
	}
}
