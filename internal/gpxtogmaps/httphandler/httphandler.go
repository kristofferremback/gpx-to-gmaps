package httphandler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kristofferostlund/recommendli/pkg/srv"

	"github.com/kristofferostlund/gpx-to-gmaps/internal/gpxtogmaps"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/geo"
	"github.com/kristofferostlund/gpx-to-gmaps/pkg/gmapsurl"
)

type handler struct {
	s *gpxtogmaps.Service

	baseURL string
}

func New(s *gpxtogmaps.Service, baseURL string) http.Handler {
	h := &handler{s, baseURL}

	r := chi.NewRouter()

	r.Post("/convert-gpx", h.postGPXToGMaps())
	r.Get("/static-map", h.getGenerateStaticMap())
	r.Handle("/*", h.notFoundHandler())

	return r
}

func (handler) notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.JSONError(w, fmt.Errorf("Not found"), srv.Status(http.StatusNotFound))
	}
}

func (h *handler) postGPXToGMaps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicleType, ok := map[string]gmapsurl.VehicleType{
			"car":     gmapsurl.Car,
			"bike":    gmapsurl.Bike,
			"walking": gmapsurl.Walking,
		}[r.FormValue("vehicle_type")]
		if !ok {
			srv.JSONError(w, fmt.Errorf("invalid vehicle type %v", r.FormValue("vehicle_type")), srv.Status(http.StatusBadRequest))
			return
		}
		maxPrecision, err := strconv.Atoi(r.FormValue("max_precision"))
		if err != nil {
			srv.JSONError(w, fmt.Errorf("parsing max_precision: %w", err), srv.Status(http.StatusBadRequest))
			return
		}

		if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			srv.JSONError(w, err, srv.Status(http.StatusBadRequest))
			return
		}

		file, _, err := r.FormFile("gpx")
		if err != nil {
			srv.JSONError(w, err, srv.Status(http.StatusBadRequest))
			return
		}
		defer file.Close()

		polygons, err := h.s.ConvertToPolygons(file, maxPrecision)
		if err != nil {
			srv.JSONError(w, err, srv.Status(http.StatusInternalServerError))
			return
		}

		gmapsURLs := make([]string, 0, len(polygons))
		mapURLs := make([]string, 0, len(polygons))

		for _, polygon := range polygons {
			gmapsURLs = append(gmapsURLs, h.s.GoogleMapsURL(polygon, vehicleType))
			mapURLs = append(mapURLs, h.staticMapURLFor(r, polygon))
		}

		srv.JSON(w, struct {
			GmapsURLs []string `json:"google_maps_urls"`
			MapsURLs  []string `json:"maps_urls"`
		}{GmapsURLs: gmapsURLs, MapsURLs: mapURLs})
	}
}

func (h *handler) getGenerateStaticMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryPolyline := r.URL.Query().Get("polyline")
		if queryPolyline == "" {
			srv.JSONError(w, errors.New("polyline must be provided"), srv.Status(http.StatusBadRequest))
			return
		}

		pl, err := url.QueryUnescape(queryPolyline)
		if err != nil {
			srv.JSONError(w, fmt.Errorf("invalid escaping of polyline: %w", err), srv.Status(http.StatusBadRequest))
			return
		}
		polygon, err := h.s.DecodePolyline(pl)
		if err != nil {
			srv.JSONError(w, fmt.Errorf("decoding polyline: %w", err), srv.Status(http.StatusInternalServerError))
			return
		}

		w.Header().Set("content-type", "image/png")
		if err := h.s.PNG(w, polygon); err != nil {
			srv.JSONError(w, fmt.Errorf("generating static map: %w", err), srv.Status(http.StatusInternalServerError))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *handler) staticMapURLFor(r *http.Request, polygon geo.Polygon) string {
	return fmt.Sprintf("%s/static-map?polyline=%s", h.baseURL, url.QueryEscape(h.s.EncodePolyline(polygon)))
}
