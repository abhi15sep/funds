package main

import (
	"./morningStar"
	"log"
	"net/http"
)

const port = `1080`

func main() {
	http.HandleFunc(`/`, morningStar.Handler)

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
