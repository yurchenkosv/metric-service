package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/clients"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/service"
	"github.com/yurchenkosv/metric-service/pkg/finalizer"
)

var (
	cfg          = config.AgentConfig{}
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {

	fmt.Printf(" Build version: %s\n Build date: %s\n Build commit: %s\n", buildVersion, buildDate, buildCommit)

	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(
		log.Fields{
			"poolInterval": cfg.PollInterval,
			"address":      cfg.Address,
		}).Info("Starting metric agent")

	metricServerClient := clients.NewMetricServerClient(cfg.Address)
	agentService := service.NewAgentMetricService(&cfg, metricServerClient)
	if cfg.CryptoKey != "" {
		encryptionService, err2 := service.NewEncryptionService(cfg.CryptoKey)
		if err2 != nil {
			log.Fatal("cannot load public key specified: ", err2)
		}
		agentService = agentService.
			WithRSAMessagesEncryption(encryptionService)
	}

	sched := gocron.NewScheduler(time.UTC)
	_, err = sched.Every(cfg.PollInterval).
		Do(agentService.CollectMetrics, 1)
	if err != nil {
		log.Fatal("cannot start collect job", err)
	}

	_, err = sched.Every(cfg.ReportInterval).
		Do(agentService.Push)
	if err != nil {
		log.Fatal("cannot start report job", err)
	}
	sched.StartAsync()
	osSignal := make(chan os.Signal, 3)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	finalizer.Shutdown(func() {
		<-osSignal
		sched.Stop()
		fmt.Println("Program exit")
	})

}
