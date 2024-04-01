package mosaic_test

import (
	"testing"

	"github.com/maladroitthief/mosaic"
)

func Test_rectangle_AreaOfOverlap(t *testing.T) {
	type fields struct {
		R mosaic.Rectangle
		O mosaic.Rectangle
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "base case",
			fields: fields{
				R: mosaic.NewRectangle(
					mosaic.NewVector(1, 1),
					4,
					4,
				),
				O: mosaic.NewRectangle(
					mosaic.NewVector(1, 1),
					4,
					4,
				),
			},
			want: 16,
		},
		{
			name: "no overlap",
			fields: fields{
				R: mosaic.NewRectangle(
					mosaic.NewVector(2, 2),
					4,
					4,
				),
				O: mosaic.NewRectangle(
					mosaic.NewVector(10, 10),
					4,
					4,
				),
			},
			want: 0,
		},
		{
			name: "half overlap",
			fields: fields{
				R: mosaic.NewRectangle(
					mosaic.NewVector(2, 2),
					4,
					4,
				),
				O: mosaic.NewRectangle(
					mosaic.NewVector(2, 4),
					4,
					4,
				),
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.R.AreaOfOverlap(tt.fields.O)
			if got != tt.want {
				t.Errorf("rectangle.AreaOfOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}
