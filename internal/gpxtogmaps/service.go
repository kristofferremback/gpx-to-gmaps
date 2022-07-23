package gpxtogmaps

import (
	"fmt"
	"image/png"
	"io"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/gmapsurl"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/gpx"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/slices"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/staticmaps"
	"github.com/twpayne/go-polyline"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ConvertToPolygons(gpxReader io.Reader, maxSize int) ([]geo.Polygon, error) {
	g, err := gpx.Parse(gpxReader)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	polygons := slices.Map(gpx.PolygonsOf(g), func(polygon geo.Polygon) geo.Polygon {
		return geo.ReduceSize(polygon, 25)
	})

	return polygons, nil
}

func (s *Service) GoogleMapsURL(polygon geo.Polygon) string {
	return gmapsurl.Of(polygon)
}

func (s *Service) PNG(writer io.Writer, polygon geo.Polygon) error {
	img, err := staticmaps.RenderOnMap(polygon)
	if err != nil {
		return fmt.Errorf("rendering image on map: %w", err)
	}

	if err := png.Encode(writer, img); err != nil {
		return fmt.Errorf("encoding png: %w", err)
	}
	return nil
}

func (s *Service) EncodePolyline(polygon geo.Polygon) string {
	pl := polyline.EncodeCoords(slices.Map(polygon, func(p geo.Point) []float64 {
		return []float64{p.Lat, p.Lng}
	}))

	return string(pl)
}

func (s *Service) DecodePolyline(pl string) (geo.Polygon, error) {
	poly, extra, err := polyline.DecodeCoords([]byte(pl))
	if err != nil {
		return nil, fmt.Errorf("decoding polyline: %w", err)
	}
	if len(extra) > 0 {
		return nil, fmt.Errorf("unexpected remainder bits in polyline: %s", extra)
	}
	return slices.Map(poly, func(coord []float64) geo.Point {
		return geo.Point{Lat: coord[0], Lng: coord[1]}
	}), nil
}
