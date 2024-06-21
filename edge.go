package mosaic

type (
	Edge struct {
		Start  Vector
		End    Vector
		Active bool
	}
)

func (e Edge) Transform(t Transform) Edge {
	return Edge{
		Start:  e.Start.Transform(t),
		End:    e.End.Transform(t),
		Active: e.Active,
	}
}

func (e Edge) RayCount(v Vector) int {
	rayCount := 0
	start := e.Start
	end := e.End

	if start.Y == end.Y {
		// Horizontal edge, Ray is vertical
		return rayCount
	}

	if v.Y == start.Y || v.Y == end.Y {
		// Avoid edge cases
		v.Y += 0.0001
	}

	if (v.Y > start.Y && v.Y <= end.Y) || (v.Y > end.Y && v.Y <= start.Y) {
		x := start.X + (v.Y-start.Y)*(end.X-start.X)/(end.Y-start.Y)
		if x > v.X {
			rayCount++
		}
	}

	return rayCount
}

func (e Edge) XIntersect(f Edge) float64 {
	xNumerator := (e.Start.X*e.End.Y-e.Start.Y*e.End.X)*(f.Start.X-f.End.X) -
		(e.Start.X-e.End.X)*(f.Start.X*f.End.Y-f.Start.Y*f.End.X)
	denominator := (e.Start.X-e.End.X)*(f.Start.Y-f.End.Y) -
		(e.Start.Y-e.End.Y)*(f.Start.X-f.End.X)

	return xNumerator / denominator
}

func (e Edge) YIntersect(f Edge) float64 {
	yNumerator := (e.Start.X*e.End.Y-e.Start.Y*e.End.X)*(f.Start.Y-f.End.Y) -
		(e.Start.Y-e.End.Y)*(f.Start.X*f.End.Y-f.Start.Y*f.End.X)
	denominator := (e.Start.X-e.End.X)*(f.Start.Y-f.End.Y) -
		(e.Start.Y-e.End.Y)*(f.Start.X-f.End.X)

	return yNumerator / denominator
}

func (e Edge) Intersect(f Edge) Vector {
	xNumerator := (e.Start.X*e.End.Y-e.Start.Y*e.End.X)*(f.Start.X-f.End.X) -
		(e.Start.X-e.End.X)*(f.Start.X*f.End.Y-f.Start.Y*f.End.X)
	yNumerator := (e.Start.X*e.End.Y-e.Start.Y*e.End.X)*(f.Start.Y-f.End.Y) -
		(e.Start.Y-e.End.Y)*(f.Start.X*f.End.Y-f.Start.Y*f.End.X)
	denominator := (e.Start.X-e.End.X)*(f.Start.Y-f.End.Y) -
		(e.Start.Y-e.End.Y)*(f.Start.X-f.End.X)

	return Vector{X: xNumerator / denominator, Y: yNumerator / denominator}
}

// Assuming CCW orientation
func (e Edge) ContainsVector(v Vector) bool {
	return (e.End.X-e.Start.X)*(v.Y-e.Start.Y) >
		(e.End.Y-e.Start.Y)*(v.X-e.Start.X)
}
