// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hexley21/soccer-manager/internal/soccer-manager/service (interfaces: TransferRecordService)
//
// Generated by this command:
//
//	mockgen -destination=mock/mock_transfer_record.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service TransferRecordService
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockTransferRecordService is a mock of TransferRecordService interface.
type MockTransferRecordService struct {
	ctrl     *gomock.Controller
	recorder *MockTransferRecordServiceMockRecorder
	isgomock struct{}
}

// MockTransferRecordServiceMockRecorder is the mock recorder for MockTransferRecordService.
type MockTransferRecordServiceMockRecorder struct {
	mock *MockTransferRecordService
}

// NewMockTransferRecordService creates a new mock instance.
func NewMockTransferRecordService(ctrl *gomock.Controller) *MockTransferRecordService {
	mock := &MockTransferRecordService{ctrl: ctrl}
	mock.recorder = &MockTransferRecordServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferRecordService) EXPECT() *MockTransferRecordServiceMockRecorder {
	return m.recorder
}

// GetTransferRecordByID mocks base method.
func (m *MockTransferRecordService) GetTransferRecordByID(ctx context.Context, id int64) (domain.TransferRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransferRecordByID", ctx, id)
	ret0, _ := ret[0].(domain.TransferRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransferRecordByID indicates an expected call of GetTransferRecordByID.
func (mr *MockTransferRecordServiceMockRecorder) GetTransferRecordByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransferRecordByID", reflect.TypeOf((*MockTransferRecordService)(nil).GetTransferRecordByID), ctx, id)
}

// ListTransferRecords mocks base method.
func (m *MockTransferRecordService) ListTransferRecords(ctx context.Context, id int64, limit int32) ([]domain.TransferRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTransferRecords", ctx, id, limit)
	ret0, _ := ret[0].([]domain.TransferRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTransferRecords indicates an expected call of ListTransferRecords.
func (mr *MockTransferRecordServiceMockRecorder) ListTransferRecords(ctx, id, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTransferRecords", reflect.TypeOf((*MockTransferRecordService)(nil).ListTransferRecords), ctx, id, limit)
}
