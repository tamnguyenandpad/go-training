package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	R float64
}

type Rec struct {
	W float64
	H float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.R * c.R
}

func (c Circle) Perimeter() float64 {
	return math.Pi * 2 * c.R
}

func (r Rec) Area() float64 {
	return r.H * r.W
}

func (r Rec) Perimeter() float64 {
	return (r.H + r.W) * 2
}

func PrintShapeInfo(s Shape) {
	fmt.Printf("Area of %T %+v is %0.2f\n", s, s, s.Area())
	fmt.Printf("Perimeter of %T %+v is %0.2f\n", s, s, s.Perimeter())
}

func main() {
	c := Circle{
		R: 3,
	}

	r := Rec{
		W: 3,
		H: 4,
	}

	PrintShapeInfo(c)
	PrintShapeInfo(r)
}
