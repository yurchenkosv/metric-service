package service

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/repository"
)

// ServerMetricService main service to work with metrics on server side
type ServerMetricService struct {
	config            *config.ServerConfig
	repo              repository.Repository
	saveMetricsToDisk bool
}

// NewServerMetricService creates new ServerMetricService and returns pointer on it.
func NewServerMetricService(cnf *config.ServerConfig, repo repository.Repository) *ServerMetricService {
	needSave := false
	if cnf.StoreFile != "" {
		needSave = true
	}
	return &ServerMetricService{
		config:            cnf,
		repo:              repo,
		saveMetricsToDisk: needSave,
	}
}

// Shutdown method saves metrics to disk, if it's required, then call repository Shutdown() method.
// This method for gracefully shutdown operations before shutting down application
func (s ServerMetricService) Shutdown() {
	if s.config.StoreInterval != 0 && s.config.DBDsn == "" {
		ctx := context.Background()
		err := s.SaveMetricsToDisk(ctx)
		if err != nil {
			log.Error("cannot store metrics in file")
		}
	}
	s.repo.Shutdown()
}

// AddMetric method for saving model.Metric in repository. It can return error if something went wrong.
// Also, it saves metrics to disk if it is required.
func (s *ServerMetricService) AddMetric(ctx context.Context, metric model.Metric) error {
	switch metric.MType {
	case "gauge":
		err := s.repo.SaveGauge(ctx, metric.ID, *metric.Value)
		if err != nil {
			log.Error(err)
			return err
		}
	case "counter":
		err := s.repo.SaveCounter(ctx, metric.ID, *metric.Delta)
		if err != nil {
			log.Error(err)
			return err
		}
	default:
		return &errors.NoSuchMetricError{MetricName: metric.ID}
	}
	err := s.SaveMetricsToDisk(ctx)
	if err != nil {
		log.Error(err)
	}
	return nil
}

// AddMetricBatch method for save slice of model.Metrics to repository.
// It can return error if it cannot save metrics for some reason.
// It calls SaveMetricsBatch repository method to save metrics transactional.
// Also, it saves metrics to disk if it is required.
func (s *ServerMetricService) AddMetricBatch(ctx context.Context, metrics model.Metrics) error {
	var err error

	err = s.repo.SaveMetricsBatch(ctx, metrics.Metric)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.SaveMetricsToDisk(ctx)
	if err != nil {
		log.Error(err)
	}
	return nil
}

// GetMetricByKey method to get pointer to model.Metric from repository.
// If nothing could be found, then MetricNotFoundError returns.
func (s *ServerMetricService) GetMetricByKey(ctx context.Context, name string) (*model.Metric, error) {
	metric, err := s.repo.GetMetricByKey(ctx, name)
	if err != nil {
		return nil, &errors.MetricNotFoundError{MetricName: name}
	}
	return metric, nil
}

// GetAllMetrics method to get pointer to model.Metrics with all metrics in repository.
func (s *ServerMetricService) GetAllMetrics(ctx context.Context) (*model.Metrics, error) {
	metrics, err := s.repo.GetAllMetrics(ctx)
	if err != nil {
		log.Error("error getting metrics", err)
		return nil, err
	}
	return metrics, nil
}

// CreateSignedHash creates hash signed by key in config. Then this hash adds to metric.
func (s *ServerMetricService) CreateSignedHash(msg string) (string, error) {
	if s.config.HashKey == "" {
		return "", &errors.NoEncryptionKeyFoundError{}
	}
	hash := signHash(s.config.HashKey, msg)
	return hash, nil
}

// SaveMetricsToDisk method to save all metrics to file if Config.File and StoreInterval is set.
// It gets all metrics from repository, then marshall it to JSON and write to file.
func (s *ServerMetricService) SaveMetricsToDisk(ctx context.Context) error {
	var mutex sync.Mutex
	if !s.saveMetricsToDisk {
		return nil
	}

	fileLocation := s.config.StoreFile
	fileBits := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	mutex.Lock()
	file, err := os.OpenFile(fileLocation, fileBits, 0600)
	if err != nil {
		log.Fatal(err)
	}

	metrics, err := s.GetAllMetrics(ctx)
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

// LoadMetricsFromDisk method to restore metrics state from disk if Config.Restore is true.
// It reads Config.File, unmarshalls and then load model.Metrics to repository
func (s *ServerMetricService) LoadMetricsFromDisk(ctx context.Context) error {
	fileLocation := s.config.StoreFile
	data, err := os.ReadFile(fileLocation)
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

	err = s.AddMetricBatch(ctx, metrics)
	if err != nil {
		return err
	}
	return nil
}
