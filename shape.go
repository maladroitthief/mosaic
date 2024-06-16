package mosaic

type (
	ShapeType int
	Shape     interface {
		Bounds() Rectangle
		Type()
	}
)

const (
	PolygonShape ShapeType = iota
	CircleShape
	RectangleShape
	TriangleShape
)
