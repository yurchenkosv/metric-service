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
)

var (
	cfg = config.AgentConfig{}
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	poolCount := 1
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

	sched := gocron.NewScheduler(time.UTC)
	_, err = sched.Every(cfg.PollInterval).
		Do(agentService.CollectMetrics, &poolCount)
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

	<-osSignal
	sched.Stop()
	fmt.Println("Program exit")
	os.Exit(0)

}
