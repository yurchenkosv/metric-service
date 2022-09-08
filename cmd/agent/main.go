package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/clients"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfg = config.AgentConfig{}
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(
		log.Fields{
			"poolInterval": cfg.PollInterval,
			"address":      cfg.Address,
		}).Info("Starting metric agent")

	agentService := service.NewAgentMetricService(&cfg)
	metricServerClient := clients.NewMetricServerClient(cfg.Address)

	mainLoop := time.NewTicker(cfg.PollInterval)
	pushLoop := time.NewTicker(cfg.ReportInterval)
	mainLoopStop := make(chan bool)
	memMetrics := make(chan model.Metrics)
	osSignal := make(chan os.Signal, 3)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		var pollCount int
		for {
			select {
			case <-mainLoopStop:
				return
			case <-mainLoop.C:
				pollCount = 1
				agentService.CollectMetrics(pollCount)
			case <-pushLoop.C:
				memMetrics <- agentService.CollectMetrics(pollCount)
			}
		}
	}()

	go func() {
		for {
			metricServerClient.PushMetrics(<-memMetrics)
		}
	}()

	<-osSignal
	Cleanup(mainLoop, pushLoop, mainLoopStop)

}
