// Code generated by MockGen. DO NOT EDIT.
// Source: internal/dao/order.go

// Package dao is a generated GoMock package.
package dao

import (
	models "dealer/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	gorm "gorm.io/gorm"
)

// MockOrderInterface is a mock of OrderInterface interface.
type MockOrderInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOrderInterfaceMockRecorder
}

// MockOrderInterfaceMockRecorder is the mock recorder for MockOrderInterface.
type MockOrderInterfaceMockRecorder struct {
	mock *MockOrderInterface
}

// NewMockOrderInterface creates a new mock instance.
func NewMockOrderInterface(ctrl *gomock.Controller) *MockOrderInterface {
	mock := &MockOrderInterface{ctrl: ctrl}
	mock.recorder = &MockOrderInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderInterface) EXPECT() *MockOrderInterfaceMockRecorder {
	return m.recorder
}

// BulkUpdate mocks base method.
func (m *MockOrderInterface) BulkUpdate(arg0 context.Context, arg1 *gorm.DB, arg2 []*models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkUpdate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkUpdate indicates an expected call of BulkUpdate.
func (mr *MockOrderInterfaceMockRecorder) BulkUpdate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkUpdate", reflect.TypeOf((*MockOrderInterface)(nil).BulkUpdate), arg0, arg1, arg2)
}

// Insert mocks base method.
func (m *MockOrderInterface) Insert(arg0 context.Context, arg1 *gorm.DB, arg2 *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockOrderInterfaceMockRecorder) Insert(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockOrderInterface)(nil).Insert), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockOrderInterface) Update(arg0 context.Context, arg1 *gorm.DB, arg2 *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockOrderInterfaceMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockOrderInterface)(nil).Update), arg0, arg1, arg2)
}
