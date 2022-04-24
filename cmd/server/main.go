package main

import (
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/routers"
	"log"
	"net/http"
)

var (
	server = "localhost"
	port   = "8080"
)

func main() {
	router := routers.NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", server, port), router))
}
