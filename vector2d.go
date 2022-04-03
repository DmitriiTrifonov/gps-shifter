package gps_shifter

import (
	"fmt"
	"math"
)

type Vector2D struct {
	x float64
	y float64
}

func NewVector2D(x2, y2, x1, y1 float64) Vector2D {
	return Vector2D{x: x2 - x1, y: y2 - y1}
}

func (v Vector2D) String() string {
	return fmt.Sprintf("x: %f, y: %f", v.x, v.y)
}

func (v Vector2D) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}
