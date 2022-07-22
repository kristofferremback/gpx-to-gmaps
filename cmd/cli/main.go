package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tkrajina/gpxgo/gpx"
	"golang.org/x/image/font/basicfont"
)

var (
	filenameF     = flag.String("filename", "", "name of gpx file")
	outputFolderF = flag.String("output-dir", "", "folder to output the maps in")
)

type Point struct {
	Lat, Lng float64
}

func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.Lat, p.Lng)
}

func main() {
	flag.Parse()
	gpxData, err := gpxDataOf(*filenameF)
	if err != nil {
		log.Fatalf("getting gpx data from %s: %v", *filenameF, err)
	}
	polygons, err := polygonsOf(gpxData)
	if err != nil {
		log.Fatalf("getting polygons %v", err)
	}

	for _, polygon := range polygons {
		fmt.Println(googleMapsURLOf(polygon))
	}

	if outputFolderF != nil && *outputFolderF != "" {
		if err := renderPNGs(*outputFolderF, polygons); err != nil {
			log.Fatalf("rendering PNGs: %v", err)
		}
	}
}

func renderPNGs(outputFolder string, polygons [][]Point) error {
	imgs := make([]image.Image, 0, len(polygons))
	for _, polygon := range polygons {
		img, err := renderOnMap(polygon)
		if err != nil {
			return fmt.Errorf("rendering image on map: %w", err)
		}
		imgs = append(imgs, img)
	}

	if err := os.MkdirAll(outputFolder, 0o700); err != nil {
		return fmt.Errorf("creating directory for %v: %w", outputFolder, err)
	}

	for i, img := range imgs {
		fp := filepath.Join(outputFolder, fmt.Sprintf("map-%d.png", i))
		log.Printf("output to: %s", fp)
		if err := gg.SavePNG(fp, img); err != nil {
			return fmt.Errorf("saving png for %s: %w", fp, err)
		}
	}

	return nil
}

func googleMapsURLOf(polygon []Point) string {
	waypoints := make([]string, 0, len(polygon)-2)
	for _, p := range polygon {
		waypoints = append(waypoints, p.String())
	}

	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/dir/%s/data=!3m1!4b1!4m2!4m1!3e1", strings.Join(waypoints, "/"))
	return googleMapsURL
}

func renderOnMap(polygon []Point) (image.Image, error) {
	cont := sm.NewContext()
	cont.SetSize(1920, 1080)

	positions := make([]s2.LatLng, 0, len(polygon))
	for _, p := range polygon {
		positions = append(positions, s2.LatLngFromDegrees(p.Lat, p.Lng))
	}

	cont.AddObject(sm.NewPath(positions, color.Black, 1))

	for i, p := range positions {
		textImg := textImageOf(fmt.Sprint(i))
		cont.AddObject(sm.NewImageMarker(p, textImg, 0.5, 0.5))
	}

	img, err := cont.Render()
	if err != nil {
		return nil, fmt.Errorf("rendering image: %w", err)
	}

	return img, nil
}

func textImageOf(text string) image.Image {
	width, height := 20.0, 20.0

	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(color.Black)
	dc.Clear()
	dc.SetFontFace(basicfont.Face7x13)

	dc.SetColor(color.White)
	dc.DrawStringAnchored(text, height/2, width/2, 0.5, 0.5)

	return dc.Image()
}

func gpxDataOf(filename string) (*gpx.GPX, error) {
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

func polygonsOf(data *gpx.GPX) ([][]Point, error) {
	polygons := make([][]Point, 0, len(data.Tracks))
	for _, track := range data.Tracks {
		polygon := polygonOf(track)
		blurred := blur(polygon, 20)

		polygons = append(polygons, blurred)
	}

	return polygons, nil
}

func blur(polygon []Point, maxSize int) []Point {
	if len(polygon) <= maxSize {
		return polygon
	}

	return flatteningBlur(naiveBlur(polygon, 50))
}

func flatteningBlur(polygon []Point) []Point {
	out := make([]Point, 0)
	for i, p := range polygon {
		if i == 0 || i == len(polygon)-1 {
			out = append(out, p)
			continue
		}

		lastOut := out[len(out)-1]
		next := polygon[i+1]

		angle := math.Atan2(p.Lat-lastOut.Lat, p.Lng-lastOut.Lng) * (180 / math.Pi)
		angleNext := math.Atan2(next.Lat-p.Lat, next.Lng-p.Lng) * (180 / math.Pi)

		if math.Abs(angle-angleNext) > 30 {
			out = append(out, p)
		}
	}
	return out
}

// naiveBlur keeps the first, last, and every nth point in the list and returns it.
func naiveBlur(polygon []Point, nth int) []Point {
	smoothPolygon := make([]Point, 0, nth)
	atEvery := len(polygon) / nth
	for i, p := range polygon {
		if i == 0 || i == len(polygon)-1 || i%atEvery == 0 {
			smoothPolygon = append(smoothPolygon, p)
		}
	}
	return smoothPolygon
}

func polygonOf(track gpx.GPXTrack) []Point {
	polyon := make([]Point, 0)
	for _, segment := range track.Segments {
		for _, p := range segment.Points {
			polyon = append(polyon, Point{Lat: p.Latitude, Lng: p.Longitude})
		}
	}
	return polyon
}
