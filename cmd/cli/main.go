package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tkrajina/gpxgo/gpx"
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

	googleMapsURL := fmt.Sprintf("https://www.google.com/maps/dir/%s", strings.Join(waypoints, "/"))
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
	img, err := cont.Render()
	if err != nil {
		return nil, fmt.Errorf("rendering image: %w", err)
	}

	return img, nil
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
		smoothPolygon := blur(polygon, 20)

		polygons = append(polygons, smoothPolygon)
	}

	return polygons, nil
}

func blur(polygon []Point, maxSize int) []Point {
	if len(polygon) <= maxSize {
		return polygon
	}

	// TODO: Try out a more complex algorithm for this, perhaps one being clever with the shape?
	// Like basically straight lines could be made into a single point
	// and non-straight lines could get slightly higher resolution.
	smoothPolygon := make([]Point, 0, maxSize)
	atEvery := len(polygon) / maxSize
	for i, p := range polygon {
		if i%atEvery == 0 || i == len(polygon)-1 {
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

func polylineOf(polygon []Point) [][]float64 {
	pl := make([][]float64, 0, len(polygon))
	for _, p := range polygon {
		pl = append(pl, []float64{p.Lat, p.Lng})
	}
	return pl
}

func fmtTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z07:00")
}
