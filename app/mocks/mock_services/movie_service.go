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
	payloads "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
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

// AssignGenres mocks base method.
func (m *MockMovieService) AssignGenres(id uuid.UUID, genreIDs []uuid.UUID) *errors.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignGenres", id, genreIDs)
	ret0, _ := ret[0].(*errors.ApiError)
	return ret0
}

// AssignGenres indicates an expected call of AssignGenres.
func (mr *MockMovieServiceMockRecorder) AssignGenres(id, genreIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignGenres", reflect.TypeOf((*MockMovieService)(nil).AssignGenres), id, genreIDs)
}

// CreateMovie mocks base method.
func (m *MockMovieService) CreateMovie(req payloads.CreateMovieRequest, createdBy uuid.UUID) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovie", req, createdBy)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateMovie indicates an expected call of CreateMovie.
func (mr *MockMovieServiceMockRecorder) CreateMovie(req, createdBy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovie", reflect.TypeOf((*MockMovieService)(nil).CreateMovie), req, createdBy)
}

// GetMovie mocks base method.
func (m *MockMovieService) GetMovie(id uuid.UUID, includeGenres bool) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovie", id, includeGenres)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetMovie indicates an expected call of GetMovie.
func (mr *MockMovieServiceMockRecorder) GetMovie(id, includeGenres any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovie", reflect.TypeOf((*MockMovieService)(nil).GetMovie), id, includeGenres)
}

// GetMovies mocks base method.
func (m *MockMovieService) GetMovies(limit, offset int, includeGenres bool) ([]*models.Movie, *models.ResponseMeta, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovies", limit, offset, includeGenres)
	ret0, _ := ret[0].([]*models.Movie)
	ret1, _ := ret[1].(*models.ResponseMeta)
	ret2, _ := ret[2].(*errors.ApiError)
	return ret0, ret1, ret2
}

// GetMovies indicates an expected call of GetMovies.
func (mr *MockMovieServiceMockRecorder) GetMovies(limit, offset, includeGenres any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovies", reflect.TypeOf((*MockMovieService)(nil).GetMovies), limit, offset, includeGenres)
}

// UpdateMovie mocks base method.
func (m *MockMovieService) UpdateMovie(id, updatedBy uuid.UUID, req payloads.UpdateMovieRequest) (*models.Movie, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMovie", id, updatedBy, req)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// UpdateMovie indicates an expected call of UpdateMovie.
func (mr *MockMovieServiceMockRecorder) UpdateMovie(id, updatedBy, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMovie", reflect.TypeOf((*MockMovieService)(nil).UpdateMovie), id, updatedBy, req)
}
