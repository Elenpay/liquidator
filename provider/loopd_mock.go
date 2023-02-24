// Code generated by MockGen. DO NOT EDIT.
// Source: loop/looprpc/client_grpc.pb.go

// Package provider is a generated GoMock package.
package provider

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	looprpc "github.com/lightninglabs/loop/looprpc"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockSwapClientClient is a mock of SwapClientClient interface.
type MockSwapClientClient struct {
	ctrl     *gomock.Controller
	recorder *MockSwapClientClientMockRecorder
}

// MockSwapClientClientMockRecorder is the mock recorder for MockSwapClientClient.
type MockSwapClientClientMockRecorder struct {
	mock *MockSwapClientClient
}

// NewMockSwapClientClient creates a new mock instance.
func NewMockSwapClientClient(ctrl *gomock.Controller) *MockSwapClientClient {
	mock := &MockSwapClientClient{ctrl: ctrl}
	mock.recorder = &MockSwapClientClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapClientClient) EXPECT() *MockSwapClientClientMockRecorder {
	return m.recorder
}

// GetLiquidityParams mocks base method.
func (m *MockSwapClientClient) GetLiquidityParams(ctx context.Context, in *looprpc.GetLiquidityParamsRequest, opts ...grpc.CallOption) (*looprpc.LiquidityParameters, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLiquidityParams", varargs...)
	ret0, _ := ret[0].(*looprpc.LiquidityParameters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLiquidityParams indicates an expected call of GetLiquidityParams.
func (mr *MockSwapClientClientMockRecorder) GetLiquidityParams(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiquidityParams", reflect.TypeOf((*MockSwapClientClient)(nil).GetLiquidityParams), varargs...)
}

// GetLoopInQuote mocks base method.
func (m *MockSwapClientClient) GetLoopInQuote(ctx context.Context, in *looprpc.QuoteRequest, opts ...grpc.CallOption) (*looprpc.InQuoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLoopInQuote", varargs...)
	ret0, _ := ret[0].(*looprpc.InQuoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoopInQuote indicates an expected call of GetLoopInQuote.
func (mr *MockSwapClientClientMockRecorder) GetLoopInQuote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoopInQuote", reflect.TypeOf((*MockSwapClientClient)(nil).GetLoopInQuote), varargs...)
}

// GetLoopInTerms mocks base method.
func (m *MockSwapClientClient) GetLoopInTerms(ctx context.Context, in *looprpc.TermsRequest, opts ...grpc.CallOption) (*looprpc.InTermsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLoopInTerms", varargs...)
	ret0, _ := ret[0].(*looprpc.InTermsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoopInTerms indicates an expected call of GetLoopInTerms.
func (mr *MockSwapClientClientMockRecorder) GetLoopInTerms(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoopInTerms", reflect.TypeOf((*MockSwapClientClient)(nil).GetLoopInTerms), varargs...)
}

// GetLsatTokens mocks base method.
func (m *MockSwapClientClient) GetLsatTokens(ctx context.Context, in *looprpc.TokensRequest, opts ...grpc.CallOption) (*looprpc.TokensResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLsatTokens", varargs...)
	ret0, _ := ret[0].(*looprpc.TokensResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLsatTokens indicates an expected call of GetLsatTokens.
func (mr *MockSwapClientClientMockRecorder) GetLsatTokens(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLsatTokens", reflect.TypeOf((*MockSwapClientClient)(nil).GetLsatTokens), varargs...)
}

// ListSwaps mocks base method.
func (m *MockSwapClientClient) ListSwaps(ctx context.Context, in *looprpc.ListSwapsRequest, opts ...grpc.CallOption) (*looprpc.ListSwapsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListSwaps", varargs...)
	ret0, _ := ret[0].(*looprpc.ListSwapsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSwaps indicates an expected call of ListSwaps.
func (mr *MockSwapClientClientMockRecorder) ListSwaps(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSwaps", reflect.TypeOf((*MockSwapClientClient)(nil).ListSwaps), varargs...)
}

// LoopIn mocks base method.
func (m *MockSwapClientClient) LoopIn(ctx context.Context, in *looprpc.LoopInRequest, opts ...grpc.CallOption) (*looprpc.SwapResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoopIn", varargs...)
	ret0, _ := ret[0].(*looprpc.SwapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopIn indicates an expected call of LoopIn.
func (mr *MockSwapClientClientMockRecorder) LoopIn(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopIn", reflect.TypeOf((*MockSwapClientClient)(nil).LoopIn), varargs...)
}

// LoopOut mocks base method.
func (m *MockSwapClientClient) LoopOut(ctx context.Context, in *looprpc.LoopOutRequest, opts ...grpc.CallOption) (*looprpc.SwapResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoopOut", varargs...)
	ret0, _ := ret[0].(*looprpc.SwapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOut indicates an expected call of LoopOut.
func (mr *MockSwapClientClientMockRecorder) LoopOut(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOut", reflect.TypeOf((*MockSwapClientClient)(nil).LoopOut), varargs...)
}

// LoopOutQuote mocks base method.
func (m *MockSwapClientClient) LoopOutQuote(ctx context.Context, in *looprpc.QuoteRequest, opts ...grpc.CallOption) (*looprpc.OutQuoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoopOutQuote", varargs...)
	ret0, _ := ret[0].(*looprpc.OutQuoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOutQuote indicates an expected call of LoopOutQuote.
func (mr *MockSwapClientClientMockRecorder) LoopOutQuote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOutQuote", reflect.TypeOf((*MockSwapClientClient)(nil).LoopOutQuote), varargs...)
}

// LoopOutTerms mocks base method.
func (m *MockSwapClientClient) LoopOutTerms(ctx context.Context, in *looprpc.TermsRequest, opts ...grpc.CallOption) (*looprpc.OutTermsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LoopOutTerms", varargs...)
	ret0, _ := ret[0].(*looprpc.OutTermsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOutTerms indicates an expected call of LoopOutTerms.
func (mr *MockSwapClientClientMockRecorder) LoopOutTerms(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOutTerms", reflect.TypeOf((*MockSwapClientClient)(nil).LoopOutTerms), varargs...)
}

// Monitor mocks base method.
func (m *MockSwapClientClient) Monitor(ctx context.Context, in *looprpc.MonitorRequest, opts ...grpc.CallOption) (looprpc.SwapClient_MonitorClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Monitor", varargs...)
	ret0, _ := ret[0].(looprpc.SwapClient_MonitorClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Monitor indicates an expected call of Monitor.
func (mr *MockSwapClientClientMockRecorder) Monitor(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Monitor", reflect.TypeOf((*MockSwapClientClient)(nil).Monitor), varargs...)
}

// Probe mocks base method.
func (m *MockSwapClientClient) Probe(ctx context.Context, in *looprpc.ProbeRequest, opts ...grpc.CallOption) (*looprpc.ProbeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Probe", varargs...)
	ret0, _ := ret[0].(*looprpc.ProbeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Probe indicates an expected call of Probe.
func (mr *MockSwapClientClientMockRecorder) Probe(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Probe", reflect.TypeOf((*MockSwapClientClient)(nil).Probe), varargs...)
}

// SetLiquidityParams mocks base method.
func (m *MockSwapClientClient) SetLiquidityParams(ctx context.Context, in *looprpc.SetLiquidityParamsRequest, opts ...grpc.CallOption) (*looprpc.SetLiquidityParamsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetLiquidityParams", varargs...)
	ret0, _ := ret[0].(*looprpc.SetLiquidityParamsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLiquidityParams indicates an expected call of SetLiquidityParams.
func (mr *MockSwapClientClientMockRecorder) SetLiquidityParams(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLiquidityParams", reflect.TypeOf((*MockSwapClientClient)(nil).SetLiquidityParams), varargs...)
}

// SuggestSwaps mocks base method.
func (m *MockSwapClientClient) SuggestSwaps(ctx context.Context, in *looprpc.SuggestSwapsRequest, opts ...grpc.CallOption) (*looprpc.SuggestSwapsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SuggestSwaps", varargs...)
	ret0, _ := ret[0].(*looprpc.SuggestSwapsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestSwaps indicates an expected call of SuggestSwaps.
func (mr *MockSwapClientClientMockRecorder) SuggestSwaps(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestSwaps", reflect.TypeOf((*MockSwapClientClient)(nil).SuggestSwaps), varargs...)
}

// SwapInfo mocks base method.
func (m *MockSwapClientClient) SwapInfo(ctx context.Context, in *looprpc.SwapInfoRequest, opts ...grpc.CallOption) (*looprpc.SwapStatus, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SwapInfo", varargs...)
	ret0, _ := ret[0].(*looprpc.SwapStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SwapInfo indicates an expected call of SwapInfo.
func (mr *MockSwapClientClientMockRecorder) SwapInfo(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwapInfo", reflect.TypeOf((*MockSwapClientClient)(nil).SwapInfo), varargs...)
}

// MockSwapClient_MonitorClient is a mock of SwapClient_MonitorClient interface.
type MockSwapClient_MonitorClient struct {
	ctrl     *gomock.Controller
	recorder *MockSwapClient_MonitorClientMockRecorder
}

// MockSwapClient_MonitorClientMockRecorder is the mock recorder for MockSwapClient_MonitorClient.
type MockSwapClient_MonitorClientMockRecorder struct {
	mock *MockSwapClient_MonitorClient
}

// NewMockSwapClient_MonitorClient creates a new mock instance.
func NewMockSwapClient_MonitorClient(ctrl *gomock.Controller) *MockSwapClient_MonitorClient {
	mock := &MockSwapClient_MonitorClient{ctrl: ctrl}
	mock.recorder = &MockSwapClient_MonitorClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapClient_MonitorClient) EXPECT() *MockSwapClient_MonitorClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockSwapClient_MonitorClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockSwapClient_MonitorClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockSwapClient_MonitorClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockSwapClient_MonitorClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).Context))
}

// Header mocks base method.
func (m *MockSwapClient_MonitorClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockSwapClient_MonitorClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockSwapClient_MonitorClient) Recv() (*looprpc.SwapStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*looprpc.SwapStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockSwapClient_MonitorClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockSwapClient_MonitorClient) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockSwapClient_MonitorClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).RecvMsg), m)
}

// SendMsg mocks base method.
func (m_2 *MockSwapClient_MonitorClient) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockSwapClient_MonitorClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockSwapClient_MonitorClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockSwapClient_MonitorClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockSwapClient_MonitorClient)(nil).Trailer))
}

// MockSwapClientServer is a mock of SwapClientServer interface.
type MockSwapClientServer struct {
	ctrl     *gomock.Controller
	recorder *MockSwapClientServerMockRecorder
}

// MockSwapClientServerMockRecorder is the mock recorder for MockSwapClientServer.
type MockSwapClientServerMockRecorder struct {
	mock *MockSwapClientServer
}

// NewMockSwapClientServer creates a new mock instance.
func NewMockSwapClientServer(ctrl *gomock.Controller) *MockSwapClientServer {
	mock := &MockSwapClientServer{ctrl: ctrl}
	mock.recorder = &MockSwapClientServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapClientServer) EXPECT() *MockSwapClientServerMockRecorder {
	return m.recorder
}

// GetLiquidityParams mocks base method.
func (m *MockSwapClientServer) GetLiquidityParams(arg0 context.Context, arg1 *looprpc.GetLiquidityParamsRequest) (*looprpc.LiquidityParameters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLiquidityParams", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.LiquidityParameters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLiquidityParams indicates an expected call of GetLiquidityParams.
func (mr *MockSwapClientServerMockRecorder) GetLiquidityParams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLiquidityParams", reflect.TypeOf((*MockSwapClientServer)(nil).GetLiquidityParams), arg0, arg1)
}

// GetLoopInQuote mocks base method.
func (m *MockSwapClientServer) GetLoopInQuote(arg0 context.Context, arg1 *looprpc.QuoteRequest) (*looprpc.InQuoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoopInQuote", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.InQuoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoopInQuote indicates an expected call of GetLoopInQuote.
func (mr *MockSwapClientServerMockRecorder) GetLoopInQuote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoopInQuote", reflect.TypeOf((*MockSwapClientServer)(nil).GetLoopInQuote), arg0, arg1)
}

// GetLoopInTerms mocks base method.
func (m *MockSwapClientServer) GetLoopInTerms(arg0 context.Context, arg1 *looprpc.TermsRequest) (*looprpc.InTermsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoopInTerms", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.InTermsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoopInTerms indicates an expected call of GetLoopInTerms.
func (mr *MockSwapClientServerMockRecorder) GetLoopInTerms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoopInTerms", reflect.TypeOf((*MockSwapClientServer)(nil).GetLoopInTerms), arg0, arg1)
}

// GetLsatTokens mocks base method.
func (m *MockSwapClientServer) GetLsatTokens(arg0 context.Context, arg1 *looprpc.TokensRequest) (*looprpc.TokensResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLsatTokens", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.TokensResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLsatTokens indicates an expected call of GetLsatTokens.
func (mr *MockSwapClientServerMockRecorder) GetLsatTokens(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLsatTokens", reflect.TypeOf((*MockSwapClientServer)(nil).GetLsatTokens), arg0, arg1)
}

// ListSwaps mocks base method.
func (m *MockSwapClientServer) ListSwaps(arg0 context.Context, arg1 *looprpc.ListSwapsRequest) (*looprpc.ListSwapsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSwaps", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.ListSwapsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSwaps indicates an expected call of ListSwaps.
func (mr *MockSwapClientServerMockRecorder) ListSwaps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSwaps", reflect.TypeOf((*MockSwapClientServer)(nil).ListSwaps), arg0, arg1)
}

// LoopIn mocks base method.
func (m *MockSwapClientServer) LoopIn(arg0 context.Context, arg1 *looprpc.LoopInRequest) (*looprpc.SwapResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoopIn", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.SwapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopIn indicates an expected call of LoopIn.
func (mr *MockSwapClientServerMockRecorder) LoopIn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopIn", reflect.TypeOf((*MockSwapClientServer)(nil).LoopIn), arg0, arg1)
}

// LoopOut mocks base method.
func (m *MockSwapClientServer) LoopOut(arg0 context.Context, arg1 *looprpc.LoopOutRequest) (*looprpc.SwapResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoopOut", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.SwapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOut indicates an expected call of LoopOut.
func (mr *MockSwapClientServerMockRecorder) LoopOut(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOut", reflect.TypeOf((*MockSwapClientServer)(nil).LoopOut), arg0, arg1)
}

// LoopOutQuote mocks base method.
func (m *MockSwapClientServer) LoopOutQuote(arg0 context.Context, arg1 *looprpc.QuoteRequest) (*looprpc.OutQuoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoopOutQuote", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.OutQuoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOutQuote indicates an expected call of LoopOutQuote.
func (mr *MockSwapClientServerMockRecorder) LoopOutQuote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOutQuote", reflect.TypeOf((*MockSwapClientServer)(nil).LoopOutQuote), arg0, arg1)
}

// LoopOutTerms mocks base method.
func (m *MockSwapClientServer) LoopOutTerms(arg0 context.Context, arg1 *looprpc.TermsRequest) (*looprpc.OutTermsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoopOutTerms", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.OutTermsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoopOutTerms indicates an expected call of LoopOutTerms.
func (mr *MockSwapClientServerMockRecorder) LoopOutTerms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoopOutTerms", reflect.TypeOf((*MockSwapClientServer)(nil).LoopOutTerms), arg0, arg1)
}

// Monitor mocks base method.
func (m *MockSwapClientServer) Monitor(arg0 *looprpc.MonitorRequest, arg1 looprpc.SwapClient_MonitorServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Monitor", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Monitor indicates an expected call of Monitor.
func (mr *MockSwapClientServerMockRecorder) Monitor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Monitor", reflect.TypeOf((*MockSwapClientServer)(nil).Monitor), arg0, arg1)
}

// Probe mocks base method.
func (m *MockSwapClientServer) Probe(arg0 context.Context, arg1 *looprpc.ProbeRequest) (*looprpc.ProbeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Probe", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.ProbeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Probe indicates an expected call of Probe.
func (mr *MockSwapClientServerMockRecorder) Probe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Probe", reflect.TypeOf((*MockSwapClientServer)(nil).Probe), arg0, arg1)
}

// SetLiquidityParams mocks base method.
func (m *MockSwapClientServer) SetLiquidityParams(arg0 context.Context, arg1 *looprpc.SetLiquidityParamsRequest) (*looprpc.SetLiquidityParamsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLiquidityParams", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.SetLiquidityParamsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLiquidityParams indicates an expected call of SetLiquidityParams.
func (mr *MockSwapClientServerMockRecorder) SetLiquidityParams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLiquidityParams", reflect.TypeOf((*MockSwapClientServer)(nil).SetLiquidityParams), arg0, arg1)
}

// SuggestSwaps mocks base method.
func (m *MockSwapClientServer) SuggestSwaps(arg0 context.Context, arg1 *looprpc.SuggestSwapsRequest) (*looprpc.SuggestSwapsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestSwaps", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.SuggestSwapsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestSwaps indicates an expected call of SuggestSwaps.
func (mr *MockSwapClientServerMockRecorder) SuggestSwaps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestSwaps", reflect.TypeOf((*MockSwapClientServer)(nil).SuggestSwaps), arg0, arg1)
}

// SwapInfo mocks base method.
func (m *MockSwapClientServer) SwapInfo(arg0 context.Context, arg1 *looprpc.SwapInfoRequest) (*looprpc.SwapStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SwapInfo", arg0, arg1)
	ret0, _ := ret[0].(*looprpc.SwapStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SwapInfo indicates an expected call of SwapInfo.
func (mr *MockSwapClientServerMockRecorder) SwapInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwapInfo", reflect.TypeOf((*MockSwapClientServer)(nil).SwapInfo), arg0, arg1)
}

// mustEmbedUnimplementedSwapClientServer mocks base method.
func (m *MockSwapClientServer) mustEmbedUnimplementedSwapClientServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedSwapClientServer")
}

// mustEmbedUnimplementedSwapClientServer indicates an expected call of mustEmbedUnimplementedSwapClientServer.
func (mr *MockSwapClientServerMockRecorder) mustEmbedUnimplementedSwapClientServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedSwapClientServer", reflect.TypeOf((*MockSwapClientServer)(nil).mustEmbedUnimplementedSwapClientServer))
}

// MockUnsafeSwapClientServer is a mock of UnsafeSwapClientServer interface.
type MockUnsafeSwapClientServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeSwapClientServerMockRecorder
}

// MockUnsafeSwapClientServerMockRecorder is the mock recorder for MockUnsafeSwapClientServer.
type MockUnsafeSwapClientServerMockRecorder struct {
	mock *MockUnsafeSwapClientServer
}

// NewMockUnsafeSwapClientServer creates a new mock instance.
func NewMockUnsafeSwapClientServer(ctrl *gomock.Controller) *MockUnsafeSwapClientServer {
	mock := &MockUnsafeSwapClientServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeSwapClientServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeSwapClientServer) EXPECT() *MockUnsafeSwapClientServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedSwapClientServer mocks base method.
func (m *MockUnsafeSwapClientServer) mustEmbedUnimplementedSwapClientServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedSwapClientServer")
}

// mustEmbedUnimplementedSwapClientServer indicates an expected call of mustEmbedUnimplementedSwapClientServer.
func (mr *MockUnsafeSwapClientServerMockRecorder) mustEmbedUnimplementedSwapClientServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedSwapClientServer", reflect.TypeOf((*MockUnsafeSwapClientServer)(nil).mustEmbedUnimplementedSwapClientServer))
}

// MockSwapClient_MonitorServer is a mock of SwapClient_MonitorServer interface.
type MockSwapClient_MonitorServer struct {
	ctrl     *gomock.Controller
	recorder *MockSwapClient_MonitorServerMockRecorder
}

// MockSwapClient_MonitorServerMockRecorder is the mock recorder for MockSwapClient_MonitorServer.
type MockSwapClient_MonitorServerMockRecorder struct {
	mock *MockSwapClient_MonitorServer
}

// NewMockSwapClient_MonitorServer creates a new mock instance.
func NewMockSwapClient_MonitorServer(ctrl *gomock.Controller) *MockSwapClient_MonitorServer {
	mock := &MockSwapClient_MonitorServer{ctrl: ctrl}
	mock.recorder = &MockSwapClient_MonitorServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapClient_MonitorServer) EXPECT() *MockSwapClient_MonitorServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockSwapClient_MonitorServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockSwapClient_MonitorServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m_2 *MockSwapClient_MonitorServer) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockSwapClient_MonitorServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockSwapClient_MonitorServer) Send(arg0 *looprpc.SwapStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockSwapClient_MonitorServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockSwapClient_MonitorServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockSwapClient_MonitorServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockSwapClient_MonitorServer) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockSwapClient_MonitorServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockSwapClient_MonitorServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockSwapClient_MonitorServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockSwapClient_MonitorServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockSwapClient_MonitorServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockSwapClient_MonitorServer)(nil).SetTrailer), arg0)
}
