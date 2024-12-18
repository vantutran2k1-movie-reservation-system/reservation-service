// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/state_repository.go
//
// Generated by this command:
//
//	mockgen -source=app/repositories/state_repository.go -destination=app/mocks/mock_repositories/state_repository.go -package=mock_repositories
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

// MockStateRepository is a mock of StateRepository interface.
type MockStateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStateRepositoryMockRecorder
}

// MockStateRepositoryMockRecorder is the mock recorder for MockStateRepository.
type MockStateRepositoryMockRecorder struct {
	mock *MockStateRepository
}

// NewMockStateRepository creates a new mock instance.
func NewMockStateRepository(ctrl *gomock.Controller) *MockStateRepository {
	mock := &MockStateRepository{ctrl: ctrl}
	mock.recorder = &MockStateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateRepository) EXPECT() *MockStateRepositoryMockRecorder {
	return m.recorder
}

// CreateState mocks base method.
func (m *MockStateRepository) CreateState(tx *gorm.DB, state *models.State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateState", tx, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateState indicates an expected call of CreateState.
func (mr *MockStateRepositoryMockRecorder) CreateState(tx, state any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateState", reflect.TypeOf((*MockStateRepository)(nil).CreateState), tx, state)
}

// GetState mocks base method.
func (m *MockStateRepository) GetState(filter filters.StateFilter) (*models.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", filter)
	ret0, _ := ret[0].(*models.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockStateRepositoryMockRecorder) GetState(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStateRepository)(nil).GetState), filter)
}

// GetStates mocks base method.
func (m *MockStateRepository) GetStates(filter filters.StateFilter) ([]*models.State, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStates", filter)
	ret0, _ := ret[0].([]*models.State)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStates indicates an expected call of GetStates.
func (mr *MockStateRepositoryMockRecorder) GetStates(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStates", reflect.TypeOf((*MockStateRepository)(nil).GetStates), filter)
}
