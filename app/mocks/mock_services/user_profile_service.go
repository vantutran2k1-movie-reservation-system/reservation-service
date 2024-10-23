// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/user_profile_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/user_profile_service.go -destination=app/mocks/mock_services/user_profile_service.go -package=mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	multipart "mime/multipart"
	reflect "reflect"

	uuid "github.com/google/uuid"
	errors "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	payloads "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	gomock "go.uber.org/mock/gomock"
)

// MockUserProfileService is a mock of UserProfileService interface.
type MockUserProfileService struct {
	ctrl     *gomock.Controller
	recorder *MockUserProfileServiceMockRecorder
}

// MockUserProfileServiceMockRecorder is the mock recorder for MockUserProfileService.
type MockUserProfileServiceMockRecorder struct {
	mock *MockUserProfileService
}

// NewMockUserProfileService creates a new mock instance.
func NewMockUserProfileService(ctrl *gomock.Controller) *MockUserProfileService {
	mock := &MockUserProfileService{ctrl: ctrl}
	mock.recorder = &MockUserProfileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserProfileService) EXPECT() *MockUserProfileServiceMockRecorder {
	return m.recorder
}

// CreateUserProfile mocks base method.
func (m *MockUserProfileService) CreateUserProfile(userID uuid.UUID, req payloads.CreateUserProfileRequest) (*models.UserProfile, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserProfile", userID, req)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateUserProfile indicates an expected call of CreateUserProfile.
func (mr *MockUserProfileServiceMockRecorder) CreateUserProfile(userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserProfile", reflect.TypeOf((*MockUserProfileService)(nil).CreateUserProfile), userID, req)
}

// DeleteProfilePicture mocks base method.
func (m *MockUserProfileService) DeleteProfilePicture(userID uuid.UUID) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProfilePicture", userID)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// DeleteProfilePicture indicates an expected call of DeleteProfilePicture.
func (mr *MockUserProfileServiceMockRecorder) DeleteProfilePicture(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProfilePicture", reflect.TypeOf((*MockUserProfileService)(nil).DeleteProfilePicture), userID)
}

// GetProfileByUserID mocks base method.
func (m *MockUserProfileService) GetProfileByUserID(userID uuid.UUID) (*models.UserProfile, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByUserID", userID)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetProfileByUserID indicates an expected call of GetProfileByUserID.
func (mr *MockUserProfileServiceMockRecorder) GetProfileByUserID(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByUserID", reflect.TypeOf((*MockUserProfileService)(nil).GetProfileByUserID), userID)
}

// UpdateProfilePicture mocks base method.
func (m *MockUserProfileService) UpdateProfilePicture(userID uuid.UUID, file *multipart.FileHeader) (*models.UserProfile, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfilePicture", userID, file)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// UpdateProfilePicture indicates an expected call of UpdateProfilePicture.
func (mr *MockUserProfileServiceMockRecorder) UpdateProfilePicture(userID, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfilePicture", reflect.TypeOf((*MockUserProfileService)(nil).UpdateProfilePicture), userID, file)
}

// UpdateUserProfile mocks base method.
func (m *MockUserProfileService) UpdateUserProfile(userID uuid.UUID, req payloads.UpdateUserProfileRequest) (*models.UserProfile, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", userID, req)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockUserProfileServiceMockRecorder) UpdateUserProfile(userID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockUserProfileService)(nil).UpdateUserProfile), userID, req)
}
