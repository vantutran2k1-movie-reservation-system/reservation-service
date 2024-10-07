// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/movie_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/movie_service.go -destination=app/mocks/mock_services/movie_service.go -package=mock_services
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

// MockMovieService is a mock of MovieService interface.
type MockMovieService struct {
	ctrl     *gomock.Controller
	recorder *MockMovieServiceMockRecorder
}

// MockMovieServiceMockRecorder is the mock recorder for MockMovieService.
type MockMovieServiceMockRecorder struct {
	mock *MockMovieService
}

// NewMockMovieService creates a new mock instance.
func NewMockMovieService(ctrl *gomock.Controller) *MockMovieService {
	mock := &MockMovieService{ctrl: ctrl}
	mock.recorder = &MockMovieServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMovieService) EXPECT() *MockMovieServiceMockRecorder {
	return m.recorder
}

// CreateMovie mocks base method.
func (m *MockMovieService) CreateMovie(title string, description *string, releaseDate string, duration int, language *string, rating *float64, createdBy uuid.UUID) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovie", title, description, releaseDate, duration, language, rating, createdBy)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateMovie indicates an expected call of CreateMovie.
func (mr *MockMovieServiceMockRecorder) CreateMovie(title, description, releaseDate, duration, language, rating, createdBy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovie", reflect.TypeOf((*MockMovieService)(nil).CreateMovie), title, description, releaseDate, duration, language, rating, createdBy)
}

// GetMovie mocks base method.
func (m *MockMovieService) GetMovie(id uuid.UUID) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovie", id)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetMovie indicates an expected call of GetMovie.
func (mr *MockMovieServiceMockRecorder) GetMovie(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovie", reflect.TypeOf((*MockMovieService)(nil).GetMovie), id)
}

// GetMovies mocks base method.
func (m *MockMovieService) GetMovies(limit, offset int) ([]*models.Movie, *models.ResponseMeta, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", limit, offset)
	ret0, _ := ret[0].([]*models.Movie)
	ret1, _ := ret[1].(*models.ResponseMeta)
	ret2, _ := ret[2].(*errors.ApiError)
	return ret0, ret1, ret2
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockMovieServiceMockRecorder) GetMovies(limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockMovieService)(nil).GetMovies), limit, offset)
}

// UpdateMovie mocks base method.
func (m *MockMovieService) UpdateMovie(id, updatedBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMovie", id, updatedBy, title, description, releaseDate, duration, language, rating)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// UpdateMovie indicates an expected call of UpdateMovie.
func (mr *MockMovieServiceMockRecorder) UpdateMovie(id, updatedBy, title, description, releaseDate, duration, language, rating any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMovie", reflect.TypeOf((*MockMovieService)(nil).UpdateMovie), id, updatedBy, title, description, releaseDate, duration, language, rating)
}