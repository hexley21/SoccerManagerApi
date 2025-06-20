// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hexley21/soccer-manager/internal/soccer-manager/service (interfaces: GlobeService)
//
// Generated by this command:
//
//	mockgen -destination=mock/mock_globe.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service GlobeService
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockGlobeService is a mock of GlobeService interface.
type MockGlobeService struct {
	ctrl     *gomock.Controller
	recorder *MockGlobeServiceMockRecorder
	isgomock struct{}
}

// MockGlobeServiceMockRecorder is the mock recorder for MockGlobeService.
type MockGlobeServiceMockRecorder struct {
	mock *MockGlobeService
}

// NewMockGlobeService creates a new mock instance.
func NewMockGlobeService(ctrl *gomock.Controller) *MockGlobeService {
	mock := &MockGlobeService{ctrl: ctrl}
	mock.recorder = &MockGlobeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobeService) EXPECT() *MockGlobeServiceMockRecorder {
	return m.recorder
}

// ListCountries mocks base method.
func (m *MockGlobeService) ListCountries(ctx context.Context) ([]domain.CountryCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCountries", ctx)
	ret0, _ := ret[0].([]domain.CountryCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCountries indicates an expected call of ListCountries.
func (mr *MockGlobeServiceMockRecorder) ListCountries(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCountries", reflect.TypeOf((*MockGlobeService)(nil).ListCountries), ctx)
}

// ListLocales mocks base method.
func (m *MockGlobeService) ListLocales(ctx context.Context) ([]domain.LocaleCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLocales", ctx)
	ret0, _ := ret[0].([]domain.LocaleCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLocales indicates an expected call of ListLocales.
func (mr *MockGlobeServiceMockRecorder) ListLocales(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLocales", reflect.TypeOf((*MockGlobeService)(nil).ListLocales), ctx)
}
