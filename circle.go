package mosaic

type (
	Circle struct {
		Position Vector
		Radius   float64
		Bounds   Rectangle
	}
)

func NewCircle(position Vector, radius float64) Circle {
	return Circle{
		Position: position,
		Radius:   radius,
		Bounds:   NewRectangle(position, radius, radius),
	}
}

func (c Circle) Update() Circle {
	c.Bounds = NewRectangle(c.Position, c.Radius, c.Radius)
	return c
}

func (c Circle) Intersects(d Circle) (normal Vector, depth float64) {
	distance := c.Position.Distance(d.Position)
	radii := c.Radius + d.Radius

	if distance >= radii {
		return Vector{}, 0.0
	}

	normal = d.Position.Subtract(c.Position).Normalize()
	depth = radii - distance

	return normal, depth
}

func (c Circle) Contains(d Circle) bool {
	return c.Radius >= c.Position.Distance(d.Position)+d.Radius
}
