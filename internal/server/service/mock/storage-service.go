// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blokhinnv/gophkeeper/internal/server/service (interfaces: StorageService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/blokhinnv/gophkeeper/internal/server/models"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockStorageService is a mock of StorageService interface.
type MockStorageService struct {
	ctrl     *gomock.Controller
	recorder *MockStorageServiceMockRecorder
}

// MockStorageServiceMockRecorder is the mock recorder for MockStorageService.
type MockStorageServiceMockRecorder struct {
	mock *MockStorageService
}

// NewMockStorageService creates a new mock instance.
func NewMockStorageService(ctrl *gomock.Controller) *MockStorageService {
	mock := &MockStorageService{ctrl: ctrl}
	mock.recorder = &MockStorageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageService) EXPECT() *MockStorageServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockStorageService) Delete(arg0 context.Context, arg1 models.CollectionName, arg2 string, arg3 primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageServiceMockRecorder) Delete(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorageService)(nil).Delete), arg0, arg1, arg2, arg3)
}

// GetAll mocks base method.
func (m *MockStorageService) GetAll(arg0 context.Context, arg1 models.CollectionName, arg2 string) ([]models.UntypedRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.UntypedRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorageServiceMockRecorder) GetAll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorageService)(nil).GetAll), arg0, arg1, arg2)
}

// Store mocks base method.
func (m *MockStorageService) Store(arg0 context.Context, arg1 models.CollectionName, arg2 models.UntypedRecord) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockStorageServiceMockRecorder) Store(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockStorageService)(nil).Store), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockStorageService) Update(arg0 context.Context, arg1 models.CollectionName, arg2 string, arg3 primitive.ObjectID, arg4 interface{}, arg5 models.Metadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockStorageServiceMockRecorder) Update(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorageService)(nil).Update), arg0, arg1, arg2, arg3, arg4, arg5)
}