package repository

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/yurchenkosv/metric-service/internal/model"
)

func initContainers(t *testing.T, ctx context.Context) testcontainers.Container {

	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		t.Error(err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{port.Port() + "/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "metric_service",
		},
		WaitingFor: wait.ForListeningPort(port),
		AutoRemove: true,
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	return postgres
}

func TestNewPostgresRepo(t *testing.T) {
	type args struct {
		dbURI string
	}
	tests := []struct {
		name   string
		args   args
		before func(t *testing.T, ctx context.Context) testcontainers.Container
		want   *PostgresRepo
	}{
		{
			name:   "should create connection to database and return repo",
			args:   args{},
			before: initContainers,
			want:   &PostgresRepo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.args.dbURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.args.dbURI)
			assert.IsType(t, tt.want, repo)
		})
	}
}

func TestPostgresRepo_GetAllMetrics(t *testing.T) {
	type fields struct {
		DBURI string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name            string
		fields          fields
		want            *model.Metrics
		before          func(t *testing.T, ctx context.Context) testcontainers.Container
		beforeCondition func(conn *sqlx.DB)
		wantErr         bool
	}{
		{
			name: "should sucessfuly return metrics",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			want: &model.Metrics{Metric: []model.Metric{
				{
					ID:    "RandomValue",
					MType: "counter",
					Delta: model.NewCounter(123),
					Value: nil,
				},
			}},
			before: initContainers,
			beforeCondition: func(conn *sqlx.DB) {
				qry := `
					INSERT INTO metrics(
					     id,
					     metric_id,
					     metric_type,
					     metric_delta
						) 
					VALUES (
							 1,
					        'RandomValue',
					        'counter',
					        123
					        )
				`
				_, err := conn.Exec(qry)
				if err != nil {
					fmt.Print(err)
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")
			tt.beforeCondition(repo.Conn)

			got, err := repo.GetAllMetrics(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_GetMetricByKey(t *testing.T) {
	type fields struct {
		DBURI string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		before          func(t *testing.T, ctx context.Context) testcontainers.Container
		beforeCondition func(conn *sqlx.DB)
		want            *model.Metric
		wantErr         bool
	}{
		{
			name: "should successfully return metric by name",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			args:   args{name: "RandomValue"},
			before: initContainers,
			beforeCondition: func(conn *sqlx.DB) {
				qry := `
					INSERT INTO metrics(
					     id,
					     metric_id,
					     metric_type,
					     metric_delta
						) 
					VALUES (
							 1,
					        'RandomValue',
					        'counter',
					        123
					        )
				`
				_, err := conn.Exec(qry)
				if err != nil {
					fmt.Print(err)
				}
			},
			want: &model.Metric{
				ID:    "RandomValue",
				MType: "counter",
				Delta: model.NewCounter(123),
				Value: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")
			tt.beforeCondition(repo.Conn)

			got, err := repo.GetMetricByKey(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricByKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_Ping(t *testing.T) {
	type fields struct {
		DBURI string
	}
	tests := []struct {
		name    string
		before  func(t *testing.T, ctx context.Context) testcontainers.Container
		fields  fields
		wantErr bool
	}{
		{
			name:   "should successfully send ping to database",
			before: initContainers,
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			if err := repo.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepo_SaveCounter(t *testing.T) {
	type fields struct {
		DBURI string
	}
	type args struct {
		name    string
		counter model.Counter
	}
	tests := []struct {
		name    string
		fields  fields
		before  func(t *testing.T, ctx context.Context) testcontainers.Container
		after   func(conn *sqlx.DB, metricID string) *model.Counter
		args    args
		wantErr bool
	}{
		{
			name: "should sucessfully save counter in db",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			before: initContainers,
			args: args{
				name:    "RandomValue",
				counter: 123,
			},
			wantErr: false,
			after: func(conn *sqlx.DB, metricID string) *model.Counter {
				var counter model.Counter
				qry := "SELECT metric_delta FROM metrics WHERE metric_id=$1"
				result := conn.QueryRow(qry, metricID)
				err := result.Scan(&counter)
				if err != nil {
					fmt.Print(err)
				}
				return &counter
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")

			if err := repo.SaveCounter(ctx, tt.args.name, tt.args.counter); (err != nil) != tt.wantErr {
				t.Errorf("SaveCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if counter := (*tt.after(repo.Conn, tt.args.name)); counter != tt.args.counter {
				t.Errorf("result in db was %v, expected %v", counter, tt.args.counter)
			}
		})
	}
}

func TestPostgresRepo_SaveGauge(t *testing.T) {
	type fields struct {
		DBURI string
	}
	type args struct {
		name  string
		gauge model.Gauge
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		before  func(t *testing.T, ctx context.Context) testcontainers.Container
		after   func(conn *sqlx.DB, metricID string) *model.Gauge
		wantErr bool
	}{
		{
			name: "should sucessfuly save gauge metric",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			args: args{
				name:  "RandomGauge",
				gauge: 500.123,
			},
			before: initContainers,
			after: func(conn *sqlx.DB, metricID string) *model.Gauge {
				var gauge model.Gauge
				qry := "SELECT metric_value FROM metrics WHERE metric_id=$1"
				result := conn.QueryRow(qry, metricID)
				err := result.Scan(&gauge)
				if err != nil {
					fmt.Print(err)
				}
				return &gauge

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")

			if err := repo.SaveGauge(ctx, tt.args.name, tt.args.gauge); (err != nil) != tt.wantErr {
				t.Errorf("SaveGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gauge := (*tt.after(repo.Conn, tt.args.name)); gauge != tt.args.gauge {
				t.Errorf("result in db was %v, expected %v", gauge, tt.args.gauge)
			}
		})
	}
}

func TestPostgresRepo_SaveMetricsBatch(t *testing.T) {
	type fields struct {
		DBURI string
	}
	type args struct {
		metrics []model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		before  func(t *testing.T, ctx context.Context) testcontainers.Container
		after   func(conn *sqlx.DB) []model.Metric
		wantErr bool
		want    []model.Metric
	}{
		{
			name: "shold successfully save metric batch",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			args: args{
				metrics: []model.Metric{
					{
						ID:    "RandomCounter",
						MType: "counter",
						Delta: model.NewCounter(100),
					},
					{
						ID:    "RandomGauge",
						MType: "gauge",
						Value: model.NewGauge(500.25),
					},
				},
			},
			before: initContainers,
			after: func(conn *sqlx.DB) []model.Metric {
				var metrics []model.Metric
				qry := "SELECT metric_id, metric_type, metric_delta, metric_value FROM metrics WHERE true"
				rows, err := conn.Query(qry)
				if err != nil {
					fmt.Println(err)
				}
				err = rows.Err()
				if err != nil {
					fmt.Println(err)
				}
				for rows.Next() {
					var (
						metricID    string
						metricType  string
						metricDelta *model.Counter
						metricValue *model.Gauge
						metric      model.Metric
					)

					err2 := rows.Scan(&metricID, &metricType, &metricDelta, &metricValue)
					if err2 != nil {
						fmt.Println(err2)
					}
					metric.ID = metricID
					metric.MType = metricType
					metric.Delta = metricDelta
					metric.Value = metricValue
					metrics = append(metrics, metric)
				}
				return metrics
			},
			wantErr: false,
			want: []model.Metric{
				{
					ID:    "RandomCounter",
					MType: "counter",
					Delta: model.NewCounter(100),
				},
				{
					ID:    "RandomGauge",
					MType: "gauge",
					Value: model.NewGauge(500.25),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")
			if err := repo.SaveMetricsBatch(ctx, tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("SaveMetricsBatch() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := tt.after(repo.Conn)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricByKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_Shutdown(t *testing.T) {
	type fields struct {
		DBURI string
	}
	tests := []struct {
		name   string
		fields fields
		before func(t *testing.T, ctx context.Context) testcontainers.Container
	}{
		{
			name: "should call shutdown method",
			fields: fields{
				DBURI: "postgresql://postgres:postgres@%s/metric_service?sslmode=disable",
			},
			before: initContainers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			postgres := tt.before(t, ctx)
			defer postgres.Terminate(ctx)
			endpoint, err := postgres.Endpoint(ctx, "")
			if err != nil {
				t.Error(err)
			}
			tt.fields.DBURI = fmt.Sprintf("postgresql://postgres:postgres@%s/metric_service?sslmode=disable", endpoint)
			repo := NewPostgresRepo(tt.fields.DBURI)
			repo.Migrate("../../db/migrations")

			repo.Shutdown()
		})
	}
}
