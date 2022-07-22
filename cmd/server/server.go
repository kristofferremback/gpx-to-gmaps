package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
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
	addr     = flag.String("addr", defaultAddr, "HTTP address")
	logLevel = flag.String("log-level", logging.LevelInfo.String(), "log level")
)

func main() {
	flag.Parse()

	log := logging.New(logging.GetLevelByName(*logLevel), logging.FormatConsolePretty)
	defer log.Sync() //nolint: errcheck

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(1 * time.Minute))

	r.Get("/status", getStatus())

	r.Mount("/api", httphandler.New(gpxtogmaps.NewService()))

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
