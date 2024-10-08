// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/genre_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/genre_service.go -destination=app/mocks/mock_services/genre_service.go -package=mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	errors "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
)

// MockGenreService is a mock of GenreService interface.
type MockGenreService struct {
	ctrl     *gomock.Controller
	recorder *MockGenreServiceMockRecorder
}

// MockGenreServiceMockRecorder is the mock recorder for MockGenreService.
type MockGenreServiceMockRecorder struct {
	mock *MockGenreService
}

// NewMockGenreService creates a new mock instance.
func NewMockGenreService(ctrl *gomock.Controller) *MockGenreService {
	mock := &MockGenreService{ctrl: ctrl}
	mock.recorder = &MockGenreServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenreService) EXPECT() *MockGenreServiceMockRecorder {
	return m.recorder
}

// CreateGenre mocks base method.
func (m *MockGenreService) CreateGenre(name string) (*models.Genre, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGenre", name)
	ret0, _ := ret[0].(*models.Genre)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateGenre indicates an expected call of CreateGenre.
func (mr *MockGenreServiceMockRecorder) CreateGenre(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGenre", reflect.TypeOf((*MockGenreService)(nil).CreateGenre), name)
}

// GetGenre mocks base method.
func (m *MockGenreService) GetGenre(id uuid.UUID) (*models.Genre, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenre", id)
	ret0, _ := ret[0].(*models.Genre)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetGenre indicates an expected call of GetGenre.
func (mr *MockGenreServiceMockRecorder) GetGenre(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenre", reflect.TypeOf((*MockGenreService)(nil).GetGenre), id)
}
