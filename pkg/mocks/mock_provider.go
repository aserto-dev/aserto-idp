// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aserto-dev/aserto-idp/pkg/provider (interfaces: Provider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	grpcplugin "github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
	gomock "github.com/golang/mock/gomock"
)

// MockProvider is a mock of Provider interface.
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance.
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetName mocks base method.
func (m *MockProvider) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockProviderMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockProvider)(nil).GetName))
}

// GetPath mocks base method.
func (m *MockProvider) GetPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPath indicates an expected call of GetPath.
func (mr *MockProviderMockRecorder) GetPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPath", reflect.TypeOf((*MockProvider)(nil).GetPath))
}

// Kill mocks base method.
func (m *MockProvider) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockProviderMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockProvider)(nil).Kill))
}

// PluginClient mocks base method.
func (m *MockProvider) PluginClient() (grpcplugin.PluginClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PluginClient")
	ret0, _ := ret[0].(grpcplugin.PluginClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PluginClient indicates an expected call of PluginClient.
func (mr *MockProviderMockRecorder) PluginClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PluginClient", reflect.TypeOf((*MockProvider)(nil).PluginClient))
}
