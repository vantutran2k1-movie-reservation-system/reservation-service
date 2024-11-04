// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/user_profile_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/user_profile_repository.go -destination=app/mocks/mock_repositories/user_profile_repository.go -package=mock_repositories
//

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	reflect "reflect"

	filters "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockUserProfileRepository is a mock of UserProfileRepository interface.
type MockUserProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserProfileRepositoryMockRecorder
}

// MockUserProfileRepositoryMockRecorder is the mock recorder for MockUserProfileRepository.
type MockUserProfileRepositoryMockRecorder struct {
	mock *MockUserProfileRepository
}

// NewMockUserProfileRepository creates a new mock instance.
func NewMockUserProfileRepository(ctrl *gomock.Controller) *MockUserProfileRepository {
	mock := &MockUserProfileRepository{ctrl: ctrl}
	mock.recorder = &MockUserProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserProfileRepository) EXPECT() *MockUserProfileRepositoryMockRecorder {
	return m.recorder
}

// CreateUserProfile mocks base method.
func (m *MockUserProfileRepository) CreateUserProfile(tx *gorm.DB, profile *models.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserProfile", tx, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserProfile indicates an expected call of CreateUserProfile.
func (mr *MockUserProfileRepositoryMockRecorder) CreateUserProfile(tx, profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserProfile", reflect.TypeOf((*MockUserProfileRepository)(nil).CreateUserProfile), tx, profile)
}

// GetProfile mocks base method.
func (m *MockUserProfileRepository) GetProfile(filter filters.UserProfileFilter) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", filter)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockUserProfileRepositoryMockRecorder) GetProfile(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockUserProfileRepository)(nil).GetProfile), filter)
}

// UpdateProfilePicture mocks base method.
func (m *MockUserProfileRepository) UpdateProfilePicture(tx *gorm.DB, profile *models.UserProfile, url *string) (*models.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfilePicture", tx, profile, url)
	ret0, _ := ret[0].(*models.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfilePicture indicates an expected call of UpdateProfilePicture.
func (mr *MockUserProfileRepositoryMockRecorder) UpdateProfilePicture(tx, profile, url any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfilePicture", reflect.TypeOf((*MockUserProfileRepository)(nil).UpdateProfilePicture), tx, profile, url)
}

// UpdateUserProfile mocks base method.
func (m *MockUserProfileRepository) UpdateUserProfile(tx *gorm.DB, profile *models.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", tx, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockUserProfileRepositoryMockRecorder) UpdateUserProfile(tx, profile any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockUserProfileRepository)(nil).UpdateUserProfile), tx, profile)
}
