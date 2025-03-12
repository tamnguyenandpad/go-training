package internal_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/tuannguyenandpadcojp/go-training/week2/day2/internal"
	mock_internal "github.com/tuannguyenandpadcojp/go-training/week2/day2/internal/pkg/worker/mock"
)

type expectErrFunc func(t *testing.T, err error, prefixMsg ...string)

func expectNilErr(t *testing.T, err error, prefixMsg ...string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s expected nil error, got %v", strings.Join(prefixMsg, " "), err)
	}
}

func expectJoinErrs(t *testing.T, err error, n int, prefixMsg ...string) {
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

func Test_asyncGreetingService_Greeting(t *testing.T) {
	type args struct {
		ctx         context.Context
		names       []string
		bannedNames map[string]struct{}
	}
	tests := []struct {
		name      string
		args      args
		prepare   func(pool *mock_internal.MockWorkerPool)
		expectErr expectErrFunc
	}{
		{
			name: "Submit Greeting job successfully",
			args: args{
				ctx:   context.Background(),
				names: []string{"Alice", "Bob"},
			},
			prepare: func(pool *mock_internal.MockWorkerPool) {
				// Expect to call Submit 2 times for 2 names in the list
				// All calls to Submit should success by returning nil
				pool.EXPECT().Submit(gomock.Any()).Return(nil).Times(2)
			},
			expectErr: expectNilErr,
		},
		{
			name: "Submit Greeting 1 job failed",
			args: args{
				ctx:   context.Background(),
				names: []string{"Alice", "Bob"},
			},
			prepare: func(pool *mock_internal.MockWorkerPool) {
				// Expect to call Submit 2 times for 2 names in the list
				// The first call to Submit should success by returning nil
				// The second call to Submit should fail by returning an error
				pool.EXPECT().Submit(gomock.Any()).Return(nil).Times(1)
				pool.EXPECT().Submit(gomock.Any()).Return(errors.New("submit error")).Times(1)
			},
			expectErr: func(t *testing.T, err error, prefixMsg ...string) {
				expectJoinErrs(t, err, 1, prefixMsg...)
			},
		},
		{
			name: "Submit Greeting 2 job failed",
			args: args{
				ctx:   context.Background(),
				names: []string{"Alice", "Bob", "Charlie"},
			},
			prepare: func(pool *mock_internal.MockWorkerPool) {
				// Expect to call Submit 2 times for 2 names in the list
				// The first call to Submit should success by returning nil
				// The second call to Submit should fail by returning an error
				pool.EXPECT().Submit(gomock.Any()).Return(nil).Times(1)
				pool.EXPECT().Submit(gomock.Any()).Return(errors.New("submit error")).Times(2)
			},
			expectErr: func(t *testing.T, err error, prefixMsg ...string) {
				expectJoinErrs(t, err, 2, prefixMsg...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			pool := mock_internal.NewMockWorkerPool(mockCtrl)
			if tt.prepare != nil {
				tt.prepare(pool)
			}
			s := internal.NewService(pool, tt.args.bannedNames)
			err := s.Greeting(tt.args.ctx, tt.args.names)
			tt.expectErr(t, err, "asyncGreetingService_Greeting: ")
		})
	}
}
