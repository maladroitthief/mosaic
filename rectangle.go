package mosaic

type Rectangle struct {
	Position Vector
	Height   float64
	Width    float64
}

func NewRectangle(position Vector, w, h float64) Rectangle {
	return Rectangle{
		Position: position,
		Width:    w,
		Height:   h,
	}
}

func (r Rectangle) MinPoint() Vector {
	return Vector{X: r.Position.X - r.Width/2, Y: r.Position.Y - r.Height/2}
}

func (r Rectangle) MaxPoint() Vector {
	return Vector{X: r.Position.X + r.Width/2, Y: r.Position.Y + r.Height/2}
}

func (r Rectangle) Contains(x, y float64) bool {
	if x < r.MinPoint().X || x > r.MaxPoint().X {
		return false
	}

	if y < r.MinPoint().Y || y > r.MaxPoint().Y {
		return false
	}

	return true
}

func (r Rectangle) Intersects(s Rectangle) bool {
	d1x := s.MinPoint().X - r.MaxPoint().X
	d1y := s.MinPoint().Y - r.MaxPoint().Y
	d2x := r.MinPoint().X - s.MaxPoint().X
	d2y := r.MinPoint().Y - s.MaxPoint().Y

	if d1x > 0.0 || d1y > 0.0 {
		return false
	}

	if d2x > 0.0 || d2y > 0.0 {
		return false
	}

	return true
}

func (r Rectangle) Scale(c float64) Rectangle {
	width := r.Width * c
	height := r.Height * c

	return NewRectangle(r.Position, height, width)
}

func (r Rectangle) ToPolygon() Polygon {
	return NewPolygon(
		r.Position,
		[]Vector{
			{X: -r.Width / 2, Y: -r.Height / 2},
			{X: -r.Width / 2, Y: r.Height / 2},
			{X: r.Width / 2, Y: r.Height / 2},
			{X: r.Width / 2, Y: -r.Height / 2},
		})
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) AreaOfOverlap(o Rectangle) float64 {
	if r.Intersects(o) == false {
		return 0.0
	}

	rMin := r.MinPoint()
	rMax := r.MaxPoint()

	oMin := o.MinPoint()
	oMax := o.MaxPoint()

	x := min(rMax.X, oMax.X) - max(rMin.X, oMin.X)
	y := min(rMax.Y, oMax.Y) - max(rMin.Y, oMin.Y)

	return x * y
}
