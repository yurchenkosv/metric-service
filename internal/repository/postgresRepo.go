package repository

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/model"
)

type PostgresRepo struct {
	Conn  *sqlx.DB // Conn stores connection to current postgres instance
	DBURI string   // DBURI connection string to initialize postgres connection
}

// NewPostgresRepo initializes connection to postgres.
// It's also configures connection and returns pointer to current PostgresRepo.
func NewPostgresRepo(dbURI string) *PostgresRepo {
	conn, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(5)

	return &PostgresRepo{
		Conn:  conn,
		DBURI: dbURI,
	}
}

// Migrate executes migrations to set database to initial state
func (repo *PostgresRepo) Migrate(migrationsPath string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		repo.DBURI)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}

// Shutdown made to drop idle connections, wait current connections successfully done.
// Then connection to postgres instance will be terminated.
func (repo *PostgresRepo) Shutdown() {
	repo.Conn.SetMaxIdleConns(-1)
	repo.Conn.Close()
}

// SaveCounter saves in database counter metric.
// When conflict occurs, query updates current value of metric
func (repo *PostgresRepo) SaveCounter(name string, counter model.Counter) error {
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
	_, err := repo.Conn.Exec(query, name, "counter", int(counter))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// SaveGauge saves in database gauge metric.
// When conflict occurs, query updates current value of metric
func (repo *PostgresRepo) SaveGauge(name string, gauge model.Gauge) error {
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
	_, err := repo.Conn.Exec(query, name, "gauge", float64(gauge))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetMetricByKey selects all data in DB by key provided by name string parameter.
// Returns model.Metric as single value because of metric key is unic
func (repo *PostgresRepo) GetMetricByKey(name string) (*model.Metric, error) {
	var metric model.Metric
	query := `
		SELECT metric_id, metric_type, metric_delta, metric_value
		FROM metrics 
		WHERE metric_id = $1
		`
	err := repo.Conn.QueryRow(query, name).
		Scan(&metric.ID,
			&metric.MType,
			&metric.Delta,
			&metric.Value)
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// GetAllMetrics selects all metrics in DB and returns it as pointer to model.Metrics
func (repo *PostgresRepo) GetAllMetrics() (*model.Metrics, error) {
	var metrics model.Metrics

	query := "SELECT metric_id, metric_type, metric_delta, metric_value FROM metrics"

	result, err := repo.Conn.Query(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = result.Err()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var metricID, metricType string
		var metricDelta *model.Counter
		var metricValue *model.Gauge
		err = result.Scan(&metricID, &metricType, &metricDelta, &metricValue)
		if err != nil {
			continue
		}
		metrics.Metric = append(metrics.Metric, model.Metric{
			ID:    metricID,
			MType: metricType,
			Delta: metricDelta,
			Value: metricValue,
		})
	}
	return &metrics, nil
}

// Ping creates connection to DB and do select on it.
// When ping unsuccessful error returns. It means that database is unhealthy
func (repo *PostgresRepo) Ping() error {
	return repo.Conn.Ping()
}

// SaveMetricsBatch saves slice of model.Metric in DB in one transaction.
// When error occurs, transaction rollbacks.
func (repo *PostgresRepo) SaveMetricsBatch(metrics []model.Metric) error {
	tx, err := repo.Conn.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	for i := range metrics {
		query := `
		INSERT INTO metrics(
			metric_id,
			metric_type,
			metric_delta,
			metric_value,
			hash
		)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (metric_id) DO UPDATE
		SET metric_delta=metrics.metric_delta+$3,
			metric_value=$4,
			hash=$5;
		`
		_, err = tx.Exec(
			query,
			metrics[i].ID,
			metrics[i].MType,
			metrics[i].Delta,
			metrics[i].Value,
			metrics[i].Hash,
		)
		if err != nil {
			log.Error(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return tx.Rollback()
	}
	return nil
}
