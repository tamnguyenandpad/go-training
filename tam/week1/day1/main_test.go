package main

import (
	"fmt"
	"testing"
)

func Test_fib(t *testing.T) {
	tests := [][2]int{
		{-1, -1},
		{0, -1},
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 2},
		{5, 3},
		{6, 5},
	}
	for _, tt := range tests {
		tc := fmt.Sprintf("Fibonacci number at position %d is %d", tt[0], tt[1])
		t.Run(tc, func(t *testing.T) {
			got := fib(tt[0])
			if got != tt[1] {
				t.Errorf("Fibonacci number at position %d is %d - got:%d", tt[0], tt[1], got)
			}
		})
	}
}
