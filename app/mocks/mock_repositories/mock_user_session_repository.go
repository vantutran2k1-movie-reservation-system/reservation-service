// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/user_session_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/user_session_repository.go -destination=app/mocks/mock_repositories/mock_user_session_repository.go -package=mock_repositories
//

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
)

// MockUserSessionRepository is a mock of UserSessionRepository interface.
type MockUserSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserSessionRepositoryMockRecorder
}

// MockUserSessionRepositoryMockRecorder is the mock recorder for MockUserSessionRepository.
type MockUserSessionRepositoryMockRecorder struct {
	mock *MockUserSessionRepository
}

// NewMockUserSessionRepository creates a new mock instance.
func NewMockUserSessionRepository(ctrl *gomock.Controller) *MockUserSessionRepository {
	mock := &MockUserSessionRepository{ctrl: ctrl}
	mock.recorder = &MockUserSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserSessionRepository) EXPECT() *MockUserSessionRepositoryMockRecorder {
	return m.recorder
}

// CreateUserSession mocks base method.
func (m *MockUserSessionRepository) CreateUserSession(sessionID string, expiration time.Duration, session *models.UserSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserSession", sessionID, expiration, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserSession indicates an expected call of CreateUserSession.
func (mr *MockUserSessionRepositoryMockRecorder) CreateUserSession(sessionID, expiration, session any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserSession", reflect.TypeOf((*MockUserSessionRepository)(nil).CreateUserSession), sessionID, expiration, session)
}

// DeleteUserSession mocks base method.
func (m *MockUserSessionRepository) DeleteUserSession(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSession", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSession indicates an expected call of DeleteUserSession.
func (mr *MockUserSessionRepositoryMockRecorder) DeleteUserSession(sessionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSession", reflect.TypeOf((*MockUserSessionRepository)(nil).DeleteUserSession), sessionID)
}

// DeleteUserSessions mocks base method.
func (m *MockUserSessionRepository) DeleteUserSessions(userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSessions", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSessions indicates an expected call of DeleteUserSessions.
func (mr *MockUserSessionRepositoryMockRecorder) DeleteUserSessions(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSessions", reflect.TypeOf((*MockUserSessionRepository)(nil).DeleteUserSessions), userID)
}

// GetUserSession mocks base method.
func (m *MockUserSessionRepository) GetUserSession(sessionID string) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSession", sessionID)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSession indicates an expected call of GetUserSession.
func (mr *MockUserSessionRepositoryMockRecorder) GetUserSession(sessionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSession", reflect.TypeOf((*MockUserSessionRepository)(nil).GetUserSession), sessionID)
}

// GetUserSessionID mocks base method.
func (m *MockUserSessionRepository) GetUserSessionID(tokenValue string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSessionID", tokenValue)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetUserSessionID indicates an expected call of GetUserSessionID.
func (mr *MockUserSessionRepositoryMockRecorder) GetUserSessionID(tokenValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSessionID", reflect.TypeOf((*MockUserSessionRepository)(nil).GetUserSessionID), tokenValue)
}