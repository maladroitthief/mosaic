package mosaic_test

import "math"

func WithinTolerance(x, y, tolerance float64) bool {
	if x == y {
		return true
	}

	delta := x - y
	if y == 0 {
		return delta < tolerance
	}

	return (delta / math.Abs(y)) < tolerance
}
