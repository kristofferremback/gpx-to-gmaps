package gmapsurl

import (
	"fmt"
	"strings"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
)

func Of(polygon geo.Polygon) string {
	waypoints := make([]string, 0, len(polygon)-2)
	for _, p := range polygon {
		waypoints = append(waypoints, p.String())
	}

	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/dir/%s/data=!3m1!4b1!4m2!4m1!3e1", strings.Join(waypoints, "/"))
	return googleMapsURL
}
