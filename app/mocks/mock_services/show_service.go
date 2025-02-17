// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/show_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/show_service.go -destination=app/mocks/mock_services/show_service.go -package=mock_services
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	constants "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	errors "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	models "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	payloads "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	gomock "go.uber.org/mock/gomock"
)

// MockShowService is a mock of ShowService interface.
type MockShowService struct {
	ctrl     *gomock.Controller
	recorder *MockShowServiceMockRecorder
}

// MockShowServiceMockRecorder is the mock recorder for MockShowService.
type MockShowServiceMockRecorder struct {
	mock *MockShowService
}

// NewMockShowService creates a new mock instance.
func NewMockShowService(ctrl *gomock.Controller) *MockShowService {
	mock := &MockShowService{ctrl: ctrl}
	mock.recorder = &MockShowServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShowService) EXPECT() *MockShowServiceMockRecorder {
	return m.recorder
}

// CreateShow mocks base method.
func (m *MockShowService) CreateShow(req payloads.CreateShowRequest) (*models.Show, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShow", req)
	ret0, _ := ret[0].(*models.Show)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateShow indicates an expected call of CreateShow.
func (mr *MockShowServiceMockRecorder) CreateShow(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShow", reflect.TypeOf((*MockShowService)(nil).CreateShow), req)
}

// GetShow mocks base method.
func (m *MockShowService) GetShow(id uuid.UUID, userEmail *string) (*models.Show, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShow", id, userEmail)
	ret0, _ := ret[0].(*models.Show)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetShow indicates an expected call of GetShow.
func (mr *MockShowServiceMockRecorder) GetShow(id, userEmail any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShow", reflect.TypeOf((*MockShowService)(nil).GetShow), id, userEmail)
}

// GetShows mocks base method.
func (m *MockShowService) GetShows(status constants.ShowStatus, limit, offset int) ([]*models.Show, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShows", status, limit, offset)
	ret0, _ := ret[0].([]*models.Show)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// GetShows indicates an expected call of GetShows.
func (mr *MockShowServiceMockRecorder) GetShows(status, limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShows", reflect.TypeOf((*MockShowService)(nil).GetShows), status, limit, offset)
}

// ScheduleUpdateShowStatus mocks base method.
func (m *MockShowService) ScheduleUpdateShowStatus() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduleUpdateShowStatus")
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleUpdateShowStatus indicates an expected call of ScheduleUpdateShowStatus.
func (mr *MockShowServiceMockRecorder) ScheduleUpdateShowStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleUpdateShowStatus", reflect.TypeOf((*MockShowService)(nil).ScheduleUpdateShowStatus))
}
