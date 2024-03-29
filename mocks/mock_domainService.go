// Code generated by MockGen. DO NOT EDIT.
// Source: microservice/internal/app/drivers/rest (interfaces: DomainSvc)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "microservice/models"
	reflect "reflect"
)

// MockDomainService is a mock of DomainSvc interface
type MockDomainService struct {
	ctrl     *gomock.Controller
	recorder *MockDomainServiceMockRecorder
}

// MockDomainServiceMockRecorder is the mock recorder for MockDomainService
type MockDomainServiceMockRecorder struct {
	mock *MockDomainService
}

// NewMockDomainService creates a new mock instance
func NewMockDomainService(ctrl *gomock.Controller) *MockDomainService {
	mock := &MockDomainService{ctrl: ctrl}
	mock.recorder = &MockDomainServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDomainService) EXPECT() *MockDomainServiceMockRecorder {
	return m.recorder
}

// AddDocument mocks base method
func (m *MockDomainService) AddDocument(arg0 context.Context, arg1 models.Document) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDocument", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDocument indicates an expected call of AddDocument
func (mr *MockDomainServiceMockRecorder) AddDocument(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDocument", reflect.TypeOf((*MockDomainService)(nil).AddDocument), arg0, arg1)
}

// GetDocument mocks base method
func (m *MockDomainService) GetDocument(arg0 context.Context, arg1 string) (models.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDocument", arg0, arg1)
	ret0, _ := ret[0].(models.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDocument indicates an expected call of GetDocument
func (mr *MockDomainServiceMockRecorder) GetDocument(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDocument", reflect.TypeOf((*MockDomainService)(nil).GetDocument), arg0, arg1)
}

// Teardown mocks base method
func (m *MockDomainService) Teardown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Teardown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Teardown indicates an expected call of Teardown
func (mr *MockDomainServiceMockRecorder) Teardown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Teardown", reflect.TypeOf((*MockDomainService)(nil).Teardown), arg0)
}
