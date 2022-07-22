package gpx

import (
	"fmt"
	"os"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
	"github.com/tkrajina/gpxgo/gpx"
)

func ReadGPXFile(filename string) (*gpx.GPX, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	data, err := gpx.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing gpx data: %w", err)
	}
	return data, nil
}

func PolygonsOf(data *gpx.GPX) ([]geo.Polygon, error) {
	polygons := make([]geo.Polygon, 0, len(data.Tracks))
	for _, track := range data.Tracks {
		polygon := PolygonOf(track)
		blurred := geo.ReduceSize(polygon, 25)

		polygons = append(polygons, blurred)
	}

	return polygons, nil
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
