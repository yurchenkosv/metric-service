package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"github.com/yurchenkosv/metric-service/internal/service"

	"github.com/yurchenkosv/metric-service/internal/routers"
)

var (
	cfg  = config.NewServerConfig()
	repo repository.Repository
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(
		log.Fields{
			"address": cfg.Address,
		}).Info("Starting metric server")

	if cfg.DBDsn != "" {
		repo = repository.NewPostgresRepo(cfg.DBDsn)
		repo.Migrate("db/migrations")
	} else {
		repo = repository.NewMapRepo()
	}

	metricService := service.NewServerMetricService(cfg, repo)
	if cfg.Restore {
		err := metricService.LoadMetricsFromDisk()
		if err != nil {
			log.Fatal("cannot read metrics from file")
		}
	}

	sched := gocron.NewScheduler(time.UTC)
	if cfg.StoreInterval != 0 && cfg.DBDsn == "" {
		_, err := sched.Every(cfg.StoreInterval).
			Do(metricService.SaveMetricsToDisk)
		if err != nil {
			log.Error("cannot save metrics to disk", err)
		}
		sched.StartAsync()
	}

	router := routers.NewRouter(cfg, repo)
	server := &http.Server{Addr: cfg.Address, Handler: router}
	go func(server *http.Server) {
		log.Warn(server.ListenAndServe())
	}(server)

	<-osSignal
	log.Warn("shuting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}

	sched.Stop()
	metricService.Shutdown()
	os.Exit(0)
}
