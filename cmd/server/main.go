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
	"time"

	"github.com/yurchenkosv/metric-service/internal/routers"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg        = types.Config{}
	mapStorage = storage.NewMapStorage()
	storeLoop  *time.Ticker
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Restore {
		mapStorage = functions.ReadMetricsFromDisk(&cfg, &mapStorage)
	}
}

func main() {
	osSignal := make(chan os.Signal, 1)
	storeLoopStop := make(chan bool)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	go func() {
		<-osSignal
		if cfg.StoreInterval != 0 {
			storeLoopStop <- true
		}
		functions.FlushMetricsToDisk(&cfg, mapStorage)
		os.Exit(0)
	}()

	if cfg.StoreInterval != 0 {
		storeLoop = time.NewTicker(cfg.StoreInterval)
		go func() {
			for {
				select {
				case <-storeLoopStop:
					return
				case <-storeLoop.C:
					functions.FlushMetricsToDisk(&cfg, mapStorage)
				}
			}

		}()
	}

	router := routers.NewRouter(&cfg, &mapStorage)
	server := &http.Server{Addr: cfg.Address, Handler: router}
	log.Fatal(server.ListenAndServe())
}
