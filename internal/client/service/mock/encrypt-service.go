// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blokhinnv/gophkeeper/internal/client/service (interfaces: EncryptService)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/blokhinnv/gophkeeper/internal/client/models"
	gomock "github.com/golang/mock/gomock"
)

// MockEncryptService is a mock of EncryptService interface.
type MockEncryptService struct {
	ctrl     *gomock.Controller
	recorder *MockEncryptServiceMockRecorder
}

// MockEncryptServiceMockRecorder is the mock recorder for MockEncryptService.
type MockEncryptServiceMockRecorder struct {
	mock *MockEncryptService
}

// NewMockEncryptService creates a new mock instance.
func NewMockEncryptService(ctrl *gomock.Controller) *MockEncryptService {
	mock := &MockEncryptService{ctrl: ctrl}
	mock.recorder = &MockEncryptServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncryptService) EXPECT() *MockEncryptServiceMockRecorder {
	return m.recorder
}

// FromEncryptedFile mocks base method.
func (m *MockEncryptService) FromEncryptedFile(arg0, arg1 string) (*models.SyncResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromEncryptedFile", arg0, arg1)
	ret0, _ := ret[0].(*models.SyncResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FromEncryptedFile indicates an expected call of FromEncryptedFile.
func (mr *MockEncryptServiceMockRecorder) FromEncryptedFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromEncryptedFile", reflect.TypeOf((*MockEncryptService)(nil).FromEncryptedFile), arg0, arg1)
}

// ToEncryptedFile mocks base method.
func (m *MockEncryptService) ToEncryptedFile(arg0 *models.SyncResponse, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToEncryptedFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ToEncryptedFile indicates an expected call of ToEncryptedFile.
func (mr *MockEncryptServiceMockRecorder) ToEncryptedFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToEncryptedFile", reflect.TypeOf((*MockEncryptService)(nil).ToEncryptedFile), arg0, arg1, arg2)
}
