// Automatically generated by MockGen. DO NOT EDIT!
// Source: os.go

package os

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
	os "os"
)

// Mock of OS interface
type MockOS struct {
	ctrl     *gomock.Controller
	recorder *_MockOSRecorder
}

// Recorder for MockOS (not exported)
type _MockOSRecorder struct {
	mock *MockOS
}

func NewMockOS(ctrl *gomock.Controller) *MockOS {
	mock := &MockOS{ctrl: ctrl}
	mock.recorder = &_MockOSRecorder{mock}
	return mock
}

func (_m *MockOS) EXPECT() *_MockOSRecorder {
	return _m.recorder
}

func (_m *MockOS) Args() []string {
	ret := _m.ctrl.Call(_m, "Args")
	ret0, _ := ret[0].([]string)
	return ret0
}

func (_mr *_MockOSRecorder) Args() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Args")
}

func (_m *MockOS) Getenv(key string) string {
	ret := _m.ctrl.Call(_m, "Getenv", key)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockOSRecorder) Getenv(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Getenv", arg0)
}

func (_m *MockOS) LookupEnv(key string) (string, bool) {
	ret := _m.ctrl.Call(_m, "LookupEnv", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

func (_mr *_MockOSRecorder) LookupEnv(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "LookupEnv", arg0)
}

func (_m *MockOS) Exit(code int) {
	_m.ctrl.Call(_m, "Exit", code)
}

func (_mr *_MockOSRecorder) Exit(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exit", arg0)
}

func (_m *MockOS) Stdin() io.Reader {
	ret := _m.ctrl.Call(_m, "Stdin")
	ret0, _ := ret[0].(io.Reader)
	return ret0
}

func (_mr *_MockOSRecorder) Stdin() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stdin")
}

func (_m *MockOS) Stdout() io.Writer {
	ret := _m.ctrl.Call(_m, "Stdout")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

func (_mr *_MockOSRecorder) Stdout() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stdout")
}

func (_m *MockOS) Stderr() io.Writer {
	ret := _m.ctrl.Call(_m, "Stderr")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

func (_mr *_MockOSRecorder) Stderr() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stderr")
}

func (_m *MockOS) Stat(name string) (os.FileInfo, error) {
	ret := _m.ctrl.Call(_m, "Stat", name)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOSRecorder) Stat(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stat", arg0)
}

func (_m *MockOS) IsNotExist(err error) bool {
	ret := _m.ctrl.Call(_m, "IsNotExist", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockOSRecorder) IsNotExist(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsNotExist", arg0)
}

func (_m *MockOS) Remove(name string) error {
	ret := _m.ctrl.Call(_m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOSRecorder) Remove(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Remove", arg0)
}

func (_m *MockOS) Rename(oldpath string, newpath string) error {
	ret := _m.ctrl.Call(_m, "Rename", oldpath, newpath)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockOSRecorder) Rename(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Rename", arg0, arg1)
}

func (_m *MockOS) Open(name string) (*os.File, error) {
	ret := _m.ctrl.Call(_m, "Open", name)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOSRecorder) Open(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Open", arg0)
}

func (_m *MockOS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	ret := _m.ctrl.Call(_m, "OpenFile", name, flag, perm)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockOSRecorder) OpenFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "OpenFile", arg0, arg1, arg2)
}

func (_m *MockOS) NewDirReader(dir string) DirReader {
	ret := _m.ctrl.Call(_m, "NewDirReader", dir)
	ret0, _ := ret[0].(DirReader)
	return ret0
}

func (_mr *_MockOSRecorder) NewDirReader(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewDirReader", arg0)
}
