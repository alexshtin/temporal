// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by MockGen. DO NOT EDIT.
// Source: visibilityInterfaces.go

// Package persistence is a generated GoMock package.
package persistence

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockVisibilityManager is a mock of VisibilityManager interface.
type MockVisibilityManager struct {
	ctrl     *gomock.Controller
	recorder *MockVisibilityManagerMockRecorder
}

// MockVisibilityManagerMockRecorder is the mock recorder for MockVisibilityManager.
type MockVisibilityManagerMockRecorder struct {
	mock *MockVisibilityManager
}

// NewMockVisibilityManager creates a new mock instance.
func NewMockVisibilityManager(ctrl *gomock.Controller) *MockVisibilityManager {
	mock := &MockVisibilityManager{ctrl: ctrl}
	mock.recorder = &MockVisibilityManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVisibilityManager) EXPECT() *MockVisibilityManagerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockVisibilityManager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockVisibilityManagerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockVisibilityManager)(nil).Close))
}

// CountWorkflowExecutions mocks base method.
func (m *MockVisibilityManager) CountWorkflowExecutions(request *CountWorkflowExecutionsRequest) (*CountWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountWorkflowExecutions", request)
	ret0, _ := ret[0].(*CountWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountWorkflowExecutions indicates an expected call of CountWorkflowExecutions.
func (mr *MockVisibilityManagerMockRecorder) CountWorkflowExecutions(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountWorkflowExecutions", reflect.TypeOf((*MockVisibilityManager)(nil).CountWorkflowExecutions), request)
}

// DeleteWorkflowExecution mocks base method.
func (m *MockVisibilityManager) DeleteWorkflowExecution(request *VisibilityDeleteWorkflowExecutionRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWorkflowExecution", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWorkflowExecution indicates an expected call of DeleteWorkflowExecution.
func (mr *MockVisibilityManagerMockRecorder) DeleteWorkflowExecution(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWorkflowExecution", reflect.TypeOf((*MockVisibilityManager)(nil).DeleteWorkflowExecution), request)
}

// GetClosedWorkflowExecution mocks base method.
func (m *MockVisibilityManager) GetClosedWorkflowExecution(request *GetClosedWorkflowExecutionRequest) (*GetClosedWorkflowExecutionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClosedWorkflowExecution", request)
	ret0, _ := ret[0].(*GetClosedWorkflowExecutionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClosedWorkflowExecution indicates an expected call of GetClosedWorkflowExecution.
func (mr *MockVisibilityManagerMockRecorder) GetClosedWorkflowExecution(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClosedWorkflowExecution", reflect.TypeOf((*MockVisibilityManager)(nil).GetClosedWorkflowExecution), request)
}

// GetName mocks base method.
func (m *MockVisibilityManager) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockVisibilityManagerMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockVisibilityManager)(nil).GetName))
}

// ListClosedWorkflowExecutions mocks base method.
func (m *MockVisibilityManager) ListClosedWorkflowExecutions(request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClosedWorkflowExecutions", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClosedWorkflowExecutions indicates an expected call of ListClosedWorkflowExecutions.
func (mr *MockVisibilityManagerMockRecorder) ListClosedWorkflowExecutions(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClosedWorkflowExecutions", reflect.TypeOf((*MockVisibilityManager)(nil).ListClosedWorkflowExecutions), request)
}

// ListClosedWorkflowExecutionsByStatus mocks base method.
func (m *MockVisibilityManager) ListClosedWorkflowExecutionsByStatus(request *ListClosedWorkflowExecutionsByStatusRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClosedWorkflowExecutionsByStatus", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClosedWorkflowExecutionsByStatus indicates an expected call of ListClosedWorkflowExecutionsByStatus.
func (mr *MockVisibilityManagerMockRecorder) ListClosedWorkflowExecutionsByStatus(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClosedWorkflowExecutionsByStatus", reflect.TypeOf((*MockVisibilityManager)(nil).ListClosedWorkflowExecutionsByStatus), request)
}

// ListClosedWorkflowExecutionsByType mocks base method.
func (m *MockVisibilityManager) ListClosedWorkflowExecutionsByType(request *ListWorkflowExecutionsByTypeRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClosedWorkflowExecutionsByType", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClosedWorkflowExecutionsByType indicates an expected call of ListClosedWorkflowExecutionsByType.
func (mr *MockVisibilityManagerMockRecorder) ListClosedWorkflowExecutionsByType(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClosedWorkflowExecutionsByType", reflect.TypeOf((*MockVisibilityManager)(nil).ListClosedWorkflowExecutionsByType), request)
}

// ListClosedWorkflowExecutionsByWorkflowID mocks base method.
func (m *MockVisibilityManager) ListClosedWorkflowExecutionsByWorkflowID(request *ListWorkflowExecutionsByWorkflowIDRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClosedWorkflowExecutionsByWorkflowID", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClosedWorkflowExecutionsByWorkflowID indicates an expected call of ListClosedWorkflowExecutionsByWorkflowID.
func (mr *MockVisibilityManagerMockRecorder) ListClosedWorkflowExecutionsByWorkflowID(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClosedWorkflowExecutionsByWorkflowID", reflect.TypeOf((*MockVisibilityManager)(nil).ListClosedWorkflowExecutionsByWorkflowID), request)
}

// ListOpenWorkflowExecutions mocks base method.
func (m *MockVisibilityManager) ListOpenWorkflowExecutions(request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOpenWorkflowExecutions", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOpenWorkflowExecutions indicates an expected call of ListOpenWorkflowExecutions.
func (mr *MockVisibilityManagerMockRecorder) ListOpenWorkflowExecutions(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOpenWorkflowExecutions", reflect.TypeOf((*MockVisibilityManager)(nil).ListOpenWorkflowExecutions), request)
}

// ListOpenWorkflowExecutionsByType mocks base method.
func (m *MockVisibilityManager) ListOpenWorkflowExecutionsByType(request *ListWorkflowExecutionsByTypeRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOpenWorkflowExecutionsByType", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOpenWorkflowExecutionsByType indicates an expected call of ListOpenWorkflowExecutionsByType.
func (mr *MockVisibilityManagerMockRecorder) ListOpenWorkflowExecutionsByType(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOpenWorkflowExecutionsByType", reflect.TypeOf((*MockVisibilityManager)(nil).ListOpenWorkflowExecutionsByType), request)
}

// ListOpenWorkflowExecutionsByWorkflowID mocks base method.
func (m *MockVisibilityManager) ListOpenWorkflowExecutionsByWorkflowID(request *ListWorkflowExecutionsByWorkflowIDRequest) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOpenWorkflowExecutionsByWorkflowID", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOpenWorkflowExecutionsByWorkflowID indicates an expected call of ListOpenWorkflowExecutionsByWorkflowID.
func (mr *MockVisibilityManagerMockRecorder) ListOpenWorkflowExecutionsByWorkflowID(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOpenWorkflowExecutionsByWorkflowID", reflect.TypeOf((*MockVisibilityManager)(nil).ListOpenWorkflowExecutionsByWorkflowID), request)
}

// ListWorkflowExecutions mocks base method.
func (m *MockVisibilityManager) ListWorkflowExecutions(request *ListWorkflowExecutionsRequestV2) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkflowExecutions", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkflowExecutions indicates an expected call of ListWorkflowExecutions.
func (mr *MockVisibilityManagerMockRecorder) ListWorkflowExecutions(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkflowExecutions", reflect.TypeOf((*MockVisibilityManager)(nil).ListWorkflowExecutions), request)
}

// RecordWorkflowExecutionClosed mocks base method.
func (m *MockVisibilityManager) RecordWorkflowExecutionClosed(request *RecordWorkflowExecutionClosedRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordWorkflowExecutionClosed", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordWorkflowExecutionClosed indicates an expected call of RecordWorkflowExecutionClosed.
func (mr *MockVisibilityManagerMockRecorder) RecordWorkflowExecutionClosed(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordWorkflowExecutionClosed", reflect.TypeOf((*MockVisibilityManager)(nil).RecordWorkflowExecutionClosed), request)
}

// RecordWorkflowExecutionStarted mocks base method.
func (m *MockVisibilityManager) RecordWorkflowExecutionStarted(request *RecordWorkflowExecutionStartedRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordWorkflowExecutionStarted", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordWorkflowExecutionStarted indicates an expected call of RecordWorkflowExecutionStarted.
func (mr *MockVisibilityManagerMockRecorder) RecordWorkflowExecutionStarted(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordWorkflowExecutionStarted", reflect.TypeOf((*MockVisibilityManager)(nil).RecordWorkflowExecutionStarted), request)
}

// ScanWorkflowExecutions mocks base method.
func (m *MockVisibilityManager) ScanWorkflowExecutions(request *ListWorkflowExecutionsRequestV2) (*ListWorkflowExecutionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanWorkflowExecutions", request)
	ret0, _ := ret[0].(*ListWorkflowExecutionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ScanWorkflowExecutions indicates an expected call of ScanWorkflowExecutions.
func (mr *MockVisibilityManagerMockRecorder) ScanWorkflowExecutions(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanWorkflowExecutions", reflect.TypeOf((*MockVisibilityManager)(nil).ScanWorkflowExecutions), request)
}

// UpsertWorkflowExecution mocks base method.
func (m *MockVisibilityManager) UpsertWorkflowExecution(request *UpsertWorkflowExecutionRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertWorkflowExecution", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertWorkflowExecution indicates an expected call of UpsertWorkflowExecution.
func (mr *MockVisibilityManagerMockRecorder) UpsertWorkflowExecution(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertWorkflowExecution", reflect.TypeOf((*MockVisibilityManager)(nil).UpsertWorkflowExecution), request)
}
