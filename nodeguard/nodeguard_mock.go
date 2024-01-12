// Code generated by MockGen. DO NOT EDIT.
// Source: nodeguard/nodeguard_grpc.pb.go

// Package nodeguard is a generated GoMock package.
package nodeguard

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockNodeGuardServiceClient is a mock of NodeGuardServiceClient interface.
type MockNodeGuardServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNodeGuardServiceClientMockRecorder
}

// MockNodeGuardServiceClientMockRecorder is the mock recorder for MockNodeGuardServiceClient.
type MockNodeGuardServiceClientMockRecorder struct {
	mock *MockNodeGuardServiceClient
}

// NewMockNodeGuardServiceClient creates a new mock instance.
func NewMockNodeGuardServiceClient(ctrl *gomock.Controller) *MockNodeGuardServiceClient {
	mock := &MockNodeGuardServiceClient{ctrl: ctrl}
	mock.recorder = &MockNodeGuardServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeGuardServiceClient) EXPECT() *MockNodeGuardServiceClientMockRecorder {
	return m.recorder
}

// AddLiquidityRule mocks base method.
func (m *MockNodeGuardServiceClient) AddLiquidityRule(ctx context.Context, in *AddLiquidityRuleRequest, opts ...grpc.CallOption) (*AddLiquidityRuleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddLiquidityRule", varargs...)
	ret0, _ := ret[0].(*AddLiquidityRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLiquidityRule indicates an expected call of AddLiquidityRule.
func (mr *MockNodeGuardServiceClientMockRecorder) AddLiquidityRule(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLiquidityRule", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).AddLiquidityRule), varargs...)
}

// AddNode mocks base method.
func (m *MockNodeGuardServiceClient) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddNode", varargs...)
	ret0, _ := ret[0].(*AddNodeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNode indicates an expected call of AddNode.
func (mr *MockNodeGuardServiceClientMockRecorder) AddNode(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNode", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).AddNode), varargs...)
}

// CloseChannel mocks base method.
func (m *MockNodeGuardServiceClient) CloseChannel(ctx context.Context, in *CloseChannelRequest, opts ...grpc.CallOption) (*CloseChannelResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CloseChannel", varargs...)
	ret0, _ := ret[0].(*CloseChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseChannel indicates an expected call of CloseChannel.
func (mr *MockNodeGuardServiceClientMockRecorder) CloseChannel(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseChannel", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).CloseChannel), varargs...)
}

// GetAvailableUtxos mocks base method.
func (m *MockNodeGuardServiceClient) GetAvailableUtxos(ctx context.Context, in *GetAvailableUtxosRequest, opts ...grpc.CallOption) (*GetAvailableUtxosResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAvailableUtxos", varargs...)
	ret0, _ := ret[0].(*GetAvailableUtxosResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableUtxos indicates an expected call of GetAvailableUtxos.
func (mr *MockNodeGuardServiceClientMockRecorder) GetAvailableUtxos(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableUtxos", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetAvailableUtxos), varargs...)
}

// GetAvailableWallets mocks base method.
func (m *MockNodeGuardServiceClient) GetAvailableWallets(ctx context.Context, in *GetAvailableWalletsRequest, opts ...grpc.CallOption) (*GetAvailableWalletsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAvailableWallets", varargs...)
	ret0, _ := ret[0].(*GetAvailableWalletsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableWallets indicates an expected call of GetAvailableWallets.
func (mr *MockNodeGuardServiceClientMockRecorder) GetAvailableWallets(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableWallets", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetAvailableWallets), varargs...)
}

// GetChannel mocks base method.
func (m *MockNodeGuardServiceClient) GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*GetChannelResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetChannel", varargs...)
	ret0, _ := ret[0].(*GetChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannel indicates an expected call of GetChannel.
func (mr *MockNodeGuardServiceClientMockRecorder) GetChannel(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannel", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetChannel), varargs...)
}

// GetChannelOperationRequest mocks base method.
func (m *MockNodeGuardServiceClient) GetChannelOperationRequest(ctx context.Context, in *GetChannelOperationRequestRequest, opts ...grpc.CallOption) (*GetChannelOperationRequestResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetChannelOperationRequest", varargs...)
	ret0, _ := ret[0].(*GetChannelOperationRequestResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannelOperationRequest indicates an expected call of GetChannelOperationRequest.
func (mr *MockNodeGuardServiceClientMockRecorder) GetChannelOperationRequest(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelOperationRequest", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetChannelOperationRequest), varargs...)
}

// GetLiquidityRules mocks base method.
func (m *MockNodeGuardServiceClient) GetLiquidityRules(ctx context.Context, in *GetLiquidityRulesRequest, opts ...grpc.CallOption) (*GetLiquidityRulesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLiquidityRules", varargs...)
	ret0, _ := ret[0].(*GetLiquidityRulesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLiquidityRules indicates an expected call of GetLiquidityRules.
func (mr *MockNodeGuardServiceClientMockRecorder) GetLiquidityRules(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiquidityRules", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetLiquidityRules), varargs...)
}

// GetNewWalletAddress mocks base method.
func (m *MockNodeGuardServiceClient) GetNewWalletAddress(ctx context.Context, in *GetNewWalletAddressRequest, opts ...grpc.CallOption) (*GetNewWalletAddressResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNewWalletAddress", varargs...)
	ret0, _ := ret[0].(*GetNewWalletAddressResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewWalletAddress indicates an expected call of GetNewWalletAddress.
func (mr *MockNodeGuardServiceClientMockRecorder) GetNewWalletAddress(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewWalletAddress", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetNewWalletAddress), varargs...)
}

// GetNodes mocks base method.
func (m *MockNodeGuardServiceClient) GetNodes(ctx context.Context, in *GetNodesRequest, opts ...grpc.CallOption) (*GetNodesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNodes", varargs...)
	ret0, _ := ret[0].(*GetNodesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes.
func (mr *MockNodeGuardServiceClientMockRecorder) GetNodes(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetNodes), varargs...)
}

// GetWithdrawalsRequestStatus mocks base method.
func (m *MockNodeGuardServiceClient) GetWithdrawalsRequestStatus(ctx context.Context, in *GetWithdrawalsRequestStatusRequest, opts ...grpc.CallOption) (*GetWithdrawalsRequestStatusResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetWithdrawalsRequestStatus", varargs...)
	ret0, _ := ret[0].(*GetWithdrawalsRequestStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsRequestStatus indicates an expected call of GetWithdrawalsRequestStatus.
func (mr *MockNodeGuardServiceClientMockRecorder) GetWithdrawalsRequestStatus(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsRequestStatus", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).GetWithdrawalsRequestStatus), varargs...)
}

// OpenChannel mocks base method.
func (m *MockNodeGuardServiceClient) OpenChannel(ctx context.Context, in *OpenChannelRequest, opts ...grpc.CallOption) (*OpenChannelResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OpenChannel", varargs...)
	ret0, _ := ret[0].(*OpenChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenChannel indicates an expected call of OpenChannel.
func (mr *MockNodeGuardServiceClientMockRecorder) OpenChannel(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenChannel", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).OpenChannel), varargs...)
}

// RequestWithdrawal mocks base method.
func (m *MockNodeGuardServiceClient) RequestWithdrawal(ctx context.Context, in *RequestWithdrawalRequest, opts ...grpc.CallOption) (*RequestWithdrawalResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RequestWithdrawal", varargs...)
	ret0, _ := ret[0].(*RequestWithdrawalResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestWithdrawal indicates an expected call of RequestWithdrawal.
func (mr *MockNodeGuardServiceClientMockRecorder) RequestWithdrawal(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestWithdrawal", reflect.TypeOf((*MockNodeGuardServiceClient)(nil).RequestWithdrawal), varargs...)
}

// MockNodeGuardServiceServer is a mock of NodeGuardServiceServer interface.
type MockNodeGuardServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNodeGuardServiceServerMockRecorder
}

// MockNodeGuardServiceServerMockRecorder is the mock recorder for MockNodeGuardServiceServer.
type MockNodeGuardServiceServerMockRecorder struct {
	mock *MockNodeGuardServiceServer
}

// NewMockNodeGuardServiceServer creates a new mock instance.
func NewMockNodeGuardServiceServer(ctrl *gomock.Controller) *MockNodeGuardServiceServer {
	mock := &MockNodeGuardServiceServer{ctrl: ctrl}
	mock.recorder = &MockNodeGuardServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeGuardServiceServer) EXPECT() *MockNodeGuardServiceServerMockRecorder {
	return m.recorder
}

// AddLiquidityRule mocks base method.
func (m *MockNodeGuardServiceServer) AddLiquidityRule(arg0 context.Context, arg1 *AddLiquidityRuleRequest) (*AddLiquidityRuleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLiquidityRule", arg0, arg1)
	ret0, _ := ret[0].(*AddLiquidityRuleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLiquidityRule indicates an expected call of AddLiquidityRule.
func (mr *MockNodeGuardServiceServerMockRecorder) AddLiquidityRule(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLiquidityRule", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).AddLiquidityRule), arg0, arg1)
}

// AddNode mocks base method.
func (m *MockNodeGuardServiceServer) AddNode(arg0 context.Context, arg1 *AddNodeRequest) (*AddNodeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNode", arg0, arg1)
	ret0, _ := ret[0].(*AddNodeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNode indicates an expected call of AddNode.
func (mr *MockNodeGuardServiceServerMockRecorder) AddNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNode", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).AddNode), arg0, arg1)
}

// CloseChannel mocks base method.
func (m *MockNodeGuardServiceServer) CloseChannel(arg0 context.Context, arg1 *CloseChannelRequest) (*CloseChannelResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseChannel", arg0, arg1)
	ret0, _ := ret[0].(*CloseChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseChannel indicates an expected call of CloseChannel.
func (mr *MockNodeGuardServiceServerMockRecorder) CloseChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseChannel", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).CloseChannel), arg0, arg1)
}

// GetAvailableUtxos mocks base method.
func (m *MockNodeGuardServiceServer) GetAvailableUtxos(arg0 context.Context, arg1 *GetAvailableUtxosRequest) (*GetAvailableUtxosResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableUtxos", arg0, arg1)
	ret0, _ := ret[0].(*GetAvailableUtxosResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableUtxos indicates an expected call of GetAvailableUtxos.
func (mr *MockNodeGuardServiceServerMockRecorder) GetAvailableUtxos(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableUtxos", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetAvailableUtxos), arg0, arg1)
}

// GetAvailableWallets mocks base method.
func (m *MockNodeGuardServiceServer) GetAvailableWallets(arg0 context.Context, arg1 *GetAvailableWalletsRequest) (*GetAvailableWalletsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableWallets", arg0, arg1)
	ret0, _ := ret[0].(*GetAvailableWalletsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableWallets indicates an expected call of GetAvailableWallets.
func (mr *MockNodeGuardServiceServerMockRecorder) GetAvailableWallets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableWallets", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetAvailableWallets), arg0, arg1)
}

// GetChannel mocks base method.
func (m *MockNodeGuardServiceServer) GetChannel(arg0 context.Context, arg1 *GetChannelRequest) (*GetChannelResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannel", arg0, arg1)
	ret0, _ := ret[0].(*GetChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannel indicates an expected call of GetChannel.
func (mr *MockNodeGuardServiceServerMockRecorder) GetChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannel", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetChannel), arg0, arg1)
}

// GetChannelOperationRequest mocks base method.
func (m *MockNodeGuardServiceServer) GetChannelOperationRequest(arg0 context.Context, arg1 *GetChannelOperationRequestRequest) (*GetChannelOperationRequestResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannelOperationRequest", arg0, arg1)
	ret0, _ := ret[0].(*GetChannelOperationRequestResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannelOperationRequest indicates an expected call of GetChannelOperationRequest.
func (mr *MockNodeGuardServiceServerMockRecorder) GetChannelOperationRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannelOperationRequest", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetChannelOperationRequest), arg0, arg1)
}

// GetLiquidityRules mocks base method.
func (m *MockNodeGuardServiceServer) GetLiquidityRules(arg0 context.Context, arg1 *GetLiquidityRulesRequest) (*GetLiquidityRulesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLiquidityRules", arg0, arg1)
	ret0, _ := ret[0].(*GetLiquidityRulesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLiquidityRules indicates an expected call of GetLiquidityRules.
func (mr *MockNodeGuardServiceServerMockRecorder) GetLiquidityRules(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiquidityRules", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetLiquidityRules), arg0, arg1)
}

// GetNewWalletAddress mocks base method.
func (m *MockNodeGuardServiceServer) GetNewWalletAddress(arg0 context.Context, arg1 *GetNewWalletAddressRequest) (*GetNewWalletAddressResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewWalletAddress", arg0, arg1)
	ret0, _ := ret[0].(*GetNewWalletAddressResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewWalletAddress indicates an expected call of GetNewWalletAddress.
func (mr *MockNodeGuardServiceServerMockRecorder) GetNewWalletAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewWalletAddress", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetNewWalletAddress), arg0, arg1)
}

// GetNodes mocks base method.
func (m *MockNodeGuardServiceServer) GetNodes(arg0 context.Context, arg1 *GetNodesRequest) (*GetNodesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes", arg0, arg1)
	ret0, _ := ret[0].(*GetNodesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes.
func (mr *MockNodeGuardServiceServerMockRecorder) GetNodes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetNodes), arg0, arg1)
}

// GetWithdrawalsRequestStatus mocks base method.
func (m *MockNodeGuardServiceServer) GetWithdrawalsRequestStatus(arg0 context.Context, arg1 *GetWithdrawalsRequestStatusRequest) (*GetWithdrawalsRequestStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsRequestStatus", arg0, arg1)
	ret0, _ := ret[0].(*GetWithdrawalsRequestStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsRequestStatus indicates an expected call of GetWithdrawalsRequestStatus.
func (mr *MockNodeGuardServiceServerMockRecorder) GetWithdrawalsRequestStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsRequestStatus", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).GetWithdrawalsRequestStatus), arg0, arg1)
}

// OpenChannel mocks base method.
func (m *MockNodeGuardServiceServer) OpenChannel(arg0 context.Context, arg1 *OpenChannelRequest) (*OpenChannelResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenChannel", arg0, arg1)
	ret0, _ := ret[0].(*OpenChannelResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenChannel indicates an expected call of OpenChannel.
func (mr *MockNodeGuardServiceServerMockRecorder) OpenChannel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenChannel", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).OpenChannel), arg0, arg1)
}

// RequestWithdrawal mocks base method.
func (m *MockNodeGuardServiceServer) RequestWithdrawal(arg0 context.Context, arg1 *RequestWithdrawalRequest) (*RequestWithdrawalResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestWithdrawal", arg0, arg1)
	ret0, _ := ret[0].(*RequestWithdrawalResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestWithdrawal indicates an expected call of RequestWithdrawal.
func (mr *MockNodeGuardServiceServerMockRecorder) RequestWithdrawal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestWithdrawal", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).RequestWithdrawal), arg0, arg1)
}

// mustEmbedUnimplementedNodeGuardServiceServer mocks base method.
func (m *MockNodeGuardServiceServer) mustEmbedUnimplementedNodeGuardServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNodeGuardServiceServer")
}

// mustEmbedUnimplementedNodeGuardServiceServer indicates an expected call of mustEmbedUnimplementedNodeGuardServiceServer.
func (mr *MockNodeGuardServiceServerMockRecorder) mustEmbedUnimplementedNodeGuardServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNodeGuardServiceServer", reflect.TypeOf((*MockNodeGuardServiceServer)(nil).mustEmbedUnimplementedNodeGuardServiceServer))
}

// MockUnsafeNodeGuardServiceServer is a mock of UnsafeNodeGuardServiceServer interface.
type MockUnsafeNodeGuardServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNodeGuardServiceServerMockRecorder
}

// MockUnsafeNodeGuardServiceServerMockRecorder is the mock recorder for MockUnsafeNodeGuardServiceServer.
type MockUnsafeNodeGuardServiceServerMockRecorder struct {
	mock *MockUnsafeNodeGuardServiceServer
}

// NewMockUnsafeNodeGuardServiceServer creates a new mock instance.
func NewMockUnsafeNodeGuardServiceServer(ctrl *gomock.Controller) *MockUnsafeNodeGuardServiceServer {
	mock := &MockUnsafeNodeGuardServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNodeGuardServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNodeGuardServiceServer) EXPECT() *MockUnsafeNodeGuardServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNodeGuardServiceServer mocks base method.
func (m *MockUnsafeNodeGuardServiceServer) mustEmbedUnimplementedNodeGuardServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNodeGuardServiceServer")
}

// mustEmbedUnimplementedNodeGuardServiceServer indicates an expected call of mustEmbedUnimplementedNodeGuardServiceServer.
func (mr *MockUnsafeNodeGuardServiceServerMockRecorder) mustEmbedUnimplementedNodeGuardServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNodeGuardServiceServer", reflect.TypeOf((*MockUnsafeNodeGuardServiceServer)(nil).mustEmbedUnimplementedNodeGuardServiceServer))
}
