package main

import (
	"math"
	"testing"
)

func Test_ShapeInfo(t *testing.T) {
	type want struct {
		area      float64
		perimeter float64
	}

	tests := []struct {
		name  string
		shape Shape
		want  want
	}{
		{
			name:  "Circle",
			shape: Circle{R: 3},
			want: want{
				area:      math.Pi * 9,
				perimeter: 6 * math.Pi,
			},
		},
		{
			name:  "Rec",
			shape: Rec{W: 3, H: 4},
			want: want{
				area:      12,
				perimeter: 14,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.shape.Area(); got != tt.want.area {
				t.Errorf("Area of %T: got %0.2f, want %0.2f", tt.shape, got, tt.want.area)
			}
			if got := tt.shape.Perimeter(); got != tt.want.perimeter {
				t.Errorf("Perimeter of %T: got %0.2f, want %0.2f", tt.shape, got, tt.want.perimeter)
			}
		})
	}
}
