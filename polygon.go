package mosaic

import (
	"math"
)

type (
	Edge struct {
		start  Vector
		end    Vector
		active bool
	}
	Polygon struct {
		Position Vector
		rawEdges []Edge
		Edges    []Edge
		Planes   []Plane
		Bounds   Rectangle
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
			start:  vectors[i],
			end:    vectors[(i+1)%len(vectors)],
			active: true,
		}
	}

	return p.Update()
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
		vectors[i] = q.rawEdges[i].start.Clone()
	}

	return NewPolygon(position, vectors)
}

func (p Polygon) Clone() Polygon {
	return p.Copy(p)
}

func (p Polygon) SetEdge(start, end Vector, active bool) Polygon {
	for i := range p.rawEdges {
		if p.rawEdges[i].start == start && p.rawEdges[i].end == end {
			p.rawEdges[i].active = active
			p.Edges[i].active = active
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
		start := p.Edges[i].start
		end := p.Edges[i].end

		if start.Y == end.Y {
			// Horizontal edge, Ray is vertical
			continue
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
	}

	return rayCount%2 == 1
}

func (p Polygon) Intersects(q Polygon) (normal Vector, depth float64) {
	depth = math.MaxFloat64

	for _, plane := range p.Planes {
		minP, maxP := projectVectors(plane.Normal, p.Edges)
		minQ, maxQ := projectVectors(plane.Normal, q.Edges)

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
		minP, maxP := projectVectors(plane.Normal, p.Edges)
		minQ, maxQ := projectVectors(plane.Normal, q.Edges)

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
		if !p.ContainsVector(v.start) {
			contained = false
			break
		}
	}

	if contained == true {
		return Vector{}, 0.0
	}

	for _, plane := range p.Planes {
		minP, maxP := projectVectors(plane.Normal, p.Edges)
		minQ, maxQ := projectVectors(plane.Normal, q.Edges)

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
		minP, maxP := projectVectors(plane.Normal, p.Edges)
		minQ, maxQ := projectVectors(plane.Normal, q.Edges)

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

func projectVectors(axis Vector, edges []Edge) (min, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, edge := range edges {
		projection := edge.start.DotProduct(axis)

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
		p.Edges[i].start = p.Position.Add(p.rawEdges[i].start)
		p.Edges[i].end = p.Position.Add(p.rawEdges[i].end)
		p.Edges[i].active = p.rawEdges[i].active
	}

	return p.Edges
}

func (p Polygon) calcPlanes() []Plane {
	planes := make([]Plane, len(p.Edges))
	for i := 0; i < len(planes); i++ {
		if p.Edges[i].active {
			planes[i] = NewPlane(p.Edges[i].start, p.Edges[i].end)
		}
	}

	return planes
}

func (p Polygon) calcBounds() Rectangle {
	minHeight, maxHeight := math.MaxFloat64, 0.0
	minWidth, maxWidth := math.MaxFloat64, 0.0

	for i := 0; i < len(p.Edges); i++ {
		minWidth = min(minWidth, p.Edges[i].start.X)
		maxWidth = max(maxWidth, p.Edges[i].start.X)

		minHeight = min(minHeight, p.Edges[i].start.Y)
		maxHeight = max(maxHeight, p.Edges[i].start.Y)
	}

	return NewRectangle(p.Position, maxWidth-minWidth, maxHeight-minHeight)
}
