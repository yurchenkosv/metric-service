package repository

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/model"
)

type PostgresRepo struct {
	Conn  *sqlx.DB
	DBURI string
}

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

func (repo *PostgresRepo) Shutdown() {
	repo.Conn.SetMaxIdleConns(-1)
	repo.Conn.Close()
}

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

func (repo *PostgresRepo) GetAllMetrics() (*model.Metrics, error) {
	var metrics model.Metrics

	query := "SELECT metric_id, metric_type, metric_delta, metric_value, hash FROM metrics"

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

func (repo *PostgresRepo) Ping() error {
	return repo.Conn.Ping()
}

//func (repo *PostgresRepo) InsertMetrics(metrics []types.Metric) {
//	conn, err := pgx.Connect(context.Background(), repo.Conn)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//	}
//	defer conn.Close(context.Background())
//	tx, err := conn.Begin(context.Background())
//	if err != nil {
//		log.Println(err)
//	}
//	for i := range metrics {
//		query := `
//		INSERT INTO metrics(
//			metric_id,
//			metric_type,
//			metric_delta,
//			metric_value,
//			hash
//		)
//		VALUES($1, $2, $3, $4, $5)
//		ON CONFLICT (metric_id) DO UPDATE
//		SET metric_delta=metrics.metric_delta+$3,
//			metric_value=$4,
//			hash=$5;
//		`
//		_, err = tx.Exec(context.Background(),
//			query,
//			metrics[i].ID,
//			metrics[i].MType,
//			metrics[i].Delta,
//			metrics[i].Value,
//			metrics[i].Hash,
//		)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//	err = tx.Commit(context.Background())
//	if err != nil {
//		log.Println(err)
//		tx.Rollback(context.Background())
//	}
//}
