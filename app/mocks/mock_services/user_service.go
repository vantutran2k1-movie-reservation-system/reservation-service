// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/user_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/user_service.go -destination=app/mocks/mock_services/user_service.go -package=mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	errors "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	payloads "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreatePasswordResetToken mocks base method.
func (m *MockUserService) CreatePasswordResetToken(req payloads.CreatePasswordResetTokenRequest) (*models.PasswordResetToken, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePasswordResetToken", req)
	ret0, _ := ret[0].(*models.PasswordResetToken)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreatePasswordResetToken indicates an expected call of CreatePasswordResetToken.
func (mr *MockUserServiceMockRecorder) CreatePasswordResetToken(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePasswordResetToken", reflect.TypeOf((*MockUserService)(nil).CreatePasswordResetToken), req)
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(req payloads.CreateUserRequest) (*models.User, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", req)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), req)
}

// GetUser mocks base method.
func (m *MockUserService) GetUser(id uuid.UUID, includeProfile bool) (*models.User, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id, includeProfile)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserServiceMockRecorder) GetUser(id, includeProfile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserService)(nil).GetUser), id, includeProfile)
}

// LoginUser mocks base method.
func (m *MockUserService) LoginUser(req payloads.LoginUserRequest) (*models.LoginToken, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", req)
	ret0, _ := ret[0].(*models.LoginToken)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserServiceMockRecorder) LoginUser(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserService)(nil).LoginUser), req)
}

// LogoutUser mocks base method.
func (m *MockUserService) LogoutUser(tokenValue string) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", tokenValue)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser.
func (mr *MockUserServiceMockRecorder) LogoutUser(tokenValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockUserService)(nil).LogoutUser), tokenValue)
}

// ResetUserPassword mocks base method.
func (m *MockUserService) ResetUserPassword(resetToken string, request payloads.ResetPasswordRequest) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetUserPassword", resetToken, request)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// ResetUserPassword indicates an expected call of ResetUserPassword.
func (mr *MockUserServiceMockRecorder) ResetUserPassword(resetToken, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetUserPassword", reflect.TypeOf((*MockUserService)(nil).ResetUserPassword), resetToken, request)
}

// UpdateUserPassword mocks base method.
func (m *MockUserService) UpdateUserPassword(userID uuid.UUID, req payloads.UpdatePasswordRequest) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", userID, req)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockUserServiceMockRecorder) UpdateUserPassword(userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockUserService)(nil).UpdateUserPassword), userID, req)
}

// UserExistsByEmail mocks base method.
func (m *MockUserService) UserExistsByEmail(email string) (bool, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExistsByEmail", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// UserExistsByEmail indicates an expected call of UserExistsByEmail.
func (mr *MockUserServiceMockRecorder) UserExistsByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExistsByEmail", reflect.TypeOf((*MockUserService)(nil).UserExistsByEmail), email)
}

// VerifyUser mocks base method.
func (m *MockUserService) VerifyUser(token string) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUser", token)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// VerifyUser indicates an expected call of VerifyUser.
func (mr *MockUserServiceMockRecorder) VerifyUser(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUser", reflect.TypeOf((*MockUserService)(nil).VerifyUser), token)
}
