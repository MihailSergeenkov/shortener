// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jackc/pgx/v5 (interfaces: BatchResults)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

// MockBatchResults is a mock of BatchResults interface.
type MockBatchResults struct {
	ctrl     *gomock.Controller
	recorder *MockBatchResultsMockRecorder
}

// MockBatchResultsMockRecorder is the mock recorder for MockBatchResults.
type MockBatchResultsMockRecorder struct {
	mock *MockBatchResults
}

// NewMockBatchResults creates a new mock instance.
func NewMockBatchResults(ctrl *gomock.Controller) *MockBatchResults {
	mock := &MockBatchResults{ctrl: ctrl}
	mock.recorder = &MockBatchResultsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchResults) EXPECT() *MockBatchResultsMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockBatchResults) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockBatchResultsMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBatchResults)(nil).Close))
}

// Exec mocks base method.
func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec")
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockBatchResultsMockRecorder) Exec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockBatchResults)(nil).Exec))
}

// Query mocks base method.
func (m *MockBatchResults) Query() (pgx.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query")
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockBatchResultsMockRecorder) Query() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockBatchResults)(nil).Query))
}

// QueryRow mocks base method.
func (m *MockBatchResults) QueryRow() pgx.Row {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryRow")
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockBatchResultsMockRecorder) QueryRow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockBatchResults)(nil).QueryRow))
}