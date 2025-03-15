package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	mock_worker_pool "github.com/tuannguyenandpadcojp/go-training/tam/week2/day2/internal/pkg/worker_pool/mock"
	"go.uber.org/mock/gomock"
)

func Test_Service_Greet(t *testing.T) {
	type fields struct {
		worker_pool *mock_worker_pool.MockWorkerPool
	}
	type args struct {
		ctx         context.Context
		names       []string
		bannedNames map[string]struct{}
	}
	type handleErrorFunc func(t *testing.T, err error, prefixMsg ...string)

	expectNilErr := func(t *testing.T, err error, prefixMsg ...string) {
		t.Helper()
		if err != nil {
			t.Errorf("%s expected nil error, got %v", strings.Join(prefixMsg, " "), err)
		}
	}

	expectJoinErrs := func(t *testing.T, err error, n int, prefixMsg ...string) {
		t.Helper()
		if err == nil {
			t.Errorf("%s expected error, got nil", strings.Join(prefixMsg, " "))
		}
		type joinErrors interface {
			Unwrap() []error
		}
		uw, ok := err.(joinErrors)
		if !ok {
			t.Errorf("%s expected error to implement Unwrap method, got %T", strings.Join(prefixMsg, " "), err)
		}
		if len(uw.Unwrap()) != n {
			t.Errorf("%s expected %d errors, got %d", strings.Join(prefixMsg, " "), n, len(uw.Unwrap()))
		}
	}

	type testCase struct {
		prepare func(f *fields)
		args    args
		wantErr handleErrorFunc
	}

	tests := map[string]testCase{
		"Submit Greeting job successfully": {
			args: args{
				ctx:   context.Background(),
				names: []string{"Tam", "Tai"},
			},
			prepare: func(f *fields) {
				f.worker_pool.EXPECT().Submit(gomock.Any()).Return(nil).Times(2)
			},
			wantErr: expectNilErr,
		},
		"Submit Greeting 1 job successfully, 1 job failed": {
			args: args{
				ctx:   context.Background(),
				names: []string{"Tam", "Tai"},
			},
			prepare: func(f *fields) {
				f.worker_pool.EXPECT().Submit(gomock.Any()).Return(nil)
				f.worker_pool.EXPECT().Submit(gomock.Any()).Return(errors.New("error")).Times(1)
			},
			wantErr: func(t *testing.T, err error, prefixMsg ...string) {
				expectJoinErrs(t, err, 1, prefixMsg...)
			},
		},
		"Submit Greeting 2 jobs failed": {
			args: args{
				ctx:         context.Background(),
				names:       []string{"Tam", "Tai"},
				bannedNames: map[string]struct{}{},
			},
			prepare: func(f *fields) {
				f.worker_pool.EXPECT().Submit(gomock.Any()).Return(errors.New("error")).Times(2)
			},
			wantErr: func(t *testing.T, err error, prefixMsg ...string) {
				expectJoinErrs(t, err, 2, prefixMsg...)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			f := fields{
				worker_pool: mock_worker_pool.NewMockWorkerPool(mockCtrl),
			}
			if tc.prepare != nil {
				tc.prepare(&f)
			}
			s := NewGreetingService(f.worker_pool, tc.args.bannedNames)
			err := s.Greet(tc.args.ctx, tc.args.names)
			tc.wantErr(t, err, "Service_Greet")
		})
	}
}
