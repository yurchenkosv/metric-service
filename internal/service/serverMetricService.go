package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"io/ioutil"
	"os"
	"sync"
)

type ServerMetricService struct {
	config *config.ServerConfig
	repo   repository.Repository
}

func NewServerMetricService(cnf *config.ServerConfig, repo repository.Repository) *ServerMetricService {
	return &ServerMetricService{
		config: cnf,
		repo:   repo,
	}
}

func (s ServerMetricService) Shutdown() {
	if s.config.StoreInterval != 0 && s.config.DBDsn == "" {
		err := s.SaveMetricsToDisk()
		if err != nil {
			log.Error("cannot store metrics in file")
		}
	}
	s.repo.Shutdown()
}

func (s *ServerMetricService) AddMetric(metric model.Metric) error {
	switch metric.MType {
	case "gauge":
		err := s.repo.SaveGauge(metric.ID, *metric.Value)
		if err != nil {
			log.Error(err)
			return err
		}
	case "counter":
		err := s.repo.SaveCounter(metric.ID, *metric.Delta)
		if err != nil {
			log.Error(err)
			return err
		}
	default:
		return &errors.NoSuchMetricError{MetricName: metric.ID}
	}
	err := s.SaveMetricsToDisk()
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (s *ServerMetricService) AddMetricBatch(metrics model.Metrics) error {
	var err error

	err = s.repo.SaveMetricsBatch(metrics.Metric)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.SaveMetricsToDisk()
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (s *ServerMetricService) GetMetricByKey(name string) (*model.Metric, error) {
	metric, err := s.repo.GetMetricByKey(name)
	if err != nil {
		return nil, &errors.MetricNotFoundError{MetricName: name}
	}
	return metric, nil
}

func (s *ServerMetricService) GetAllMetrics() (*model.Metrics, error) {
	metrics, err := s.repo.GetAllMetrics()
	if err != nil {
		log.Error("error getting metrics", err)
		return nil, err
	}
	return metrics, nil
}

func (s *ServerMetricService) CreateSignedHash(msg string) (string, error) {
	if s.config.HashKey == "" {
		return "", &errors.NoEncryptionKeyFoundError{}
	}
	hash := signHash(s.config.HashKey, msg)
	return hash, nil
}

func (s *ServerMetricService) SaveMetricsToDisk() error {
	var mutex sync.Mutex
	if s.config.StoreFile == "" {
		return nil
	}

	fileLocation := s.config.StoreFile
	fileBits := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	mutex.Lock()
	file, err := os.OpenFile(fileLocation, fileBits, 0600)
	if err != nil {
		log.Fatal(err)
	}

	metrics, err := s.GetAllMetrics()
	if err != nil {
		log.Error(err)
		return err
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		log.Error(err)
		return err
	}

	file.Close()
	mutex.Unlock()
	return nil
}

func (s *ServerMetricService) LoadMetricsFromDisk() error {
	fileLocation := s.config.StoreFile

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Println(err)
		os.Create(fileLocation)
		return nil
	}
	metrics := model.Metrics{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = s.AddMetricBatch(metrics)
	if err != nil {
		return err
	}
	return nil
}
