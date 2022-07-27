package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kristofferostlund/gpx-to-gmaps/internal/gpxtogmaps"
	"github.com/kristofferostlund/gpx-to-gmaps/internal/gpxtogmaps/httphandler"
	"github.com/kristofferostlund/recommendli/pkg/logging"
	"github.com/kristofferostlund/recommendli/pkg/srv"
)

const defaultAddr = ":9876"

var (
	addr       = flag.String("addr", defaultAddr, "HTTP address")
	logLevel   = flag.String("log-level", logging.LevelInfo.String(), "log level")
	appBaseURL = flag.String("app-base-url", fmt.Sprintf("http://localhost%s", defaultAddr), "base URL the app is exposed on, incl protocol")
)

func main() {
	flag.Parse()

	log := logging.New(logging.GetLevelByName(*logLevel), logging.FormatConsolePretty)
	defer log.Sync() //nolint: errcheck

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				middleware.Logger(h).ServeHTTP(w, r)
				return
			}
			h.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(1 * time.Minute))

	r.Get("/status", getStatus())
	r.Mount("/api", httphandler.New(log, gpxtogmaps.NewService(), fmt.Sprintf("%s/api", *appBaseURL)))

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/*", srv.RedirectOn404(fs, "/index.html"))

	errs := make(chan error, 2)
	go func() {
		log.Info("Starting server", "addr", *addr)
		errs <- http.ListenAndServe(*addr, r)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c
		errs <- nil
	}()

	if err := <-errs; err != nil && err != http.ErrServerClosed {
		log.Fatal("Shutting down", err)
	}

	log.Info("Server shutdown")
}

func getStatus() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK")) //nolint
	})
}
