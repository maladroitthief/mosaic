package mosaic

type (
	ShapeType int
	Shape     interface {
		Bounds() Rectangle
		Type()
	}
)

const (
	CircleShape ShapeType = iota
	TriangleShape
	RectangleShape
	PolygonShape
)
