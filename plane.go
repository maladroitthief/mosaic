package mosaic

type Plane struct {
	Normal Vector
	// Distance from the origin
	Distance float64
}

func NewPlane(v, w Vector) Plane {
	normal := v.RightNormal(w)
	distance := normal.DotProduct(w)

	return Plane{
		Normal:   normal,
		Distance: distance,
	}
}

func (p Plane) DistanceTo(v Vector) float64 {
	return p.Normal.DotProduct(v) - p.Distance
}

func (p Plane) Invert() Plane {
	p.Normal.Scale(-1)
	p.Distance = p.Distance * -1
	return p
}
