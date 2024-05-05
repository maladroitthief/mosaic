package mosaic

import (
	"fmt"
	"math"
)

type Polygon struct {
	Position    Vector
	Vectors     []Vector
	CalcVectors []Vector
	Planes      []Plane
	Bounds      Rectangle
}

// NewPolygon accepts an array of vectors in CCW rotation
func NewPolygon(position Vector, vectors []Vector) Polygon {
	p := Polygon{
		Position: position,
		Vectors:  vectors,
	}

	return p.Update()
}

func (p Polygon) Update() Polygon {
	p.CalcVectors = p.calcVectors()
	p.Planes = p.calcPlanes()
	p.Bounds = p.calcBounds()
	return p
}

func (p Polygon) Copy(q Polygon) Polygon {
	position := q.Position.Clone()
	vectors := make([]Vector, len(q.Vectors))
	for i := 0; i < len(vectors); i++ {
		vectors[i] = q.Vectors[i].Clone()
	}

	return NewPolygon(position, vectors)
}

func (p Polygon) Clone() Polygon {
	return p.Copy(p)
}

func (p Polygon) CheckPosition(position Vector) Polygon {
	return p.Clone().SetPosition(position)
}

func (p Polygon) SetPosition(position Vector) Polygon {
	if p.Position == position {
		return p
	}

	p.Position = position

	return p.Update()
}

func (p Polygon) Info() string {
	return fmt.Sprintf("%+v", p.Vectors)
}

func (p Polygon) Add(v Vector) Polygon {
	q := p.Clone()
	q.Position = q.Position.Add(v)
	return q.Update()
}

func (p Polygon) ContainsVector(v Vector) bool {
	rayCount := 0
	for i := 0; i < len(p.CalcVectors); i++ {
		j := (i + 1) % len(p.CalcVectors)
		if intersectsRay(v, p.CalcVectors[i], p.CalcVectors[j]) {
			rayCount++
		}
	}
	return rayCount%2 == 1
}

func intersectsRay(point, v1, v2 Vector) bool {
	if v1.Y == v2.Y {
		// Horizontal edge, Ray is vertical
		return false
	}

	if point.Y == v1.Y || point.Y == v2.Y {
		// Avoid edge cases
		point.Y += 0.0001
	}

	if (point.Y > v1.Y && point.Y <= v2.Y) || (point.Y > v2.Y && point.Y <= v1.Y) {
		x := v1.X + (point.Y-v1.Y)*(v2.X-v1.X)/(v2.Y-v1.Y)
		return x > point.X
	}
	return false
}

func (p Polygon) Intersects(q Polygon) (normal Vector, depth float64) {
	depth = math.MaxFloat64

	for _, plane := range p.Planes {
		minP, maxP := projectVectors(plane.Normal, p.CalcVectors)
		minQ, maxQ := projectVectors(plane.Normal, q.CalcVectors)

		if minP >= maxQ || minQ >= maxP {
			return Vector{}, 0.0
		}

		planeDistance := math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	for _, plane := range q.Planes {
		minP, maxP := projectVectors(plane.Normal, p.CalcVectors)
		minQ, maxQ := projectVectors(plane.Normal, q.CalcVectors)

		if minP >= maxQ || minQ >= maxP {
			return Vector{}, 0.0
		}

		planeDistance := math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	if normal.DotProduct(q.Position.Subtract(p.Position)) < 0 {
		normal = normal.Invert()
	}

	return normal, depth
}

func (p Polygon) ContainsPolygon(q Polygon) (normal Vector, depth float64) {
	depth = math.MaxFloat64
	contained := true

	for _, v := range q.CalcVectors {
		if !p.ContainsVector(v) {
			contained = false
			break
		}
	}

	if contained == true {
		return Vector{}, 0.0
	}

	for _, plane := range p.Planes {
		minP, maxP := projectVectors(plane.Normal, p.CalcVectors)
		minQ, maxQ := projectVectors(plane.Normal, q.CalcVectors)

		planeDistance := maxQ - minQ - math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth && planeDistance > 0 {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	for _, plane := range q.Planes {
		minP, maxP := projectVectors(plane.Normal, p.CalcVectors)
		minQ, maxQ := projectVectors(plane.Normal, q.CalcVectors)

		planeDistance := maxQ - minQ - math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < depth && planeDistance > 0 {
			depth = planeDistance
			normal = plane.Normal
		}
	}

	if normal.DotProduct(q.Position.Subtract(p.Position)) > 0 {
		normal = normal.Invert()
	}

	return normal, depth
}

func projectVectors(axis Vector, vectors []Vector) (min, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, v := range vectors {
		projection := v.DotProduct(axis)

		if projection < min {
			min = projection
		}
		if projection > max {
			max = projection
		}
	}

	return min, max
}

func (p Polygon) calcVectors() []Vector {
	vectors := make([]Vector, len(p.Vectors))
	for i := 0; i < len(vectors); i++ {
		vectors[i] = p.Position.Add(p.Vectors[i])
	}

	return vectors
}

func (p Polygon) calcPlanes() []Plane {
	planes := make([]Plane, len(p.CalcVectors))
	for i := 0; i < len(planes)-1; i++ {
		planes[i] = NewPlane(p.CalcVectors[i], p.CalcVectors[i+1])
	}
	planes[len(planes)-1] = NewPlane(p.CalcVectors[len(planes)-1], p.CalcVectors[0])

	return planes
}

func (p Polygon) calcBounds() Rectangle {
	minHeight, maxHeight := math.MaxFloat64, 0.0
	minWidth, maxWidth := math.MaxFloat64, 0.0

	for i := 0; i < len(p.CalcVectors); i++ {
		minWidth = min(minWidth, p.CalcVectors[i].X)
		maxWidth = max(maxWidth, p.CalcVectors[i].X)

		minHeight = min(minHeight, p.CalcVectors[i].Y)
		maxHeight = max(maxHeight, p.CalcVectors[i].Y)
	}

	return NewRectangle(p.Position, maxWidth-minWidth, maxHeight-minHeight)
}
