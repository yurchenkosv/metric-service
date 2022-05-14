package main

import (
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"

	"github.com/yurchenkosv/metric-service/internal/routers"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg = types.Config{}
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := routers.NewRouter()
	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
