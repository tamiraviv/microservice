// Code generated by MockGen. DO NOT EDIT.
// Source: microservice/internal/app (interfaces: Configuration)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockConfigurationService is a mock of Configuration interface
type MockConfigurationService struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationServiceMockRecorder
}

// MockConfigurationServiceMockRecorder is the mock recorder for MockConfigurationService
type MockConfigurationServiceMockRecorder struct {
	mock *MockConfigurationService
}

// NewMockConfigurationService creates a new mock instance
func NewMockConfigurationService(ctrl *gomock.Controller) *MockConfigurationService {
	mock := &MockConfigurationService{ctrl: ctrl}
	mock.recorder = &MockConfigurationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurationService) EXPECT() *MockConfigurationServiceMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockConfigurationService) Get(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockConfigurationServiceMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigurationService)(nil).Get), arg0)
}

// GetBool mocks base method
func (m *MockConfigurationService) GetBool(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBool", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBool indicates an expected call of GetBool
func (mr *MockConfigurationServiceMockRecorder) GetBool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBool", reflect.TypeOf((*MockConfigurationService)(nil).GetBool), arg0)
}

// GetDuration mocks base method
func (m *MockConfigurationService) GetDuration(arg0 string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDuration", arg0)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDuration indicates an expected call of GetDuration
func (mr *MockConfigurationServiceMockRecorder) GetDuration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDuration", reflect.TypeOf((*MockConfigurationService)(nil).GetDuration), arg0)
}

// GetInt mocks base method
func (m *MockConfigurationService) GetInt(arg0 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt indicates an expected call of GetInt
func (mr *MockConfigurationServiceMockRecorder) GetInt(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockConfigurationService)(nil).GetInt), arg0)
}

// GetString mocks base method
func (m *MockConfigurationService) GetString(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetString indicates an expected call of GetString
func (mr *MockConfigurationServiceMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockConfigurationService)(nil).GetString), arg0)
}

// IsSet mocks base method
func (m *MockConfigurationService) IsSet(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSet", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSet indicates an expected call of IsSet
func (mr *MockConfigurationServiceMockRecorder) IsSet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSet", reflect.TypeOf((*MockConfigurationService)(nil).IsSet), arg0)
}

// Set mocks base method
func (m *MockConfigurationService) Set(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockConfigurationServiceMockRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockConfigurationService)(nil).Set), arg0, arg1)
}
