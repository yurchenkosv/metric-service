// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\repository\repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/yurchenkosv/metric-service/internal/model"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAllMetrics mocks base method.
func (m *MockRepository) GetAllMetrics(ctx context.Context) (*model.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics", ctx)
	ret0, _ := ret[0].(*model.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockRepositoryMockRecorder) GetAllMetrics(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockRepository)(nil).GetAllMetrics), ctx)
}

// GetMetricByKey mocks base method.
func (m *MockRepository) GetMetricByKey(arg0 string, arg1 context.Context) (*model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByKey", arg0, arg1)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByKey indicates an expected call of GetMetricByKey.
func (mr *MockRepositoryMockRecorder) GetMetricByKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByKey", reflect.TypeOf((*MockRepository)(nil).GetMetricByKey), arg0, arg1)
}

// Migrate mocks base method.
func (m *MockRepository) Migrate(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Migrate", arg0)
}

// Migrate indicates an expected call of Migrate.
func (mr *MockRepositoryMockRecorder) Migrate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockRepository)(nil).Migrate), arg0)
}

// Ping mocks base method.
func (m *MockRepository) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping), ctx)
}

// SaveCounter mocks base method.
func (m *MockRepository) SaveCounter(arg0 string, arg1 model.Counter, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCounter", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCounter indicates an expected call of SaveCounter.
func (mr *MockRepositoryMockRecorder) SaveCounter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCounter", reflect.TypeOf((*MockRepository)(nil).SaveCounter), arg0, arg1, arg2)
}

// SaveGauge mocks base method.
func (m *MockRepository) SaveGauge(arg0 string, arg1 model.Gauge, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveGauge", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveGauge indicates an expected call of SaveGauge.
func (mr *MockRepositoryMockRecorder) SaveGauge(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGauge", reflect.TypeOf((*MockRepository)(nil).SaveGauge), arg0, arg1, arg2)
}

// SaveMetricsBatch mocks base method.
func (m *MockRepository) SaveMetricsBatch(arg0 []model.Metric, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetricsBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetricsBatch indicates an expected call of SaveMetricsBatch.
func (mr *MockRepositoryMockRecorder) SaveMetricsBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetricsBatch", reflect.TypeOf((*MockRepository)(nil).SaveMetricsBatch), arg0, arg1)
}

// Shutdown mocks base method.
func (m *MockRepository) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRepositoryMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRepository)(nil).Shutdown))
}
