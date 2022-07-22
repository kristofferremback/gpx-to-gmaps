package geo

import (
	"fmt"
	"math"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/slices"
)

type Polygon []Point

type Point struct {
	Lat, Lng float64
}

func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.Lat, p.Lng)
}

func ReduceSize(polygon Polygon, maxSize int) Polygon {
	if len(polygon) <= maxSize {
		return polygon
	}

	size := 100
	for {
		out := smoothenStraightishLines(slices.PickSpaced(polygon, size))
		if len(out) <= maxSize || maxSize == 0 {
			return out
		}
		size--
	}
}

func smoothenStraightishLines(polygon Polygon) Polygon {
	out := make(Polygon, 0)
	for i, p := range polygon {
		if i == 0 || i == len(polygon)-1 {
			out = append(out, p)
			continue
		}

		prev := out[len(out)-1]
		next := polygon[i+1]

		angle := math.Atan2(p.Lat-prev.Lat, p.Lng-prev.Lng) * (180 / math.Pi)
		angleNext := math.Atan2(next.Lat-p.Lat, next.Lng-p.Lng) * (180 / math.Pi)

		if math.Abs(angle-angleNext) >= 45 {
			out = append(out, p)
		}
	}
	return out
}
