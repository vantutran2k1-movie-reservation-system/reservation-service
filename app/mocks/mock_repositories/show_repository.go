// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/show_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/show_repository.go -destination=app/mocks/mock_repositories/show_repository.go -package=mock_repositories
//

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	constants "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	filters "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockShowRepository is a mock of ShowRepository interface.
type MockShowRepository struct {
	ctrl     *gomock.Controller
	recorder *MockShowRepositoryMockRecorder
}

// MockShowRepositoryMockRecorder is the mock recorder for MockShowRepository.
type MockShowRepositoryMockRecorder struct {
	mock *MockShowRepository
}

// NewMockShowRepository creates a new mock instance.
func NewMockShowRepository(ctrl *gomock.Controller) *MockShowRepository {
	mock := &MockShowRepository{ctrl: ctrl}
	mock.recorder = &MockShowRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShowRepository) EXPECT() *MockShowRepositoryMockRecorder {
	return m.recorder
}

// CreateShow mocks base method.
func (m *MockShowRepository) CreateShow(tx *gorm.DB, show *models.Show) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShow", tx, show)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateShow indicates an expected call of CreateShow.
func (mr *MockShowRepositoryMockRecorder) CreateShow(tx, show any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShow", reflect.TypeOf((*MockShowRepository)(nil).CreateShow), tx, show)
}

// GetShow mocks base method.
func (m *MockShowRepository) GetShow(filter filters.ShowFilter) (*models.Show, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShow", filter)
	ret0, _ := ret[0].(*models.Show)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShow indicates an expected call of GetShow.
func (mr *MockShowRepositoryMockRecorder) GetShow(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShow", reflect.TypeOf((*MockShowRepository)(nil).GetShow), filter)
}

// GetShows mocks base method.
func (m *MockShowRepository) GetShows(filter filters.ShowFilter) ([]*models.Show, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShows", filter)
	ret0, _ := ret[0].([]*models.Show)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShows indicates an expected call of GetShows.
func (mr *MockShowRepositoryMockRecorder) GetShows(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShows", reflect.TypeOf((*MockShowRepository)(nil).GetShows), filter)
}

// IsShowInValidTimeRange mocks base method.
func (m *MockShowRepository) IsShowInValidTimeRange(theaterId uuid.UUID, startTime, endTime time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsShowInValidTimeRange", theaterId, startTime, endTime)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsShowInValidTimeRange indicates an expected call of IsShowInValidTimeRange.
func (mr *MockShowRepositoryMockRecorder) IsShowInValidTimeRange(theaterId, startTime, endTime any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsShowInValidTimeRange", reflect.TypeOf((*MockShowRepository)(nil).IsShowInValidTimeRange), theaterId, startTime, endTime)
}

// ScheduleActivateShows mocks base method.
func (m *MockShowRepository) ScheduleActivateShows(tx *gorm.DB, beforeStart time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleActivateShows", tx, beforeStart)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleActivateShows indicates an expected call of ScheduleActivateShows.
func (mr *MockShowRepositoryMockRecorder) ScheduleActivateShows(tx, beforeStart any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleActivateShows", reflect.TypeOf((*MockShowRepository)(nil).ScheduleActivateShows), tx, beforeStart)
}

// ScheduleCompleteShows mocks base method.
func (m *MockShowRepository) ScheduleCompleteShows(tx *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleCompleteShows", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleCompleteShows indicates an expected call of ScheduleCompleteShows.
func (mr *MockShowRepositoryMockRecorder) ScheduleCompleteShows(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleCompleteShows", reflect.TypeOf((*MockShowRepository)(nil).ScheduleCompleteShows), tx)
}

// UpdateShowStatus mocks base method.
func (m *MockShowRepository) UpdateShowStatus(tx *gorm.DB, showId uuid.UUID, status constants.ShowStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShowStatus", tx, showId, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShowStatus indicates an expected call of UpdateShowStatus.
func (mr *MockShowRepositoryMockRecorder) UpdateShowStatus(tx, showId, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShowStatus", reflect.TypeOf((*MockShowRepository)(nil).UpdateShowStatus), tx, showId, status)
}
