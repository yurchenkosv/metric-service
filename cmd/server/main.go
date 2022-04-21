package main

import (
	"github.com/yurchenkosv/metric-service/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/update/", handlers.HandleMetric)
	server := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
