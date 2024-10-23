// Code generated by MockGen. DO NOT EDIT.
// Source: app/services/theater_service.go
//
// Generated by this command:
//
//	mockgen -source=app/services/theater_service.go -destination=app/mocks/mock_services/theater_service.go -package=mock_services
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

// MockTheaterService is a mock of TheaterService interface.
type MockTheaterService struct {
	ctrl     *gomock.Controller
	recorder *MockTheaterServiceMockRecorder
}

// MockTheaterServiceMockRecorder is the mock recorder for MockTheaterService.
type MockTheaterServiceMockRecorder struct {
	mock *MockTheaterService
}

// NewMockTheaterService creates a new mock instance.
func NewMockTheaterService(ctrl *gomock.Controller) *MockTheaterService {
	mock := &MockTheaterService{ctrl: ctrl}
	mock.recorder = &MockTheaterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTheaterService) EXPECT() *MockTheaterServiceMockRecorder {
	return m.recorder
}

// CreateTheater mocks base method.
func (m *MockTheaterService) CreateTheater(req payloads.CreateTheaterRequest) (*models.Theater, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTheater", req)
	ret0, _ := ret[0].(*models.Theater)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateTheater indicates an expected call of CreateTheater.
func (mr *MockTheaterServiceMockRecorder) CreateTheater(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTheater", reflect.TypeOf((*MockTheaterService)(nil).CreateTheater), req)
}

// CreateTheaterLocation mocks base method.
func (m *MockTheaterService) CreateTheaterLocation(theaterID uuid.UUID, req payloads.CreateTheaterLocationRequest) (*models.TheaterLocation, *errors.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTheaterLocation", theaterID, req)
	ret0, _ := ret[0].(*models.TheaterLocation)
	ret1, _ := ret[1].(*errors.ApiError)
	return ret0, ret1
}

// CreateTheaterLocation indicates an expected call of CreateTheaterLocation.
func (mr *MockTheaterServiceMockRecorder) CreateTheaterLocation(theaterID, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTheaterLocation", reflect.TypeOf((*MockTheaterService)(nil).CreateTheaterLocation), theaterID, req)
}
