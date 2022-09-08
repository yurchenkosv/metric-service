package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/config"
	migration "github.com/yurchenkosv/metric-service/internal/migrate"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"github.com/yurchenkosv/metric-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yurchenkosv/metric-service/internal/routers"
)

var (
	cfg       = config.NewServerConfig()
	storeLoop *time.Ticker
	repo      repository.Repository
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	osSignal := make(chan os.Signal, 1)
	storeLoopStop := make(chan bool)
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}

	metricService := service.NewServerMetricService(cfg, repo)
	log.WithFields(
		log.Fields{
			"address": cfg.Address,
		}).Info("Starting metric server")

	if cfg.DBDsn != "" {
		migration.Migrate(cfg.DBDsn)
		repo = repository.NewPostgresRepo(cfg.DBDsn)
	} else {
		repo = repository.NewMapStorage()
	}

	if cfg.Restore {
		err := metricService.LoadMetricsFromDisk()
		if err != nil {
			log.Fatal("cannot read metrics from file")
		}
	}

	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	go func() {
		<-osSignal
		if cfg.StoreInterval != 0 && cfg.DBDsn == "" {
			storeLoopStop <- true
		}
		err := metricService.SaveMetricsToDisk()
		if err != nil {
			log.Error("cannot store metrics in file")
		}
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
					err := metricService.SaveMetricsToDisk()
					if err != nil {
						log.Error("cannot store metrics in file")
					}
				}
			}

		}()
	}

	router := routers.NewRouter(cfg, repo)
	server := &http.Server{Addr: cfg.Address, Handler: router}
	log.Fatal(server.ListenAndServe())
}
