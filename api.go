package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/funds/model"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
)

const port = `1080`

var modelHandler = model.Handler()
var apiHandler http.Handler

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if len(model.ListFunds()) > 0 && model.Health() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == `/health` {
			healthHandler(w, r)
		} else {
			modelHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	alcotestConfig := alcotest.Flags(``)
	prometheusConfig := prometheus.Flags(`prometheus`)
	rateConfig := rate.Flags(`rate`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	if err := model.Init(); err != nil {
		log.Printf(`Error while initializing model: %v`, err)
	}

	log.Print(`Starting server on port ` + port)

	apiHandler = prometheus.Handler(prometheusConfig, rate.Handler(rateConfig, gziphandler.GzipHandler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler())))))
	server := &http.Server{
		Addr:    `:` + port,
		Handler: apiHandler,
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		serveError <- server.ListenAndServe()
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}