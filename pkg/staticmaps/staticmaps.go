package staticmaps

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"golang.org/x/image/font/basicfont"

	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
)

func RenderOnMap(polygon geo.Polygon) (image.Image, error) {
	cont := sm.NewContext()
	cont.SetSize(1920, 1080)

	positions := make([]s2.LatLng, 0, len(polygon))
	for _, p := range polygon {
		positions = append(positions, s2.LatLngFromDegrees(p.Lat, p.Lng))
	}

	cont.AddObject(sm.NewPath(positions, color.Black, 1))

	for i, p := range positions {
		cont.AddObject(sm.NewImageMarker(p, numberTextBox(i), 0.5, 0.5))
	}

	img, err := cont.Render()
	if err != nil {
		return nil, fmt.Errorf("rendering image: %w", err)
	}

	return img, nil
}

func numberTextBox(num int) image.Image {
	width, height := 20.0, 20.0

	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(color.Black)
	dc.Clear()
	dc.SetFontFace(basicfont.Face7x13)

	dc.SetColor(color.White)
	dc.DrawStringAnchored(strconv.Itoa(num), height/2, width/2, 0.5, 0.5)

	return dc.Image()
}
