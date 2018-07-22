package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ViBiOh/funds/pkg/model"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/db"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/rollbar"
	"github.com/ViBiOh/httputils/pkg/server"
)

func main() {
	serverConfig := httputils.Flags(``)
	alcotestConfig := alcotest.Flags(``)
	opentracingConfig := opentracing.Flags(`tracing`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	rollbarConfig := rollbar.Flags(`rollbar`)

	fundsConfig := model.Flags(``)
	dbConfig := db.Flags(`db`)

	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	serverApp := httputils.NewApp(serverConfig)
	healthcheckApp := healthcheck.NewApp()
	opentracingApp := opentracing.NewApp(opentracingConfig)
	owaspApp := owasp.NewApp(owaspConfig)
	corsApp := cors.NewApp(corsConfig)
	rollbarApp := rollbar.NewApp(rollbarConfig)
	gzipApp := gzip.NewApp()

	fundApp, err := model.NewApp(fundsConfig, dbConfig)
	if err != nil {
		log.Fatalf(`Error while creating Fund app: %v`, err)
	}

	modelHandler := model.Handler(fundApp)
	healthcheckApp.NextHealthcheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(fundApp.ListFunds()) > 0 && fundApp.Health() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))

	handler := server.ChainMiddlewares(modelHandler, opentracingApp, rollbarApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(handler, nil, healthcheckApp, rollbarApp)
}
