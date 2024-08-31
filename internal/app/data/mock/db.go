// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/data/db.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

// MockDBPooler is a mock of DBPooler interface.
type MockDBPooler struct {
	ctrl     *gomock.Controller
	recorder *MockDBPoolerMockRecorder
}

// MockDBPoolerMockRecorder is the mock recorder for MockDBPooler.
type MockDBPoolerMockRecorder struct {
	mock *MockDBPooler
}

// NewMockDBPooler creates a new mock instance.
func NewMockDBPooler(ctrl *gomock.Controller) *MockDBPooler {
	mock := &MockDBPooler{ctrl: ctrl}
	mock.recorder = &MockDBPoolerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBPooler) EXPECT() *MockDBPoolerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDBPooler) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockDBPoolerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDBPooler)(nil).Close))
}

// Exec mocks base method.
func (m *MockDBPooler) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range arguments {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDBPoolerMockRecorder) Exec(ctx, sql interface{}, arguments ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, arguments...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDBPooler)(nil).Exec), varargs...)
}

// Ping mocks base method.
func (m *MockDBPooler) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDBPoolerMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDBPooler)(nil).Ping), ctx)
}

// Query mocks base method.
func (m *MockDBPooler) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockDBPoolerMockRecorder) Query(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDBPooler)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockDBPooler) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockDBPoolerMockRecorder) QueryRow(ctx, sql interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockDBPooler)(nil).QueryRow), varargs...)
}

// SendBatch mocks base method.
func (m *MockDBPooler) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBatch", ctx, b)
	ret0, _ := ret[0].(pgx.BatchResults)
	return ret0
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockDBPoolerMockRecorder) SendBatch(ctx, b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockDBPooler)(nil).SendBatch), ctx, b)
}
