// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iostrovok/go-mockdata/test-code (interfaces: IOne)

// Package mmock is a generated GoMock package.
package mmock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIOne is a mock of IOne interface
type MockIOne struct {
	ctrl     *gomock.Controller
	recorder *MockIOneMockRecorder
}

// MockIOneMockRecorder is the mock recorder for MockIOne
type MockIOneMockRecorder struct {
	mock *MockIOne
}

// NewMockIOne creates a new mock instance
func NewMockIOne(ctrl *gomock.Controller) *MockIOne {
	mock := &MockIOne{ctrl: ctrl}
	mock.recorder = &MockIOneMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIOne) EXPECT() *MockIOneMockRecorder {
	return m.recorder
}

// FirstFunc mocks base method
func (m *MockIOne) FirstFunc(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstFunc", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FirstFunc indicates an expected call of FirstFunc
func (mr *MockIOneMockRecorder) FirstFunc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstFunc", reflect.TypeOf((*MockIOne)(nil).FirstFunc), arg0)
}

// SecondFunc mocks base method
func (m *MockIOne) SecondFunc(arg0 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecondFunc", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SecondFunc indicates an expected call of SecondFunc
func (mr *MockIOneMockRecorder) SecondFunc(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecondFunc", reflect.TypeOf((*MockIOne)(nil).SecondFunc), arg0)
}
