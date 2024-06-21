package mosaic_test

import (
	"testing"

	"github.com/maladroitthief/mosaic"
)

func Test_edge_XIntersect(t *testing.T) {
	type setup struct {
		edge mosaic.Edge
	}
	type input struct {
		edge mosaic.Edge
	}
	type want struct {
		intersect float64
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
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(300, 200),
					End:    mosaic.NewVector(100, 150),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(200, 150),
					End:    mosaic.NewVector(200, 200),
					Active: true,
				},
			},
			want: want{
				intersect: 200,
			},
		},
		{
			name: "angle case",
			setup: setup{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(5, 5),
					End:    mosaic.NewVector(25, 25),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(5, 25),
					End:    mosaic.NewVector(25, 5),
					Active: true,
				},
			},
			want: want{
				intersect: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.setup.edge.XIntersect(tt.input.edge)

			if got != tt.want.intersect {
				t.Errorf("edge.XIntersect() got = %v, want %v", got, tt.want.intersect)
			}
		})
	}
}

func Test_edge_YIntersect(t *testing.T) {
	type setup struct {
		edge mosaic.Edge
	}
	type input struct {
		edge mosaic.Edge
	}
	type want struct {
		intersect float64
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
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(300, 200),
					End:    mosaic.NewVector(100, 150),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(200, 150),
					End:    mosaic.NewVector(200, 200),
					Active: true,
				},
			},
			want: want{
				intersect: 175.0,
			},
		},
		{
			name: "reverse case",
			setup: setup{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(100, 150),
					End:    mosaic.NewVector(300, 200),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(200, 150),
					End:    mosaic.NewVector(200, 200),
					Active: true,
				},
			},
			want: want{
				intersect: 175.0,
			},
		},
		{
			name: "reverse case",
			setup: setup{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(100, 150),
					End:    mosaic.NewVector(300, 200),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(150, 150),
					End:    mosaic.NewVector(150, 200),
					Active: true,
				},
			},
			want: want{
				intersect: 162.5,
			},
		},
		{
			name: "angle case",
			setup: setup{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(5, 5),
					End:    mosaic.NewVector(25, 25),
					Active: true,
				},
			},
			input: input{
				edge: mosaic.Edge{
					Start:  mosaic.NewVector(5, 25),
					End:    mosaic.NewVector(25, 5),
					Active: true,
				},
			},
			want: want{
				intersect: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.setup.edge.YIntersect(tt.input.edge)

			if got != tt.want.intersect {
				t.Errorf("edge.YIntersect() got = %v, want %v", got, tt.want.intersect)
			}
		})
	}
}
