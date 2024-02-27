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
