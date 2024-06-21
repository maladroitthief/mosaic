package mosaic

import (
	"fmt"
	"math"
)

type (
	Polygon struct {
		Position Vector
		rawEdges []Edge
		Edges    []Edge
		Planes   []Plane
		Bounds   Rectangle
		Rotation float64
	}
)

// NewPolygon accepts an array of vectors in CCW rotation
func NewPolygon(position Vector, vectors []Vector) Polygon {
	p := Polygon{
		Position: position,
		rawEdges: make([]Edge, len(vectors)),
		Edges:    make([]Edge, len(vectors)),
	}

	for i := 0; i < len(vectors); i++ {
		p.rawEdges[i] = Edge{
			Start:  vectors[i],
			End:    vectors[(i+1)%len(vectors)],
			Active: true,
		}
	}

	return p.Update()
}

func (p Polygon) Info() string {
	return fmt.Sprintf("%+v, %+v", p.Position, p.Edges)
}

func (p Polygon) Update() Polygon {
	p.Edges = p.calcEdges()
	p.Planes = p.calcPlanes()
	p.Bounds = p.calcBounds()
	return p
}

func (p Polygon) Copy(q Polygon) Polygon {
	position := q.Position.Clone()
	vectors := make([]Vector, len(q.rawEdges))
	for i := 0; i < len(vectors); i++ {
		vectors[i] = q.rawEdges[i].Start.Clone()
	}

	return NewPolygon(position, vectors)
}

func (p Polygon) Clone() Polygon {
	return p.Copy(p)
}

func (p Polygon) SetEdge(start, end Vector, active bool) Polygon {
	for i := range p.rawEdges {
		if p.rawEdges[i].Start == start && p.rawEdges[i].End == end {
			p.rawEdges[i].Active = active
			p.Edges[i].Active = active
		}
	}

	return p.Update()
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

func (p Polygon) Add(v Vector) Polygon {
	q := p.Clone()
	q.Position = q.Position.Add(v)
	return q.Update()
}

func (p Polygon) ContainsVector(v Vector) bool {
	rayCount := 0
	for i := 0; i < len(p.Edges); i++ {
		rayCount += p.Edges[i].RayCount(v)
	}

	return rayCount%2 == 1
}

func (p Polygon) Intersects(q Polygon) (normal Vector, depth float64) {
	depth = math.MaxFloat64

	for _, plane := range p.Planes {
		minP, maxP := p.projectVectors(plane.Normal)
		minQ, maxQ := q.projectVectors(plane.Normal)

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
		minP, maxP := p.projectVectors(plane.Normal)
		minQ, maxQ := q.projectVectors(plane.Normal)

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
	xDepth := math.MaxFloat64
	xNormal := Vector{}
	yDepth := math.MaxFloat64
	yNormal := Vector{}
	contained := true

	for _, v := range q.Edges {
		if !p.ContainsVector(v.Start) {
			contained = false
			break
		}
	}

	if contained == true {
		return Vector{}, 0.0
	}

	for _, plane := range p.Planes {
		minP, maxP := p.projectVectors(plane.Normal)
		minQ, maxQ := q.projectVectors(plane.Normal)

		planeDistance := maxQ - minQ - math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < xDepth && planeDistance > 0 && plane.Normal.X != 0 {
			xDepth = planeDistance
			xNormal = plane.Normal
		}

		if planeDistance < yDepth && planeDistance > 0 && plane.Normal.Y != 0 {
			yDepth = planeDistance
			yNormal = plane.Normal
		}
	}

	for _, plane := range q.Planes {
		minP, maxP := p.projectVectors(plane.Normal)
		minQ, maxQ := q.projectVectors(plane.Normal)

		planeDistance := maxQ - minQ - math.Min(maxQ-minP, maxP-minQ)
		if planeDistance < xDepth && planeDistance > 0 && plane.Normal.X != 0 {
			xDepth = planeDistance
			xNormal = plane.Normal
		}

		if planeDistance < yDepth && planeDistance > 0 && plane.Normal.Y != 0 {
			yDepth = planeDistance
			yNormal = plane.Normal
		}
	}

	if xNormal.DotProduct(q.Position.Subtract(p.Position)) > 0 {
		xNormal = xNormal.Invert()
	}

	if yNormal.DotProduct(q.Position.Subtract(p.Position)) > 0 {
		yNormal = yNormal.Invert()
	}

	if xDepth == math.MaxFloat64 {
		normal = yNormal
		depth = yDepth
	} else if yDepth == math.MaxFloat64 {
		normal = xNormal
		depth = xDepth
	} else {
		normal = xNormal.Add(yNormal).Normalize()
		depth = math.Sqrt(math.Pow(xDepth, 2) + math.Pow(yDepth, 2))
	}

	return normal, depth
}

func (p Polygon) projectVectors(axis Vector) (min, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, edge := range p.Edges {
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

func (p Polygon) calcEdges() []Edge {
	for i := 0; i < len(p.rawEdges); i++ {
		p.Edges[i].Start = p.Position.Add(p.rawEdges[i].Start)
		p.Edges[i].End = p.Position.Add(p.rawEdges[i].End)
		p.Edges[i].Active = p.rawEdges[i].Active
	}

	return p.Edges
}

func (p Polygon) calcPlanes() []Plane {
	planes := make([]Plane, len(p.Edges))
	for i := 0; i < len(planes); i++ {
		if p.Edges[i].Active {
			planes[i] = NewPlane(p.Edges[i].Start, p.Edges[i].End)
		}
	}

	return planes
}

func (p Polygon) calcBounds() Rectangle {
	minHeight, maxHeight := math.MaxFloat64, 0.0
	minWidth, maxWidth := math.MaxFloat64, 0.0

	for i := 0; i < len(p.Edges); i++ {
		minWidth = min(minWidth, p.Edges[i].Start.X)
		maxWidth = max(maxWidth, p.Edges[i].Start.X)

		minHeight = min(minHeight, p.Edges[i].Start.Y)
		maxHeight = max(maxHeight, p.Edges[i].Start.Y)
	}

	return NewRectangle(p.Position, maxWidth-minWidth, maxHeight-minHeight)
}

// Gauss's shoelace formula
func (p Polygon) Area() float64 {
	area := 0.0

	for i := 0; i < len(p.Edges); i++ {
		area += p.Edges[i].Start.X*p.Edges[i].End.Y - p.Edges[i].Start.Y*p.Edges[i].End.X
	}

	return math.Abs(area)
}

func (p Polygon) Clip(clip Polygon) Polygon {
	subject := p.Clone()
	for i := 0; i < len(clip.Edges); i++ {
		vectors := make([]Vector, 0, len(p.Edges))
		for j := 0; j < len(subject.Edges); j++ {
			start := clip.Edges[i].ContainsVector(subject.Edges[j].Start)
			end := clip.Edges[i].ContainsVector(subject.Edges[j].End)

			switch {
			case start && end:
				vectors = append(vectors, NewVector(subject.Edges[j].End.X, subject.Edges[j].End.Y))
			case !start && end:
				vectors = append(vectors, NewVector(
					subject.Edges[j].XIntersect(clip.Edges[i]),
					subject.Edges[j].YIntersect(clip.Edges[i]),
				))
				vectors = append(vectors, NewVector(subject.Edges[j].End.X, subject.Edges[j].End.Y))
			case start && !end:
				vectors = append(vectors, NewVector(
					subject.Edges[j].XIntersect(clip.Edges[i]),
					subject.Edges[j].YIntersect(clip.Edges[i]),
				))
			default:
			}
		}

		rawVectors := make([]Vector, len(vectors))
		for k := 0; k < len(vectors); k++ {
			rawVectors[k] = p.Position.Subtract(vectors[k])
		}

		subject = NewPolygon(p.Position, rawVectors)
	}

	return subject
}
