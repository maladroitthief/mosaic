package mosaic

import "math"

type (
	Vector struct {
		X float64
		Y float64
	}
)

func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

func (v Vector) Copy(w Vector) Vector {
	return Vector{X: w.X, Y: w.Y}
}

func (v Vector) Clone() Vector {
	return v.Copy(v)
}

func (v Vector) Perpendicular() Vector {
	return Vector{X: v.Y, Y: -v.X}
}

func (v Vector) Invert() Vector {
	return Vector{X: -v.X, Y: -v.Y}
}

func (v Vector) DotProduct(w Vector) float64 {
	return (v.X * w.X) + (v.Y * w.Y)
}

func (v Vector) CrossProduct(w Vector) float64 {
	return (v.X * w.Y) - (v.Y * w.X)
}

func (v Vector) Normal(w Vector) Vector {
	return v.RightNormal(w)
}

// For clockwise order
func (v Vector) LeftNormal(w Vector) Vector {
	vn := v.Subtract(w).Normalize()
	return Vector{
		X: -vn.Y,
		Y: vn.X,
	}
}

// For counter clockwise order
func (v Vector) RightNormal(w Vector) Vector {
	vn := w.Subtract(v).Normalize()
	return Vector{
		X: vn.Y,
		Y: -vn.X,
	}
}

func (v Vector) Add(w Vector) Vector {
	return Vector{
		X: v.X + w.X,
		Y: v.Y + w.Y,
	}
}

func (v Vector) Subtract(w Vector) Vector {
	return Vector{
		X: v.X - w.X,
		Y: v.Y - w.Y,
	}
}

func (v Vector) Scale(c float64) Vector {
	return Vector{
		X: v.X * c,
		Y: v.Y * c,
	}
}

func (v Vector) ScaleXY(cx, cy float64) Vector {
	return Vector{
		X: v.X * cx,
		Y: v.Y * cy,
	}
}

func (v Vector) Projection(w Vector) Vector {
	return w.Scale(v.DotProduct(w) / w.DotProduct(w))
}

func (v Vector) UnitProjection(w Vector) Vector {
	return w.Scale(v.DotProduct(w))
}

func (v Vector) Reflect(w Vector) Vector {
	return v.Projection(w).Scale(2).Subtract(v)
}

func (v Vector) UnitReflect(w Vector) Vector {
	return v.UnitProjection(w).Scale(2).Subtract(v)
}

func (v Vector) Normalize() Vector {
	c := v.Magnitude()
	if c == 0 {
		c = 1
	}

	return v.Scale(1 / c)
}

func (v Vector) Length() float64 {
	return v.DotProduct(v)
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Length())
}

func (v Vector) Distance(w Vector) float64 {
	return math.Sqrt(math.Pow(w.X-v.X, 2) + math.Pow(w.Y-v.Y, 2))
}

func (v Vector) Transform(t Transform) Vector {
	return Vector{
		X: t.scale*(t.cos*v.X-t.sin*v.Y) + t.x,
		Y: t.scale*(t.sin*v.X+t.cos*v.Y) + t.y,
	}
}
