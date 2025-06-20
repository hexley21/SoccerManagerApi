// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hexley21/soccer-manager/internal/soccer-manager/repository (interfaces: TeamRepository)
//
// Generated by this command:
//
//	mockgen -destination=mock/mock_team.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository TeamRepository
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	repository "github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockTeamRepository is a mock of TeamRepository interface.
type MockTeamRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTeamRepositoryMockRecorder
	isgomock struct{}
}

// MockTeamRepositoryMockRecorder is the mock recorder for MockTeamRepository.
type MockTeamRepositoryMockRecorder struct {
	mock *MockTeamRepository
}

// NewMockTeamRepository creates a new mock instance.
func NewMockTeamRepository(ctrl *gomock.Controller) *MockTeamRepository {
	mock := &MockTeamRepository{ctrl: ctrl}
	mock.recorder = &MockTeamRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTeamRepository) EXPECT() *MockTeamRepositoryMockRecorder {
	return m.recorder
}

// DeleteTeamByID mocks base method.
func (m *MockTeamRepository) DeleteTeamByID(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTeamByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTeamByID indicates an expected call of DeleteTeamByID.
func (mr *MockTeamRepositoryMockRecorder) DeleteTeamByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTeamByID", reflect.TypeOf((*MockTeamRepository)(nil).DeleteTeamByID), ctx, id)
}

// GetTeamByID mocks base method.
func (m *MockTeamRepository) GetTeamByID(ctx context.Context, id int64) (repository.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamByID", ctx, id)
	ret0, _ := ret[0].(repository.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamByID indicates an expected call of GetTeamByID.
func (mr *MockTeamRepositoryMockRecorder) GetTeamByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamByID", reflect.TypeOf((*MockTeamRepository)(nil).GetTeamByID), ctx, id)
}

// GetTeamByUserID mocks base method.
func (m *MockTeamRepository) GetTeamByUserID(ctx context.Context, userID int64) (repository.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeamByUserID", ctx, userID)
	ret0, _ := ret[0].(repository.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeamByUserID indicates an expected call of GetTeamByUserID.
func (mr *MockTeamRepositoryMockRecorder) GetTeamByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeamByUserID", reflect.TypeOf((*MockTeamRepository)(nil).GetTeamByUserID), ctx, userID)
}

// InsertTeam mocks base method.
func (m *MockTeamRepository) InsertTeam(ctx context.Context, arg repository.InsertTeamParams) (repository.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTeam", ctx, arg)
	ret0, _ := ret[0].(repository.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTeam indicates an expected call of InsertTeam.
func (mr *MockTeamRepositoryMockRecorder) InsertTeam(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTeam", reflect.TypeOf((*MockTeamRepository)(nil).InsertTeam), ctx, arg)
}

// ListTeamsCursor mocks base method.
func (m *MockTeamRepository) ListTeamsCursor(ctx context.Context, arg repository.ListTeamsCursorParams) ([]repository.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTeamsCursor", ctx, arg)
	ret0, _ := ret[0].([]repository.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTeamsCursor indicates an expected call of ListTeamsCursor.
func (mr *MockTeamRepositoryMockRecorder) ListTeamsCursor(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTeamsCursor", reflect.TypeOf((*MockTeamRepository)(nil).ListTeamsCursor), ctx, arg)
}

// UpdateTeamDataByUserID mocks base method.
func (m *MockTeamRepository) UpdateTeamDataByUserID(ctx context.Context, arg repository.UpdateTeamDataByUserIDParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTeamDataByUserID", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTeamDataByUserID indicates an expected call of UpdateTeamDataByUserID.
func (mr *MockTeamRepositoryMockRecorder) UpdateTeamDataByUserID(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTeamDataByUserID", reflect.TypeOf((*MockTeamRepository)(nil).UpdateTeamDataByUserID), ctx, arg)
}
