package gpx

import (
	"fmt"
	"io"
	"os"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
	"github.com/tkrajina/gpxgo/gpx"
)

func Parse(reader io.Reader) (*gpx.GPX, error) {
	data, err := gpx.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("parsing gpx data: %w", err)
	}
	return data, nil
}

func ReadGPXFile(filename string) (*gpx.GPX, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

func PolygonsOf(data *gpx.GPX) []geo.Polygon {
	polygons := make([]geo.Polygon, 0, len(data.Tracks))
	for _, track := range data.Tracks {
		polygon := PolygonOf(track)

		polygons = append(polygons, polygon)
	}

	return polygons
}

func PolygonOf(track gpx.GPXTrack) geo.Polygon {
	polygon := make(geo.Polygon, 0)
	for _, segment := range track.Segments {
		for _, p := range segment.Points {
			polygon = append(polygon, geo.Point{Lat: p.Latitude, Lng: p.Longitude})
		}
	}
	return polygon
}
