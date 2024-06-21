package mosaic

import "math"

type Rectangle struct {
	Position Vector
	rawEdges [4]Edge
	Edges    [4]Edge
	Planes   [4]Plane
}

func NewRectangle(position Vector, w, h float64) Rectangle {
	return Rectangle{
		Position: position,
		rawEdges: [4]Edge{
			{Start: Vector{X: -w / 2, Y: -h / 2}, End: Vector{X: -w / 2, Y: h / 2}, Active: true},
			{Start: Vector{X: -w / 2, Y: h / 2}, End: Vector{X: w / 2, Y: h / 2}, Active: true},
			{Start: Vector{X: w / 2, Y: h / 2}, End: Vector{X: w / 2, Y: -h / 2}, Active: true},
			{Start: Vector{X: w / 2, Y: -h / 2}, End: Vector{X: -w / 2, Y: -h / 2}, Active: true},
		},
	}.Update()
}

func (r Rectangle) Update() Rectangle {
	for i := 0; i < 4; i++ {
		r.Edges[i].Start = r.Position.Add(r.rawEdges[i].Start)
		r.Edges[i].End = r.Position.Add(r.rawEdges[i].End)
		r.Edges[i].Active = r.rawEdges[i].Active

		if r.Edges[i].Active {
			r.Planes[i] = NewPlane(r.Edges[i].Start, r.Edges[i].End)
		}
	}

	return r
}

func (r Rectangle) Height() float64 {
	return r.Edges[0].Start.Distance(r.Edges[0].End)
}

func (r Rectangle) Width() float64 {
	return r.Edges[1].Start.Distance(r.Edges[1].End)
}

func (r Rectangle) MinPoint() Vector {
	minPoint := r.Edges[0].Start
	for i := 1; i < 4; i++ {
		if minPoint.Length() > r.Edges[i].Start.Length() {
			minPoint = r.Edges[i].Start
		}
	}

	return minPoint
}

func (r Rectangle) MaxPoint() Vector {
	maxPoint := r.Edges[0].Start
	for i := 1; i < 4; i++ {
		if maxPoint.Length() < r.Edges[i].Start.Length() {
			maxPoint = r.Edges[i].Start
		}
	}

	return maxPoint
}

func (r Rectangle) Transform(t Transform) Rectangle {
	positionTransform := Transform{
		x:     t.x,
		y:     t.y,
		scale: 1.0,
		sin:   0.0,
		cos:   1.0,
	}
	r.Position = r.Position.Transform(positionTransform)

	edgeTransform := Transform{
		x:     0,
		y:     0,
		scale: t.scale,
		sin:   t.sin,
		cos:   t.cos,
	}

	for i := 0; i < 4; i++ {
		r.rawEdges[i] = r.rawEdges[i].Transform(edgeTransform)
	}

	return r.Update()
}

func (r Rectangle) Scale(c float64) Rectangle {
	transform := NewTransform(0, 0, c, 0)
	for i := 0; i < 4; i++ {
		r.rawEdges[i] = r.rawEdges[i].Transform(transform)
	}

	return r.Update()
}

func (r Rectangle) ContainsVector(v Vector) bool {
	rayCount := 0
	for i := 0; i < len(r.Edges); i++ {
		rayCount += r.Edges[i].RayCount(v)
	}

	return rayCount%2 == 1
}

func (r Rectangle) Intersects(s Rectangle) (normal Vector, depth float64) {
	depth = math.MaxFloat64

	for _, plane := range r.Planes {
		minP, maxP := r.projectVectors(plane.Normal)
		minQ, maxQ := s.projectVectors(plane.Normal)

		if minP >= maxQ || minQ >= maxP {
			return Vector{}, 0.0
		}

		planeDistance := math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	for _, plane := range s.Planes {
		minP, maxP := r.projectVectors(plane.Normal)
		minQ, maxQ := s.projectVectors(plane.Normal)

		if minP >= maxQ || minQ >= maxP {
			return Vector{}, 0.0
		}

		planeDistance := math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	if normal.DotProduct(s.Position.Subtract(r.Position)) < 0 {
		normal = normal.Invert()
	}

	return normal, depth
}

func (r Rectangle) projectVectors(axis Vector) (min, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, edge := range r.Edges {
		projection := edge.Start.DotProduct(axis)

		if projection < min {
			min = projection
		}
		if projection > max {
			max = projection
		}
	}

	return min, max
}

func (r Rectangle) ToPolygon() Polygon {
	return NewPolygon(
		r.Position,
		[]Vector{
			r.rawEdges[0].Start,
			r.rawEdges[1].Start,
			r.rawEdges[2].Start,
			r.rawEdges[3].Start,
		})
}

func (r Rectangle) Area() float64 {
	return r.Width() * r.Height()
}

// TODO - This does not support rotation
func (r Rectangle) AreaOfOverlap(o Rectangle) float64 {
	_, depth := r.Intersects(o)
	if depth == 0.0 {
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

func (r Rectangle) ToCircle() Circle {
	rMin := r.Position.Distance(r.MinPoint())
	rMax := r.Position.Distance(r.MaxPoint())

	if rMax > rMin {
		return NewCircle(r.Position, rMax)
	}

	return NewCircle(r.Position, rMin)
}
