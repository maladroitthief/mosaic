package mosaic_test

import (
	"testing"

	"github.com/maladroitthief/mosaic"
)

func Test_vector_Magnitude(t *testing.T) {
	type fields struct {
		X float64
		Y float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "base case",
			fields: fields{X: 5, Y: 5},
			want:   7.0710678118654755,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := mosaic.Vector{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			got := v.Magnitude()
			if got != tt.want {
				t.Errorf("vector.Magnitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vector_Transform(t *testing.T) {
	type (
		setup struct {
			X float64
			Y float64
		}
		params struct {
			transform mosaic.Transform
		}
	)

	tests := []struct {
		name   string
		setup  setup
		params params
		want   mosaic.Vector
	}{
		{
			name: "do nothing",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					0,
					0,
					1,
					0,
				),
			},
			want: mosaic.Vector{
				X: 5,
				Y: 5,
			},
		},
		{
			name: "move",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					1,
					1,
					1,
					0,
				),
			},
			want: mosaic.Vector{
				X: 6,
				Y: 6,
			},
		},
		{
			name: "scale",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					0,
					0,
					2,
					0,
				),
			},
			want: mosaic.Vector{
				X: 10,
				Y: 10,
			},
		},
		{
			name: "rotate",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					0,
					0,
					1,
					90,
				),
			},
			want: mosaic.Vector{
				X: -5,
				Y: 5,
			},
		},
		{
			name: "move and scale",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					5,
					6,
					2,
					0,
				),
			},
			want: mosaic.Vector{
				X: 15,
				Y: 16,
			},
		},
		{
			name: "move and rotate",
			setup: setup{
				X: 5,
				Y: 5,
			},
			params: params{
				transform: mosaic.NewTransform(
					5,
					6,
					1,
					90,
				),
			},
			want: mosaic.Vector{
				X: 0,
				Y: 11,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := mosaic.Vector{
				X: tt.setup.X,
				Y: tt.setup.Y,
			}

			got := v.Transform(tt.params.transform)
			if got != tt.want {
				t.Errorf("vector.Transform() = %v, want %v", got, tt.want)
			}
		})
	}

}
