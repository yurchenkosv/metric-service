package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yurchenkosv/metric-service/internal/routers"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg        = types.Config{}
	mapStorage = storage.NewMapStorage()
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	osSignal := make(chan os.Signal, 3)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-osSignal
		functions.FlushMetricsToDisk(&cfg, mapStorage)
		os.Exit(0)
	}()
	router := routers.NewRouter(&cfg, &mapStorage)
	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
