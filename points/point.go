package points

import "math"

type Point struct {
	Y float64 `json:"y"`
	X float64 `json:"x"`
}

func (p Point) DistanceFromOrigin(origin Point) float64 {
	return math.Abs(float64(origin.X)-float64(p.X)) + math.Abs(float64(origin.Y)-float64(p.Y))
}
