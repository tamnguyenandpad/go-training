package integration

import (
	"net/http"
	"testing"

	"github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/config"
)

func Test_Integration_Greeting(t *testing.T) {

	type args struct {
		httpMethod string
		cfg        *config.Config
		names      []string
	}
	type expected struct {
		out func(t *testing.T, out *testInstanceHelper)
	}
	type testCase struct {
		args     args
		expected expected
	}
	tests := map[string]testCase{
		"Sent request with a GET method": {
			args: args{
				httpMethod: http.MethodGet,
				cfg:        nil,
				names:      []string{"Alice", "Bob"},
			},
			expected: expected{
				out: func(t *testing.T, out *testInstanceHelper) {
					if out.httpResponseRecorder.Code != http.StatusMethodNotAllowed {
						t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, out.httpResponseRecorder.Code)
					}
					totalSuccess, totalFailed := out.wp.Results()
					if totalSuccess != 0 || totalFailed != 0 {
						t.Errorf("expected 0 success, 0 failed, got %d success, %d failed", totalSuccess, totalFailed)
					}
				},
			},
		},
		"Greeting 2 members with default config success": {
			args: args{
				httpMethod: http.MethodPost,
				cfg:        nil,
				names:      []string{"Alice", "Bob"},
			},
			expected: expected{
				out: func(t *testing.T, out *testInstanceHelper) {
					if out.httpResponseRecorder.Code != http.StatusOK {
						t.Errorf("expected status code %d, got %d", http.StatusOK, out.httpResponseRecorder.Code)
					}
					totalSuccess, totalFailed := out.wp.Results()
					if totalSuccess != 2 || totalFailed != 0 {
						t.Errorf("expected 2 success, 0 failed, got %d success, %d failed", totalSuccess, totalFailed)
					}
				},
			},
		},
		"Greeting 5 members with non-blocking worker pool partial success": {
			args: args{
				httpMethod: http.MethodPost,
				cfg: &config.Config{
					NumWorkers:    1,
					QueueCap:      2,
					IsNonBlocking: true,
				},
				names: []string{"Alice", "Bob", "Charlie", "David", "Eve"},
			},
			expected: expected{
				out: func(t *testing.T, out *testInstanceHelper) {
					if out.httpResponseRecorder.Code != http.StatusOK {
						t.Errorf("expected status code %d, got %d", http.StatusOK, out.httpResponseRecorder.Code)
					}
					totalSuccess, totalFailed := out.wp.Results()
					if totalSuccess < 2 || totalSuccess >= 5 || totalFailed != 0 {
						t.Errorf("expected 2 success, 0 failed, got %d success, %d failed", totalSuccess, totalFailed)
					}
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(tc.args.httpMethod, "/greeting", nil)

			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			q := req.URL.Query()
			for _, name := range tc.args.names {
				q.Add("name", name)
			}
			req.URL.RawQuery = q.Encode()

			testInstance := DoHTTPRequestWithConfig(t, req, tc.args.cfg)
			testInstance.wp.Release()
			tc.expected.out(t, testInstance)
		})
	}
}
