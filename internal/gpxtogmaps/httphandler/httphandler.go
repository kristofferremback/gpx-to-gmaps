package httphandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kristofferostlund/recommendli/pkg/srv"

	"github.com/kristofferostlund/gpx-to-gmaps/internal/gpxtogmaps"
)

type handler struct {
	s *gpxtogmaps.Service
}

func New(s *gpxtogmaps.Service) http.Handler {
	h := &handler{s}

	r := chi.NewRouter()

	r.Post("/gpx-to-gmaps", h.postGPXToGMaps())
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.JSONError(w, fmt.Errorf("Not found"), srv.Status(404))
	}))

	return r
}

func (h *handler) postGPXToGMaps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.String())
		if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
			srv.JSONError(w, err, srv.Status(400))
		}

		file, _, err := r.FormFile("gpx")
		if err != nil {
			srv.JSONError(w, err, srv.Status(400))
			return
		}
		defer file.Close()

		polygons, err := h.s.ConvertToPolygons(file, 25)
		if err != nil {
			srv.JSONError(w, err, srv.Status(500))
			return
		}

		gmapsURLs := make([]string, 0, len(polygons))
		for _, polygon := range polygons {
			gmapsURLs = append(gmapsURLs, h.s.GoogleMapsURL(polygon))
		}

		srv.JSON(w, struct {
			GmapsURLs []string `json:"google_maps_urls"`
		}{GmapsURLs: gmapsURLs})
	}
}
