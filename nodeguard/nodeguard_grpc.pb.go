// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.2
// source: nodeguard.proto

package nodeguard

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeGuardServiceClient is the client API for NodeGuardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeGuardServiceClient interface {
	// Returns the liquidity rules associated to a node and its channels
	GetLiquidityRules(ctx context.Context, in *GetLiquidityRulesRequest, opts ...grpc.CallOption) (*GetLiquidityRulesResponse, error)
	// Returns a new unused BTC Address of a given wallet
	GetNewWalletAddress(ctx context.Context, in *GetNewWalletAddressRequest, opts ...grpc.CallOption) (*GetNewWalletAddressResponse, error)
	// Withdraws funds from a given wallet to a given address
	RequestWithdrawal(ctx context.Context, in *RequestWithdrawalRequest, opts ...grpc.CallOption) (*RequestWithdrawalResponse, error)
	// Adds a new node to the nodeguard
	AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error)
	// Gets a list of nodes
	GetNodes(ctx context.Context, in *GetNodesRequest, opts ...grpc.CallOption) (*GetNodesResponse, error)
	// Gets a list of available wallets
	GetAvailableWallets(ctx context.Context, in *GetAvailableWalletsRequest, opts ...grpc.CallOption) (*GetAvailableWalletsResponse, error)
	// Opens a channel to a given node
	OpenChannel(ctx context.Context, in *OpenChannelRequest, opts ...grpc.CallOption) (*OpenChannelResponse, error)
	// Closes a channel to a given node
	CloseChannel(ctx context.Context, in *CloseChannelRequest, opts ...grpc.CallOption) (*CloseChannelResponse, error)
	// Gets a list of channel operations requests by id
	GetChannelOperationRequest(ctx context.Context, in *GetChannelOperationRequestRequest, opts ...grpc.CallOption) (*GetChannelOperationRequestResponse, error)
	// Adds a liquidity rule to a channel
	AddLiquidityRule(ctx context.Context, in *AddLiquidityRuleRequest, opts ...grpc.CallOption) (*AddLiquidityRuleResponse, error)
	// Gets a list of available UTXOs for a wallet
	GetAvailableUtxos(ctx context.Context, in *GetAvailableUtxosRequest, opts ...grpc.CallOption) (*GetAvailableUtxosResponse, error)
	// Gets the status for the provided withdrawals request ids
	GetWithdrawalsRequestStatus(ctx context.Context, in *GetWithdrawalsRequestStatusRequest, opts ...grpc.CallOption) (*GetWithdrawalsRequestStatusResponse, error)
	// Gets a channel by id
	GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*GetChannelResponse, error)
}

type nodeGuardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeGuardServiceClient(cc grpc.ClientConnInterface) NodeGuardServiceClient {
	return &nodeGuardServiceClient{cc}
}

func (c *nodeGuardServiceClient) GetLiquidityRules(ctx context.Context, in *GetLiquidityRulesRequest, opts ...grpc.CallOption) (*GetLiquidityRulesResponse, error) {
	out := new(GetLiquidityRulesResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetLiquidityRules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetNewWalletAddress(ctx context.Context, in *GetNewWalletAddressRequest, opts ...grpc.CallOption) (*GetNewWalletAddressResponse, error) {
	out := new(GetNewWalletAddressResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetNewWalletAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) RequestWithdrawal(ctx context.Context, in *RequestWithdrawalRequest, opts ...grpc.CallOption) (*RequestWithdrawalResponse, error) {
	out := new(RequestWithdrawalResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/RequestWithdrawal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error) {
	out := new(AddNodeResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/AddNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetNodes(ctx context.Context, in *GetNodesRequest, opts ...grpc.CallOption) (*GetNodesResponse, error) {
	out := new(GetNodesResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetAvailableWallets(ctx context.Context, in *GetAvailableWalletsRequest, opts ...grpc.CallOption) (*GetAvailableWalletsResponse, error) {
	out := new(GetAvailableWalletsResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetAvailableWallets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) OpenChannel(ctx context.Context, in *OpenChannelRequest, opts ...grpc.CallOption) (*OpenChannelResponse, error) {
	out := new(OpenChannelResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/OpenChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) CloseChannel(ctx context.Context, in *CloseChannelRequest, opts ...grpc.CallOption) (*CloseChannelResponse, error) {
	out := new(CloseChannelResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/CloseChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetChannelOperationRequest(ctx context.Context, in *GetChannelOperationRequestRequest, opts ...grpc.CallOption) (*GetChannelOperationRequestResponse, error) {
	out := new(GetChannelOperationRequestResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetChannelOperationRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) AddLiquidityRule(ctx context.Context, in *AddLiquidityRuleRequest, opts ...grpc.CallOption) (*AddLiquidityRuleResponse, error) {
	out := new(AddLiquidityRuleResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/AddLiquidityRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetAvailableUtxos(ctx context.Context, in *GetAvailableUtxosRequest, opts ...grpc.CallOption) (*GetAvailableUtxosResponse, error) {
	out := new(GetAvailableUtxosResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetAvailableUtxos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetWithdrawalsRequestStatus(ctx context.Context, in *GetWithdrawalsRequestStatusRequest, opts ...grpc.CallOption) (*GetWithdrawalsRequestStatusResponse, error) {
	out := new(GetWithdrawalsRequestStatusResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetWithdrawalsRequestStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeGuardServiceClient) GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*GetChannelResponse, error) {
	out := new(GetChannelResponse)
	err := c.cc.Invoke(ctx, "/nodeguard.NodeGuardService/GetChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeGuardServiceServer is the server API for NodeGuardService service.
// All implementations must embed UnimplementedNodeGuardServiceServer
// for forward compatibility
type NodeGuardServiceServer interface {
	// Returns the liquidity rules associated to a node and its channels
	GetLiquidityRules(context.Context, *GetLiquidityRulesRequest) (*GetLiquidityRulesResponse, error)
	// Returns a new unused BTC Address of a given wallet
	GetNewWalletAddress(context.Context, *GetNewWalletAddressRequest) (*GetNewWalletAddressResponse, error)
	// Withdraws funds from a given wallet to a given address
	RequestWithdrawal(context.Context, *RequestWithdrawalRequest) (*RequestWithdrawalResponse, error)
	// Adds a new node to the nodeguard
	AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error)
	// Gets a list of nodes
	GetNodes(context.Context, *GetNodesRequest) (*GetNodesResponse, error)
	// Gets a list of available wallets
	GetAvailableWallets(context.Context, *GetAvailableWalletsRequest) (*GetAvailableWalletsResponse, error)
	// Opens a channel to a given node
	OpenChannel(context.Context, *OpenChannelRequest) (*OpenChannelResponse, error)
	// Closes a channel to a given node
	CloseChannel(context.Context, *CloseChannelRequest) (*CloseChannelResponse, error)
	// Gets a list of channel operations requests by id
	GetChannelOperationRequest(context.Context, *GetChannelOperationRequestRequest) (*GetChannelOperationRequestResponse, error)
	// Adds a liquidity rule to a channel
	AddLiquidityRule(context.Context, *AddLiquidityRuleRequest) (*AddLiquidityRuleResponse, error)
	// Gets a list of available UTXOs for a wallet
	GetAvailableUtxos(context.Context, *GetAvailableUtxosRequest) (*GetAvailableUtxosResponse, error)
	// Gets the status for the provided withdrawals request ids
	GetWithdrawalsRequestStatus(context.Context, *GetWithdrawalsRequestStatusRequest) (*GetWithdrawalsRequestStatusResponse, error)
	// Gets a channel by id
	GetChannel(context.Context, *GetChannelRequest) (*GetChannelResponse, error)
	mustEmbedUnimplementedNodeGuardServiceServer()
}

// UnimplementedNodeGuardServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNodeGuardServiceServer struct {
}

func (UnimplementedNodeGuardServiceServer) GetLiquidityRules(context.Context, *GetLiquidityRulesRequest) (*GetLiquidityRulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLiquidityRules not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetNewWalletAddress(context.Context, *GetNewWalletAddressRequest) (*GetNewWalletAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewWalletAddress not implemented")
}
func (UnimplementedNodeGuardServiceServer) RequestWithdrawal(context.Context, *RequestWithdrawalRequest) (*RequestWithdrawalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestWithdrawal not implemented")
}
func (UnimplementedNodeGuardServiceServer) AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNode not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetNodes(context.Context, *GetNodesRequest) (*GetNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodes not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetAvailableWallets(context.Context, *GetAvailableWalletsRequest) (*GetAvailableWalletsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableWallets not implemented")
}
func (UnimplementedNodeGuardServiceServer) OpenChannel(context.Context, *OpenChannelRequest) (*OpenChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenChannel not implemented")
}
func (UnimplementedNodeGuardServiceServer) CloseChannel(context.Context, *CloseChannelRequest) (*CloseChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseChannel not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetChannelOperationRequest(context.Context, *GetChannelOperationRequestRequest) (*GetChannelOperationRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannelOperationRequest not implemented")
}
func (UnimplementedNodeGuardServiceServer) AddLiquidityRule(context.Context, *AddLiquidityRuleRequest) (*AddLiquidityRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLiquidityRule not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetAvailableUtxos(context.Context, *GetAvailableUtxosRequest) (*GetAvailableUtxosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableUtxos not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetWithdrawalsRequestStatus(context.Context, *GetWithdrawalsRequestStatusRequest) (*GetWithdrawalsRequestStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWithdrawalsRequestStatus not implemented")
}
func (UnimplementedNodeGuardServiceServer) GetChannel(context.Context, *GetChannelRequest) (*GetChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannel not implemented")
}
func (UnimplementedNodeGuardServiceServer) mustEmbedUnimplementedNodeGuardServiceServer() {}

// UnsafeNodeGuardServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeGuardServiceServer will
// result in compilation errors.
type UnsafeNodeGuardServiceServer interface {
	mustEmbedUnimplementedNodeGuardServiceServer()
}

func RegisterNodeGuardServiceServer(s grpc.ServiceRegistrar, srv NodeGuardServiceServer) {
	s.RegisterService(&NodeGuardService_ServiceDesc, srv)
}

func _NodeGuardService_GetLiquidityRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLiquidityRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetLiquidityRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetLiquidityRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetLiquidityRules(ctx, req.(*GetLiquidityRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetNewWalletAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNewWalletAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetNewWalletAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetNewWalletAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetNewWalletAddress(ctx, req.(*GetNewWalletAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_RequestWithdrawal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithdrawalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).RequestWithdrawal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/RequestWithdrawal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).RequestWithdrawal(ctx, req.(*RequestWithdrawalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/AddNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).AddNode(ctx, req.(*AddNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetNodes(ctx, req.(*GetNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetAvailableWallets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableWalletsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetAvailableWallets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetAvailableWallets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetAvailableWallets(ctx, req.(*GetAvailableWalletsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_OpenChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).OpenChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/OpenChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).OpenChannel(ctx, req.(*OpenChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_CloseChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).CloseChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/CloseChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).CloseChannel(ctx, req.(*CloseChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetChannelOperationRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelOperationRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetChannelOperationRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetChannelOperationRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetChannelOperationRequest(ctx, req.(*GetChannelOperationRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_AddLiquidityRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLiquidityRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).AddLiquidityRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/AddLiquidityRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).AddLiquidityRule(ctx, req.(*AddLiquidityRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetAvailableUtxos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableUtxosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetAvailableUtxos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetAvailableUtxos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetAvailableUtxos(ctx, req.(*GetAvailableUtxosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetWithdrawalsRequestStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWithdrawalsRequestStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetWithdrawalsRequestStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetWithdrawalsRequestStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetWithdrawalsRequestStatus(ctx, req.(*GetWithdrawalsRequestStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeGuardService_GetChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeGuardServiceServer).GetChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeguard.NodeGuardService/GetChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeGuardServiceServer).GetChannel(ctx, req.(*GetChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeGuardService_ServiceDesc is the grpc.ServiceDesc for NodeGuardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeGuardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nodeguard.NodeGuardService",
	HandlerType: (*NodeGuardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLiquidityRules",
			Handler:    _NodeGuardService_GetLiquidityRules_Handler,
		},
		{
			MethodName: "GetNewWalletAddress",
			Handler:    _NodeGuardService_GetNewWalletAddress_Handler,
		},
		{
			MethodName: "RequestWithdrawal",
			Handler:    _NodeGuardService_RequestWithdrawal_Handler,
		},
		{
			MethodName: "AddNode",
			Handler:    _NodeGuardService_AddNode_Handler,
		},
		{
			MethodName: "GetNodes",
			Handler:    _NodeGuardService_GetNodes_Handler,
		},
		{
			MethodName: "GetAvailableWallets",
			Handler:    _NodeGuardService_GetAvailableWallets_Handler,
		},
		{
			MethodName: "OpenChannel",
			Handler:    _NodeGuardService_OpenChannel_Handler,
		},
		{
			MethodName: "CloseChannel",
			Handler:    _NodeGuardService_CloseChannel_Handler,
		},
		{
			MethodName: "GetChannelOperationRequest",
			Handler:    _NodeGuardService_GetChannelOperationRequest_Handler,
		},
		{
			MethodName: "AddLiquidityRule",
			Handler:    _NodeGuardService_AddLiquidityRule_Handler,
		},
		{
			MethodName: "GetAvailableUtxos",
			Handler:    _NodeGuardService_GetAvailableUtxos_Handler,
		},
		{
			MethodName: "GetWithdrawalsRequestStatus",
			Handler:    _NodeGuardService_GetWithdrawalsRequestStatus_Handler,
		},
		{
			MethodName: "GetChannel",
			Handler:    _NodeGuardService_GetChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nodeguard.proto",
}
