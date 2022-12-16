// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	auth "github.com/rusq/slackdump/v2/auth"
)

// MockCredentials is a mock of Credentials interface.
type MockCredentials struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialsMockRecorder
}

// MockCredentialsMockRecorder is the mock recorder for MockCredentials.
type MockCredentialsMockRecorder struct {
	mock *MockCredentials
}

// NewMockCredentials creates a new mock instance.
func NewMockCredentials(ctrl *gomock.Controller) *MockCredentials {
	mock := &MockCredentials{ctrl: ctrl}
	mock.recorder = &MockCredentialsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentials) EXPECT() *MockCredentialsMockRecorder {
	return m.recorder
}

// AuthProvider mocks base method.
func (m *MockCredentials) AuthProvider(ctx context.Context, workspace string, opts ...auth.Option) (auth.Provider, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, workspace}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AuthProvider", varargs...)
	ret0, _ := ret[0].(auth.Provider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthProvider indicates an expected call of AuthProvider.
func (mr *MockCredentialsMockRecorder) AuthProvider(ctx, workspace interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, workspace}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthProvider", reflect.TypeOf((*MockCredentials)(nil).AuthProvider), varargs...)
}

// IsEmpty mocks base method.
func (m *MockCredentials) IsEmpty() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmpty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsEmpty indicates an expected call of IsEmpty.
func (mr *MockCredentialsMockRecorder) IsEmpty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmpty", reflect.TypeOf((*MockCredentials)(nil).IsEmpty))
}

// Mockcontainer is a mock of container interface.
type Mockcontainer struct {
	ctrl     *gomock.Controller
	recorder *MockcontainerMockRecorder
}

// MockcontainerMockRecorder is the mock recorder for Mockcontainer.
type MockcontainerMockRecorder struct {
	mock *Mockcontainer
}

// NewMockcontainer creates a new mock instance.
func NewMockcontainer(ctrl *gomock.Controller) *Mockcontainer {
	mock := &Mockcontainer{ctrl: ctrl}
	mock.recorder = &MockcontainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockcontainer) EXPECT() *MockcontainerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockcontainer) Create(filename string) (io.WriteCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", filename)
	ret0, _ := ret[0].(io.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockcontainerMockRecorder) Create(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockcontainer)(nil).Create), filename)
}

// Open mocks base method.
func (m *Mockcontainer) Open(filename string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", filename)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockcontainerMockRecorder) Open(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*Mockcontainer)(nil).Open), filename)
}
