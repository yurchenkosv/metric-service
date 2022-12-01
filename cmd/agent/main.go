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
			"pollInterval": cfg.PollInterval,
			"sendInterval": cfg.ReportInterval,
			"address":      cfg.Address,
		}).Info("Starting metric agent")

	ip, err := resolveIP(cfg.Address)
	if err != nil {
		log.Fatal("cannot resolve ip by bind hostname ", err)
	}
	metricServerClient := clients.NewMetricServerClient(cfg.Address).SetHeader("X-Real-IP", ip.String())
	agentService := service.NewAgentMetricService(&cfg, metricServerClient)
	if cfg.CryptoKey != "" {
		agentTLSService, err2 := service.NewAgentTLSService(cfg)
		if err2 != nil {
			log.Fatal("cannot load public key specified: ", err2)
		}
		metricServerClient.WithTLS(agentTLSService.GetTLSConfig()).SetScheme("https")
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
