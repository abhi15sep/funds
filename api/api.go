package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/funds/db"
	"github.com/ViBiOh/funds/model"
)

const port = `1080`

var modelHandler = model.Handler{}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if len(model.ListFunds()) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func fundsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == `/health` {
		healthHandler(w, r)
	} else {
		modelHandler.ServeHTTP(w, r)
	}
}

func handleGracefulClose(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	<-signals
	log.Print(`SIGTERM received`)

	if server != nil {
		log.Print(`Shutting down http server`)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Print(err)
		}
	}
}

func main() {
	url := flag.String(`c`, ``, `URL to healthcheck (check and exit)`)
	infosURL := flag.String(`infos`, ``, `Informations URL`)
	flag.Parse()

	if *url != `` {
		alcotest.Do(url)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := db.Init(); err != nil {
		log.Printf(`Error while initializing database: %v`, err)
	} else if db.Ping() {
		log.Print(`Database ready`)
	}

	if err := model.Init(*infosURL); err != nil {
		log.Printf(`Error while initializing model: %v`, err)
	}

	log.Print(`Starting server on port ` + port)

	server := &http.Server{
		Addr:    `:` + port,
		Handler: http.HandlerFunc(fundsHandler),
	}

	go server.ListenAndServe()
	handleGracefulClose(server)
}
