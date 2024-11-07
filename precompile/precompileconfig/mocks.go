// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/coreth/precompile/precompileconfig (interfaces: Predicater,Config,ChainConfig,Accepter)
//
// Generated by this command:
//
//	mockgen -package=precompileconfig -destination=precompile/precompileconfig/mocks.go github.com/ava-labs/coreth/precompile/precompileconfig Predicater,Config,ChainConfig,Accepter
//

// Package precompileconfig is a generated GoMock package.
package precompileconfig

import (
	reflect "reflect"

	commontype "github.com/ava-labs/coreth/commontype"
	common "github.com/ava-labs/libevm/common"
	gomock "go.uber.org/mock/gomock"
)

// MockPredicater is a mock of Predicater interface.
type MockPredicater struct {
	ctrl     *gomock.Controller
	recorder *MockPredicaterMockRecorder
}

// MockPredicaterMockRecorder is the mock recorder for MockPredicater.
type MockPredicaterMockRecorder struct {
	mock *MockPredicater
}

// NewMockPredicater creates a new mock instance.
func NewMockPredicater(ctrl *gomock.Controller) *MockPredicater {
	mock := &MockPredicater{ctrl: ctrl}
	mock.recorder = &MockPredicaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPredicater) EXPECT() *MockPredicaterMockRecorder {
	return m.recorder
}

// PredicateGas mocks base method.
func (m *MockPredicater) PredicateGas(arg0 []byte) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PredicateGas", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PredicateGas indicates an expected call of PredicateGas.
func (mr *MockPredicaterMockRecorder) PredicateGas(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PredicateGas", reflect.TypeOf((*MockPredicater)(nil).PredicateGas), arg0)
}

// VerifyPredicate mocks base method.
func (m *MockPredicater) VerifyPredicate(arg0 *PredicateContext, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPredicate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyPredicate indicates an expected call of VerifyPredicate.
func (mr *MockPredicaterMockRecorder) VerifyPredicate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPredicate", reflect.TypeOf((*MockPredicater)(nil).VerifyPredicate), arg0, arg1)
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Equal mocks base method.
func (m *MockConfig) Equal(arg0 Config) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Equal", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Equal indicates an expected call of Equal.
func (mr *MockConfigMockRecorder) Equal(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Equal", reflect.TypeOf((*MockConfig)(nil).Equal), arg0)
}

// IsDisabled mocks base method.
func (m *MockConfig) IsDisabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDisabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDisabled indicates an expected call of IsDisabled.
func (mr *MockConfigMockRecorder) IsDisabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDisabled", reflect.TypeOf((*MockConfig)(nil).IsDisabled))
}

// Key mocks base method.
func (m *MockConfig) Key() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key")
	ret0, _ := ret[0].(string)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockConfigMockRecorder) Key() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockConfig)(nil).Key))
}

// Timestamp mocks base method.
func (m *MockConfig) Timestamp() *uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timestamp")
	ret0, _ := ret[0].(*uint64)
	return ret0
}

// Timestamp indicates an expected call of Timestamp.
func (mr *MockConfigMockRecorder) Timestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timestamp", reflect.TypeOf((*MockConfig)(nil).Timestamp))
}

// Verify mocks base method.
func (m *MockConfig) Verify(arg0 ChainConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockConfigMockRecorder) Verify(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockConfig)(nil).Verify), arg0)
}

// MockChainConfig is a mock of ChainConfig interface.
type MockChainConfig struct {
	ctrl     *gomock.Controller
	recorder *MockChainConfigMockRecorder
}

// MockChainConfigMockRecorder is the mock recorder for MockChainConfig.
type MockChainConfigMockRecorder struct {
	mock *MockChainConfig
}

// NewMockChainConfig creates a new mock instance.
func NewMockChainConfig(ctrl *gomock.Controller) *MockChainConfig {
	mock := &MockChainConfig{ctrl: ctrl}
	mock.recorder = &MockChainConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainConfig) EXPECT() *MockChainConfigMockRecorder {
	return m.recorder
}

// AllowedFeeRecipients mocks base method.
func (m *MockChainConfig) AllowedFeeRecipients() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowedFeeRecipients")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AllowedFeeRecipients indicates an expected call of AllowedFeeRecipients.
func (mr *MockChainConfigMockRecorder) AllowedFeeRecipients() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowedFeeRecipients", reflect.TypeOf((*MockChainConfig)(nil).AllowedFeeRecipients))
}

// GetFeeConfig mocks base method.
func (m *MockChainConfig) GetFeeConfig() commontype.FeeConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeeConfig")
	ret0, _ := ret[0].(commontype.FeeConfig)
	return ret0
}

// GetFeeConfig indicates an expected call of GetFeeConfig.
func (mr *MockChainConfigMockRecorder) GetFeeConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeeConfig", reflect.TypeOf((*MockChainConfig)(nil).GetFeeConfig))
}

// IsDurango mocks base method.
func (m *MockChainConfig) IsDurango(arg0 uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDurango", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDurango indicates an expected call of IsDurango.
func (mr *MockChainConfigMockRecorder) IsDurango(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDurango", reflect.TypeOf((*MockChainConfig)(nil).IsDurango), arg0)
}

// MockAccepter is a mock of Accepter interface.
type MockAccepter struct {
	ctrl     *gomock.Controller
	recorder *MockAccepterMockRecorder
}

// MockAccepterMockRecorder is the mock recorder for MockAccepter.
type MockAccepterMockRecorder struct {
	mock *MockAccepter
}

// NewMockAccepter creates a new mock instance.
func NewMockAccepter(ctrl *gomock.Controller) *MockAccepter {
	mock := &MockAccepter{ctrl: ctrl}
	mock.recorder = &MockAccepterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccepter) EXPECT() *MockAccepterMockRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *MockAccepter) Accept(arg0 *AcceptContext, arg1 common.Hash, arg2 uint64, arg3 common.Hash, arg4 int, arg5 []common.Hash, arg6 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *MockAccepterMockRecorder) Accept(arg0, arg1, arg2, arg3, arg4, arg5, arg6 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockAccepter)(nil).Accept), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}
