// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hexley21/soccer-manager/internal/soccer-manager/repository (interfaces: UserRepository)
//
// Generated by this command:
//
//	mockgen -destination=mock/mock_user.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository UserRepository
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	repository "github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
	isgomock struct{}
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CheckUserExists mocks base method.
func (m *MockUserRepository) CheckUserExists(ctx context.Context, id int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserExists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserExists indicates an expected call of CheckUserExists.
func (mr *MockUserRepositoryMockRecorder) CheckUserExists(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserExists", reflect.TypeOf((*MockUserRepository)(nil).CheckUserExists), ctx, id)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.CreateUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(repository.CreateUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), ctx, arg)
}

// DeleteUser mocks base method.
func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserRepositoryMockRecorder) DeleteUser(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepository)(nil).DeleteUser), ctx, id)
}

// GetAuth mocks base method.
func (m *MockUserRepository) GetAuth(ctx context.Context, username string) (repository.GetAuthRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuth", ctx, username)
	ret0, _ := ret[0].(repository.GetAuthRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuth indicates an expected call of GetAuth.
func (mr *MockUserRepositoryMockRecorder) GetAuth(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuth", reflect.TypeOf((*MockUserRepository)(nil).GetAuth), ctx, username)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (repository.GetUserByIDRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, id)
	ret0, _ := ret[0].(repository.GetUserByIDRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), ctx, id)
}

// GetUserByUsername mocks base method.
func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (repository.GetUserByUsernameRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(repository.GetUserByUsernameRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsername(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsername), ctx, username)
}

// GetUserHashByID mocks base method.
func (m *MockUserRepository) GetUserHashByID(ctx context.Context, id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserHashByID", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserHashByID indicates an expected call of GetUserHashByID.
func (mr *MockUserRepositoryMockRecorder) GetUserHashByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserHashByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserHashByID), ctx, id)
}

// ListUsersCursor mocks base method.
func (m *MockUserRepository) ListUsersCursor(ctx context.Context, arg repository.ListUsersCursorParams) ([]repository.ListUsersCursorRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsersCursor", ctx, arg)
	ret0, _ := ret[0].([]repository.ListUsersCursorRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsersCursor indicates an expected call of ListUsersCursor.
func (mr *MockUserRepositoryMockRecorder) ListUsersCursor(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsersCursor", reflect.TypeOf((*MockUserRepository)(nil).ListUsersCursor), ctx, arg)
}

// UpdateUserHash mocks base method.
func (m *MockUserRepository) UpdateUserHash(ctx context.Context, arg repository.UpdateUserHashParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserHash", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserHash indicates an expected call of UpdateUserHash.
func (mr *MockUserRepositoryMockRecorder) UpdateUserHash(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserHash", reflect.TypeOf((*MockUserRepository)(nil).UpdateUserHash), ctx, arg)
}
