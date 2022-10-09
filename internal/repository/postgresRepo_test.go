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
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "metric_service",
		},
		WaitingFor: wait.ForListeningPort(port),
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
			name:   "sould create connection to database and return repo",
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

			got, err := repo.GetAllMetrics()

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
		Conn  *sqlx.DB
		DBURI string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			got, err := repo.GetMetricByKey(tt.args.name)
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
		Conn  *sqlx.DB
		DBURI string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			if err := repo.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepo_SaveCounter(t *testing.T) {
	type fields struct {
		Conn  *sqlx.DB
		DBURI string
	}
	type args struct {
		name    string
		counter model.Counter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			if err := repo.SaveCounter(tt.args.name, tt.args.counter); (err != nil) != tt.wantErr {
				t.Errorf("SaveCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepo_SaveGauge(t *testing.T) {
	type fields struct {
		Conn  *sqlx.DB
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			if err := repo.SaveGauge(tt.args.name, tt.args.gauge); (err != nil) != tt.wantErr {
				t.Errorf("SaveGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepo_SaveMetricsBatch(t *testing.T) {
	type fields struct {
		Conn  *sqlx.DB
		DBURI string
	}
	type args struct {
		metrics []model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			if err := repo.SaveMetricsBatch(tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("SaveMetricsBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepo_Shutdown(t *testing.T) {
	type fields struct {
		Conn  *sqlx.DB
		DBURI string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &PostgresRepo{
				Conn:  tt.fields.Conn,
				DBURI: tt.fields.DBURI,
			}
			repo.Shutdown()
		})
	}
}
