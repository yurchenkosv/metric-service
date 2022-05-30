package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/yurchenkosv/metric-service/internal/types"
	"log"
	"os"
)

type PostgresStorage struct {
	Conn string
}

func NewPostgresStorage(cfg *types.ServerConfig) Repository {
	conn, err := pgx.Connect(context.Background(), cfg.DBDsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	return &PostgresStorage{Conn: cfg.DBDsn}
}

func (p *PostgresStorage) AddCounter(name string, counter types.Counter) {
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	query := `
		INSERT INTO metrics(
		metric_id,
		metric_type,
		metric_delta 
		)
		VALUES($1, $2, $3)
		ON CONFLICT (metric_id) DO UPDATE
		SET metric_delta=metrics.metric_delta+$3;
	`
	_, err = conn.Exec(context.Background(), query, name, "counter", int(counter))
	if err != nil {
		log.Println(err)
	}
}

func (p *PostgresStorage) AddGauge(name string, gauge types.Gauge) {
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	query := `
		INSERT INTO metrics(
			metric_id,
			metric_type,
			metric_value 
		)
		VALUES($1, $2, $3)
		ON CONFLICT (metric_id) DO UPDATE
		SET metric_value=$3;
	`
	_, err = conn.Exec(context.Background(), query, name, "gauge", float64(gauge))
	if err != nil {
		log.Println(err)
	}
}

func (p *PostgresStorage) GetMetricByKey(name string) (string, error) {
	var counter string
	var gauge string
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	query := "SELECT metric_delta, metric_value FROM metrics WHERE metric_id = $1"
	result, err := conn.Query(context.Background(), query, name)
	if err != nil {
		return "", err
	}
	defer result.Close()

	for result.Next() {
		result.Scan(&counter, &gauge)
		if counter != "" {
			return counter, nil
		} else {
			return gauge, nil
		}
	}

	return "", ErrNotFound
}

func (p *PostgresStorage) GetCounterByKey(name string) (types.Counter, error) {
	var counter types.Counter
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	query := "SELECT metric_delta FROM metrics WHERE metric_id=$1"

	result, err := conn.Query(context.Background(), query, name)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	for result.Next() {
		result.Scan(&counter)
	}
	return counter, nil
}

func (p *PostgresStorage) GetGaugeByKey(name string) (types.Gauge, error) {
	var gauge types.Gauge
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	query := "SELECT metric_value FROM metrics WHERE metric_id=$1"

	result, err := conn.Query(context.Background(), query, name)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	for result.Next() {
		result.Scan(&gauge)
	}

	return gauge, nil
}

func (p *PostgresStorage) GetAllMetrics() string {
	var metrics string
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	query := "SELECT metric_id, metric_delta FROM metrics WHERE metric_type='counter'"

	result, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
	}
	defer result.Close()

	for result.Next() {
		var value, key string
		result.Scan(&key, &value)
		metrics = metrics + fmt.Sprintf("%s = %s \n", key, value)
	}

	query = "SELECT metric_id, metric_value FROM metrics WHERE metric_type='gauge'"
	result, err = conn.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
	}
	defer result.Close()

	for result.Next() {
		var value, key string
		result.Scan(&key, &value)
		metrics = metrics + fmt.Sprintf("%s = %s \n", key, value)
	}
	return metrics
}

func (p *PostgresStorage) AsMetrics() types.Metrics {
	var metrics types.Metrics
	conn, err := pgx.Connect(context.Background(), p.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	query := "SELECT metric_id, metric_type, metric_delta, metric_value, hash FROM metrics"

	result, err := conn.Query(context.Background(), query)
	defer result.Close()

	if err != nil {
		log.Println(err)
	}

	for result.Next() {
		var metricId, metricType, hash string
		var metricDelta *int64
		var metricValue *float64
		result.Scan(&metricId, &metricType, &metricDelta, &metricValue, &hash)
		metrics.Metric = append(metrics.Metric, types.Metric{
			ID:    metricId,
			MType: metricType,
			Delta: metricDelta,
			Value: metricValue,
			Hash:  hash,
		})
	}
	return metrics
}
