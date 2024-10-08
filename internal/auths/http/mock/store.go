// Code generated by MockGen. DO NOT EDIT.
// Source: go-invoice/internal/auths/usecase (interfaces: AuthUsecase)
//
// Generated by this command:
//
//	mockgen -package mockAuthUse -destination internal/auths/http/mock/store.go go-invoice/internal/auths/usecase AuthUsecase
//

// Package mockAuthUse is a generated GoMock package.
package mockAuthUse

import (
	domain "go-invoice/domain"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthUsecase is a mock of AuthUsecase interface.
type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

// MockAuthUsecaseMockRecorder is the mock recorder for MockAuthUsecase.
type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

// NewMockAuthUsecase creates a new mock instance.
func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	mock := &MockAuthUsecase{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthUsecase) CreateUser(arg0 domain.CreateUserRequestDto) (int, gin.H, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(gin.H)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthUsecaseMockRecorder) CreateUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthUsecase)(nil).CreateUser), arg0)
}

// FetchUser mocks base method.
func (m *MockAuthUsecase) FetchUser(arg0 string) (int, domain.UserReponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUser", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(domain.UserReponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FetchUser indicates an expected call of FetchUser.
func (mr *MockAuthUsecaseMockRecorder) FetchUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUser", reflect.TypeOf((*MockAuthUsecase)(nil).FetchUser), arg0)
}

// LoginUser mocks base method.
func (m *MockAuthUsecase) LoginUser(arg0 domain.LoginRequestDto) (int, gin.H, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(gin.H)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockAuthUsecaseMockRecorder) LoginUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockAuthUsecase)(nil).LoginUser), arg0)
}
