package mosaic_test

import (
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

func Test_polygon_SetEdge(t *testing.T) {
	type setup struct {
		polygon mosaic.Polygon
		start   mosaic.Vector
		end     mosaic.Vector
		active  bool
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
				polygon: mosaic.NewRectangle(mosaic.NewVector(10, 10), 10, 10).ToPolygon(),
				start:   mosaic.Vector{5, -5},
				end:     mosaic.Vector{-5, -5},
				active:  false,
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(10, 5), 1, 1).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{},
				depth:  0.0,
			},
		},
		{
			name: "out of bounds in X",
			setup: setup{
				polygon: mosaic.NewRectangle(mosaic.NewVector(10, 10), 10, 10).ToPolygon(),
				start:   mosaic.Vector{5, -5},
				end:     mosaic.Vector{5, 5},
				active:  false,
			},
			input: input{
				polygon: mosaic.NewRectangle(mosaic.NewVector(5, 10), 1, 1).ToPolygon(),
			},
			want: want{
				normal: mosaic.Vector{-1, 0},
				depth:  2.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup.polygon = tt.setup.polygon.SetEdge(tt.setup.start, tt.setup.end, tt.setup.active)

			got := struct {
				normal mosaic.Vector
				depth  float64
			}{}

			got.normal, got.depth = tt.setup.polygon.Intersects(tt.input.polygon)

			if got.normal != tt.want.normal {
				t.Errorf("polygon.ContainsPolygon() normal = %v, want %v", got.normal, tt.want.normal)
			}

			if got.depth != tt.want.depth {
				t.Errorf("polygon.ContainsPolygon() depth = %v, want %v", got.depth, tt.want.depth)
			}
		})
	}

}
