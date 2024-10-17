// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/state_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/state_service.go -destination=app/mocks/mock_services/state_service.go -package=mock_services
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

// MockStateService is a mock of StateService interface.
type MockStateService struct {
	ctrl     *gomock.Controller
	recorder *MockStateServiceMockRecorder
}

// MockStateServiceMockRecorder is the mock recorder for MockStateService.
type MockStateServiceMockRecorder struct {
	mock *MockStateService
}

// NewMockStateService creates a new mock instance.
func NewMockStateService(ctrl *gomock.Controller) *MockStateService {
	mock := &MockStateService{ctrl: ctrl}
	mock.recorder = &MockStateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateService) EXPECT() *MockStateServiceMockRecorder {
	return m.recorder
}

// CreateState mocks base method.
func (m *MockStateService) CreateState(countryID uuid.UUID, name string, code *string) (*models.State, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateState", countryID, name, code)
	ret0, _ := ret[0].(*models.State)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateState indicates an expected call of CreateState.
func (mr *MockStateServiceMockRecorder) CreateState(countryID, name, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateState", reflect.TypeOf((*MockStateService)(nil).CreateState), countryID, name, code)
}

// GetStatesByCountry mocks base method.
func (m *MockStateService) GetStatesByCountry(countryID uuid.UUID) ([]*models.State, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatesByCountry", countryID)
	ret0, _ := ret[0].([]*models.State)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetStatesByCountry indicates an expected call of GetStatesByCountry.
func (mr *MockStateServiceMockRecorder) GetStatesByCountry(countryID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatesByCountry", reflect.TypeOf((*MockStateService)(nil).GetStatesByCountry), countryID)
}
