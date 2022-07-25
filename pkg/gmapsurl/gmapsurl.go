package gmapsurl

import (
	"fmt"
	"strings"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
)

type VehicleType string

const (
	Car     VehicleType = "CAR"
	Bike    VehicleType = "BIKE"
	Walking VehicleType = "WALKING"

	vehicleOptsCar     = "3e0"
	vehicleOptsBike    = "3e1"
	vehicleOptsWalking = "3e2"
)

func (v VehicleType) Opts() string {
	switch v {
	case Car:
		return vehicleOptsCar
	case Bike:
		return vehicleOptsBike
	case Walking:
		return vehicleOptsWalking
	default:
		return vehicleOptsCar
	}
}

func Of(polygon geo.Polygon, vehicleType VehicleType) string {
	waypoints := make([]string, 0, len(polygon)-2)
	for _, p := range polygon {
		waypoints = append(waypoints, p.String())
	}
	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/dir/%s/data=!3m1!4b1!4m2!4m1!%s", strings.Join(waypoints, "/"), vehicleType.Opts())
	return googleMapsURL
}
