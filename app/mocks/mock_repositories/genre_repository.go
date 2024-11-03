// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/genre_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/genre_repository.go -destination=app/mocks/mock_repositories/genre_repository.go -package=mock_repositories
//

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	filters "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockGenreRepository is a mock of GenreRepository interface.
type MockGenreRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGenreRepositoryMockRecorder
}

// MockGenreRepositoryMockRecorder is the mock recorder for MockGenreRepository.
type MockGenreRepositoryMockRecorder struct {
	mock *MockGenreRepository
}

// NewMockGenreRepository creates a new mock instance.
func NewMockGenreRepository(ctrl *gomock.Controller) *MockGenreRepository {
	mock := &MockGenreRepository{ctrl: ctrl}
	mock.recorder = &MockGenreRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenreRepository) EXPECT() *MockGenreRepositoryMockRecorder {
	return m.recorder
}

// CreateGenre mocks base method.
func (m *MockGenreRepository) CreateGenre(tx *gorm.DB, genre *models.Genre) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGenre", tx, genre)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGenre indicates an expected call of CreateGenre.
func (mr *MockGenreRepositoryMockRecorder) CreateGenre(tx, genre any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGenre", reflect.TypeOf((*MockGenreRepository)(nil).CreateGenre), tx, genre)
}

// GetGenre mocks base method.
func (m *MockGenreRepository) GetGenre(filter filters.GenreFilter) (*models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenre", filter)
	ret0, _ := ret[0].(*models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenre indicates an expected call of GetGenre.
func (mr *MockGenreRepositoryMockRecorder) GetGenre(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenre", reflect.TypeOf((*MockGenreRepository)(nil).GetGenre), filter)
}

// GetGenreIDs mocks base method.
func (m *MockGenreRepository) GetGenreIDs(filter filters.GenreFilter) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenreIDs", filter)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenreIDs indicates an expected call of GetGenreIDs.
func (mr *MockGenreRepositoryMockRecorder) GetGenreIDs(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenreIDs", reflect.TypeOf((*MockGenreRepository)(nil).GetGenreIDs), filter)
}

// GetGenres mocks base method.
func (m *MockGenreRepository) GetGenres(filter filters.GenreFilter) ([]*models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", filter)
	ret0, _ := ret[0].([]*models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockGenreRepositoryMockRecorder) GetGenres(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockGenreRepository)(nil).GetGenres), filter)
}

// UpdateGenre mocks base method.
func (m *MockGenreRepository) UpdateGenre(tx *gorm.DB, genre *models.Genre) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGenre", tx, genre)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGenre indicates an expected call of UpdateGenre.
func (mr *MockGenreRepositoryMockRecorder) UpdateGenre(tx, genre any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGenre", reflect.TypeOf((*MockGenreRepository)(nil).UpdateGenre), tx, genre)
}
