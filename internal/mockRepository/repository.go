// Code generated by MockGen. DO NOT EDIT.
// Source: .\internal\repository\repository.go

// Package mock_repository is a generated GoMock package.
package mockRepository

import (
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

// GetAllMetrics mockRepository base method.
func (m *MockRepository) GetAllMetrics() (*model.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics")
	ret0, _ := ret[0].(*model.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockRepositoryMockRecorder) GetAllMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockRepository)(nil).GetAllMetrics))
}

// GetMetricByKey mockRepository base method.
func (m *MockRepository) GetMetricByKey(arg0 string) (*model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByKey", arg0)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByKey indicates an expected call of GetMetricByKey.
func (mr *MockRepositoryMockRecorder) GetMetricByKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByKey", reflect.TypeOf((*MockRepository)(nil).GetMetricByKey), arg0)
}

// Ping mockRepository base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}

// SaveCounter mockRepository base method.
func (m *MockRepository) SaveCounter(arg0 string, arg1 model.Counter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCounter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCounter indicates an expected call of SaveCounter.
func (mr *MockRepositoryMockRecorder) SaveCounter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCounter", reflect.TypeOf((*MockRepository)(nil).SaveCounter), arg0, arg1)
}

// SaveGauge mockRepository base method.
func (m *MockRepository) SaveGauge(arg0 string, arg1 model.Gauge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveGauge", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveGauge indicates an expected call of SaveGauge.
func (mr *MockRepositoryMockRecorder) SaveGauge(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGauge", reflect.TypeOf((*MockRepository)(nil).SaveGauge), arg0, arg1)
}

// SaveMetricsBatch mockRepository base method.
func (m *MockRepository) SaveMetricsBatch(arg0 []model.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetricsBatch", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetricsBatch indicates an expected call of SaveMetricsBatch.
func (mr *MockRepositoryMockRecorder) SaveMetricsBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetricsBatch", reflect.TypeOf((*MockRepository)(nil).SaveMetricsBatch), arg0)
}
