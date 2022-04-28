package main

import (
	"log"
	"net/http"

	"github.com/yurchenkosv/metric-service/internal/routers"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	server = types.URLServer{}
)

func main() {
	router := routers.NewRouter()
	serveFor := server.
		SetHost("localhost").
		SetPort("8080").
		Build()
	log.Fatal(http.ListenAndServe(serveFor, router))
}
