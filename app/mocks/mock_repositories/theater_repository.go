// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/theater_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/theater_repository.go -destination=app/mocks/mock_repositories/theater_repository.go -package=mock_repositories
//

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockTheaterRepository is a mock of TheaterRepository interface.
type MockTheaterRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTheaterRepositoryMockRecorder
}

// MockTheaterRepositoryMockRecorder is the mock recorder for MockTheaterRepository.
type MockTheaterRepositoryMockRecorder struct {
	mock *MockTheaterRepository
}

// NewMockTheaterRepository creates a new mock instance.
func NewMockTheaterRepository(ctrl *gomock.Controller) *MockTheaterRepository {
	mock := &MockTheaterRepository{ctrl: ctrl}
	mock.recorder = &MockTheaterRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTheaterRepository) EXPECT() *MockTheaterRepositoryMockRecorder {
	return m.recorder
}

// CreateTheater mocks base method.
func (m *MockTheaterRepository) CreateTheater(tx *gorm.DB, theater *models.Theater) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTheater", tx, theater)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTheater indicates an expected call of CreateTheater.
func (mr *MockTheaterRepositoryMockRecorder) CreateTheater(tx, theater any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTheater", reflect.TypeOf((*MockTheaterRepository)(nil).CreateTheater), tx, theater)
}

// GetTheater mocks base method.
func (m *MockTheaterRepository) GetTheater(id uuid.UUID) (*models.Theater, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTheater", id)
	ret0, _ := ret[0].(*models.Theater)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTheater indicates an expected call of GetTheater.
func (mr *MockTheaterRepositoryMockRecorder) GetTheater(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTheater", reflect.TypeOf((*MockTheaterRepository)(nil).GetTheater), id)
}

// GetTheaterByName mocks base method.
func (m *MockTheaterRepository) GetTheaterByName(name string) (*models.Theater, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTheaterByName", name)
	ret0, _ := ret[0].(*models.Theater)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTheaterByName indicates an expected call of GetTheaterByName.
func (mr *MockTheaterRepositoryMockRecorder) GetTheaterByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTheaterByName", reflect.TypeOf((*MockTheaterRepository)(nil).GetTheaterByName), name)
}
