// Code generated by MockGen. DO NOT EDIT.
// Source: os.go

// Package os is a generated GoMock package.
package os

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
	os0 "os"
	reflect "reflect"
)

// MockOS is a mock of OS interface
type MockOS struct {
	ctrl     *gomock.Controller
	recorder *MockOSMockRecorder
}

// MockOSMockRecorder is the mock recorder for MockOS
type MockOSMockRecorder struct {
	mock *MockOS
}

// NewMockOS creates a new mock instance
func NewMockOS(ctrl *gomock.Controller) *MockOS {
	mock := &MockOS{ctrl: ctrl}
	mock.recorder = &MockOSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOS) EXPECT() *MockOSMockRecorder {
	return m.recorder
}

// Args mocks base method
func (m *MockOS) Args() []string {
	ret := m.ctrl.Call(m, "Args")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Args indicates an expected call of Args
func (mr *MockOSMockRecorder) Args() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Args", reflect.TypeOf((*MockOS)(nil).Args))
}

// Create mocks base method
func (m *MockOS) Create(name string) (*os0.File, error) {
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(*os0.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockOSMockRecorder) Create(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOS)(nil).Create), name)
}

// Getenv mocks base method
func (m *MockOS) Getenv(key string) string {
	ret := m.ctrl.Call(m, "Getenv", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// Getenv indicates an expected call of Getenv
func (mr *MockOSMockRecorder) Getenv(key interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Getenv", reflect.TypeOf((*MockOS)(nil).Getenv), key)
}

// LookupEnv mocks base method
func (m *MockOS) LookupEnv(key string) (string, bool) {
	ret := m.ctrl.Call(m, "LookupEnv", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// LookupEnv indicates an expected call of LookupEnv
func (mr *MockOSMockRecorder) LookupEnv(key interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupEnv", reflect.TypeOf((*MockOS)(nil).LookupEnv), key)
}

// ExpandEnv mocks base method
func (m *MockOS) ExpandEnv(s string) string {
	ret := m.ctrl.Call(m, "ExpandEnv", s)
	ret0, _ := ret[0].(string)
	return ret0
}

// ExpandEnv indicates an expected call of ExpandEnv
func (mr *MockOSMockRecorder) ExpandEnv(s interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpandEnv", reflect.TypeOf((*MockOS)(nil).ExpandEnv), s)
}

// Setenv mocks base method
func (m *MockOS) Setenv(key, value string) error {
	ret := m.ctrl.Call(m, "Setenv", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setenv indicates an expected call of Setenv
func (mr *MockOSMockRecorder) Setenv(key, value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setenv", reflect.TypeOf((*MockOS)(nil).Setenv), key, value)
}

// Exit mocks base method
func (m *MockOS) Exit(code int) {
	m.ctrl.Call(m, "Exit", code)
}

// Exit indicates an expected call of Exit
func (mr *MockOSMockRecorder) Exit(code interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exit", reflect.TypeOf((*MockOS)(nil).Exit), code)
}

// Stdin mocks base method
func (m *MockOS) Stdin() io.Reader {
	ret := m.ctrl.Call(m, "Stdin")
	ret0, _ := ret[0].(io.Reader)
	return ret0
}

// Stdin indicates an expected call of Stdin
func (mr *MockOSMockRecorder) Stdin() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stdin", reflect.TypeOf((*MockOS)(nil).Stdin))
}

// Stdout mocks base method
func (m *MockOS) Stdout() io.Writer {
	ret := m.ctrl.Call(m, "Stdout")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

// Stdout indicates an expected call of Stdout
func (mr *MockOSMockRecorder) Stdout() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stdout", reflect.TypeOf((*MockOS)(nil).Stdout))
}

// Stderr mocks base method
func (m *MockOS) Stderr() io.Writer {
	ret := m.ctrl.Call(m, "Stderr")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

// Stderr indicates an expected call of Stderr
func (mr *MockOSMockRecorder) Stderr() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stderr", reflect.TypeOf((*MockOS)(nil).Stderr))
}

// Stat mocks base method
func (m *MockOS) Stat(name string) (os0.FileInfo, error) {
	ret := m.ctrl.Call(m, "Stat", name)
	ret0, _ := ret[0].(os0.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat
func (mr *MockOSMockRecorder) Stat(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockOS)(nil).Stat), name)
}

// IsNotExist mocks base method
func (m *MockOS) IsNotExist(err error) bool {
	ret := m.ctrl.Call(m, "IsNotExist", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNotExist indicates an expected call of IsNotExist
func (mr *MockOSMockRecorder) IsNotExist(err interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNotExist", reflect.TypeOf((*MockOS)(nil).IsNotExist), err)
}

// Remove mocks base method
func (m *MockOS) Remove(name string) error {
	ret := m.ctrl.Call(m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockOSMockRecorder) Remove(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockOS)(nil).Remove), name)
}

// Rename mocks base method
func (m *MockOS) Rename(oldpath, newpath string) error {
	ret := m.ctrl.Call(m, "Rename", oldpath, newpath)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rename indicates an expected call of Rename
func (mr *MockOSMockRecorder) Rename(oldpath, newpath interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rename", reflect.TypeOf((*MockOS)(nil).Rename), oldpath, newpath)
}

// Open mocks base method
func (m *MockOS) Open(name string) (*os0.File, error) {
	ret := m.ctrl.Call(m, "Open", name)
	ret0, _ := ret[0].(*os0.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open
func (mr *MockOSMockRecorder) Open(name interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockOS)(nil).Open), name)
}

// OpenFile mocks base method
func (m *MockOS) OpenFile(name string, flag int, perm os0.FileMode) (*os0.File, error) {
	ret := m.ctrl.Call(m, "OpenFile", name, flag, perm)
	ret0, _ := ret[0].(*os0.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFile indicates an expected call of OpenFile
func (mr *MockOSMockRecorder) OpenFile(name, flag, perm interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFile", reflect.TypeOf((*MockOS)(nil).OpenFile), name, flag, perm)
}

// MkdirAll mocks base method
func (m *MockOS) MkdirAll(path string, perm os0.FileMode) error {
	ret := m.ctrl.Call(m, "MkdirAll", path, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// MkdirAll indicates an expected call of MkdirAll
func (mr *MockOSMockRecorder) MkdirAll(path, perm interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*MockOS)(nil).MkdirAll), path, perm)
}

// NewDirReader mocks base method
func (m *MockOS) NewDirReader(dir string) DirReader {
	ret := m.ctrl.Call(m, "NewDirReader", dir)
	ret0, _ := ret[0].(DirReader)
	return ret0
}

// NewDirReader indicates an expected call of NewDirReader
func (mr *MockOSMockRecorder) NewDirReader(dir interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDirReader", reflect.TypeOf((*MockOS)(nil).NewDirReader), dir)
}
