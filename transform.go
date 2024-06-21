package mosaic

import "math"

type (
	Transform struct {
		x     float64
		y     float64
		sin   float64
		cos   float64
		scale float64
	}
)

func NewTransform(x, y, scale, angle float64) Transform {
	if scale == 0 {
		scale = 1
	}

	if angle == 0.0 || angle == 360.0 {
		return Transform{
			scale: scale,
			sin:   0.0,
			cos:   1.0,
			x:     x,
			y:     y,
		}
	}

	r := math.Pi * angle / 180.0

	return Transform{
		scale: scale,
		sin:   math.Sin(r),
		cos:   math.Cos(r),
		x:     x,
		y:     y,
	}
}
