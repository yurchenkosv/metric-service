package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/functions"
	migration "github.com/yurchenkosv/metric-service/internal/migrate"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yurchenkosv/metric-service/internal/routers"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg        = types.ServerConfig{}
	storeLoop  *time.Ticker
	mapStorage storage.Repository
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	osSignal := make(chan os.Signal, 1)
	storeLoopStop := make(chan bool)
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(
		log.Fields{
			"address": cfg.Address,
		}).Info("Starting metric agent")

	if cfg.DBDsn != "" {
		migration.Migrate(cfg.DBDsn)
		mapStorage = storage.NewPostgresStorage(&cfg)
	} else {
		mapStorage = storage.NewMapStorage()
	}

	if cfg.Restore {
		mapStorage = functions.ReadMetricsFromDisk(&cfg, &mapStorage)
	}

	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	go func() {
		<-osSignal
		if cfg.StoreInterval != 0 && cfg.DBDsn == "" {
			storeLoopStop <- true
		}
		functions.FlushMetricsToDisk(&cfg, mapStorage)
		os.Exit(0)
	}()

	if cfg.StoreInterval != 0 && cfg.DBDsn == "" {
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
