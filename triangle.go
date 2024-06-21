package mosaic

type (
	Triangle struct {
		Position Vector
		Rotation float64
		rawEdges [3]Edge
		Edges    [3]Edge
	}
)
