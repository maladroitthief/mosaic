package mosaic_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/maladroitthief/mosaic"
)

func Test_polygon_ContainsPolygon(t *testing.T) {
	type setup struct {
		polygon mosaic.Polygon
	}
	type input struct {
		polygon mosaic.Polygon
	}
	type want struct {
		normal mosaic.Vector
		depth  float64
	}
	tests := []struct {
		name  string
		setup setup
		input input
		want  want
	}{
		{
			name: "base case",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 5), 10, 10).ToPolygon(),
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 5), 5, 5).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{},
				depth:  0.0,
			},
		},
		{
			name: "out of bounds in X",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 5), 10, 10).ToPolygon(),
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(11, 5), 2, 2).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{-1, 0},
				depth:  2.0,
			},
		},
		{
			name: "out of bounds in Y",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 5), 10, 10).ToPolygon(),
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 11), 2, 2).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{0, -1},
				depth:  2.0,
			},
		},
		{
			name: "out of bounds in XY",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 5), 10, 10).ToPolygon(),
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(12, 12), 2, 2).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{-0.7071067811865475, -0.7071067811865475},
				depth:  4.242640687119285,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := struct {
				normal mosaic.Vector
				depth  float64
			}{}
			got.normal, got.depth = tt.setup.polygon.ContainsPolygon(tt.input.polygon)

			if got.normal != tt.want.normal {
				t.Errorf("polygon.ContainsPolygon() normal = %v, want %v", got.normal, tt.want.normal)
			}

			if got.depth != tt.want.depth {
				t.Errorf("polygon.ContainsPolygon() depth = %v, want %v", got.depth, tt.want.depth)
			}
		})
	}

}

func Test_polygon_Join(t *testing.T) {
	type setup struct {
		polygon mosaic.Polygon
	}
	type input struct {
		polygon mosaic.Polygon
	}
	type want struct {
		p []mosaic.Plane
		q []mosaic.Plane
	}
	tests := []struct {
		name  string
		setup setup
		input input
		want  want
	}{
		{
			name: "base case",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(2, 2), 4, 4).ToPolygon(),
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(6, 2), 4, 4).ToPolygon(),
			},
			want: want{
				p: []mosaic.Plane{
					{mosaic.Vector{1, -0}, 0}, {mosaic.Vector{0, -1}, -4}, {mosaic.Vector{0, 1}, 0},
				},
				q: []mosaic.Plane{
					{mosaic.Vector{0, -1}, -4}, {mosaic.Vector{-1, -0}, -8}, {mosaic.Vector{0, 1}, 0},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := struct {
				p mosaic.Polygon
				q mosaic.Polygon
			}{}
			got.p, got.q = tt.setup.polygon.Join(tt.input.polygon)

			comparePlanes := func(p mosaic.Plane, q mosaic.Plane) int {
				x := cmp.Compare(p.Distance, q.Distance)
				y := cmp.Compare(p.Normal.X, q.Normal.X) * 10
				z := cmp.Compare(p.Normal.Y, q.Normal.Y) * 100

				return x + y + z
			}

			if slices.CompareFunc(got.p.Planes, tt.want.p, comparePlanes) != 0 {
				t.Errorf("polygon.Join() p.Planes = %v, want %v", got.p.Planes, tt.want.p)
			}

			if slices.CompareFunc(got.q.Planes, tt.want.q, comparePlanes) != 0 {
				t.Errorf("polygon.Join() q.Planes = %v, want %v", got.q.Planes, tt.want.q)
			}
		})
	}

}
