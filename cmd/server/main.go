package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"github.com/yurchenkosv/metric-service/internal/service"
	"github.com/yurchenkosv/metric-service/pkg/finalizer"

	"github.com/yurchenkosv/metric-service/internal/api"
	"github.com/yurchenkosv/metric-service/internal/handlers"
	"github.com/yurchenkosv/metric-service/internal/routers"
	"google.golang.org/grpc"
	"net"
)

var (
	cfg               = config.NewServerConfig()
	repo              repository.Repository
	buildVersion      = "N/A"
	buildDate         = "N/A"
	buildCommit       = "N/A"
	mainContext       = context.Background()
	grpcServerOptions []grpc.ServerOption
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {

	fmt.Printf(" Build version: %s\n Build date: %s\n Build commit: %s\n", buildVersion, buildDate, buildCommit)

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
		err2 := metricService.LoadMetricsFromDisk(mainContext)
		if err2 != nil {
			log.Fatal("cannot read metrics from file")
		}
	}

	sched := gocron.NewScheduler(time.UTC)
	if cfg.StoreInterval != 0 && cfg.DBDsn == "" {
		_, err2 := sched.Every(cfg.StoreInterval).
			Do(metricService.SaveMetricsToDisk, mainContext)
		if err2 != nil {
			log.Error("cannot save metrics to disk", err2)
		}
		sched.StartAsync()
	}

	router := routers.NewRouter(cfg, repo)
	server := &http.Server{Addr: cfg.Address, Handler: router}

	if cfg.CryptoKey != "" {
		tlsService, err := service.NewServerTLSService(*cfg)
		if err != nil {
			log.Fatal(err)
		}
		_, err = tlsService.CreatePemCertificateFromPrivateKey(strings.Split(cfg.Address, ":")[0])
		if err != nil {
			log.Fatal("cannot create expected certificate ", err)
		}
		cert, err := tlsService.SaveCertificateToDisk()
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig, err := tlsService.GetCredentialConfig()
		if err != nil {
			log.Fatal("cannot get credential config for grpc server ", err)
		}

		grpcServerOptions = append(grpcServerOptions, grpc.Creds(tlsConfig))

		go func(server *http.Server) {
			log.Warn(server.ListenAndServeTLS(cert, cfg.CryptoKey))
		}(server)
	} else {
		go func(server *http.Server) {
			log.Warn(server.ListenAndServe())
		}(server)
	}

	grpcMetricsHanlrer := handlers.NewGRPCMetricHandler(metricService)
	grpcHealthCheckHandler := handlers.NewGRPCHealthCheckHandler(service.NewHealthCheckService(cfg, repo))
	grpcServer := grpc.NewServer(grpcServerOptions...)
	api.RegisterMetricServiceServer(grpcServer, grpcMetricsHanlrer)
	api.RegisterHealthcheckServer(grpcServer, grpcHealthCheckHandler)
	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	go func(listener net.Listener) {
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Error(err)
		}
	}(listener)

	<-osSignal
	log.Warn("shutting down server")

	ctx, cancel := context.WithTimeout(mainContext, 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
	grpcServer.GracefulStop()
	finalizer.Shutdown(func() {
		sched.Stop()
		metricService.Shutdown()
	})
}
