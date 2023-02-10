// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: looprpc/debug.proto

package looprpc

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

// DebugClient is the client API for Debug service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DebugClient interface {
	// ForceAutoLoop is intended for *testing purposes only* and will not work on
	// mainnet. This endpoint ticks our autoloop timer, triggering automated
	// dispatch of a swap if one is suggested.
	ForceAutoLoop(ctx context.Context, in *ForceAutoLoopRequest, opts ...grpc.CallOption) (*ForceAutoLoopResponse, error)
}

type debugClient struct {
	cc grpc.ClientConnInterface
}

func NewDebugClient(cc grpc.ClientConnInterface) DebugClient {
	return &debugClient{cc}
}

func (c *debugClient) ForceAutoLoop(ctx context.Context, in *ForceAutoLoopRequest, opts ...grpc.CallOption) (*ForceAutoLoopResponse, error) {
	out := new(ForceAutoLoopResponse)
	err := c.cc.Invoke(ctx, "/looprpc.Debug/ForceAutoLoop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebugServer is the server API for Debug service.
// All implementations must embed UnimplementedDebugServer
// for forward compatibility
type DebugServer interface {
	// ForceAutoLoop is intended for *testing purposes only* and will not work on
	// mainnet. This endpoint ticks our autoloop timer, triggering automated
	// dispatch of a swap if one is suggested.
	ForceAutoLoop(context.Context, *ForceAutoLoopRequest) (*ForceAutoLoopResponse, error)
	mustEmbedUnimplementedDebugServer()
}

// UnimplementedDebugServer must be embedded to have forward compatible implementations.
type UnimplementedDebugServer struct {
}

func (UnimplementedDebugServer) ForceAutoLoop(context.Context, *ForceAutoLoopRequest) (*ForceAutoLoopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForceAutoLoop not implemented")
}
func (UnimplementedDebugServer) mustEmbedUnimplementedDebugServer() {}

// UnsafeDebugServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DebugServer will
// result in compilation errors.
type UnsafeDebugServer interface {
	mustEmbedUnimplementedDebugServer()
}

func RegisterDebugServer(s grpc.ServiceRegistrar, srv DebugServer) {
	s.RegisterService(&Debug_ServiceDesc, srv)
}

func _Debug_ForceAutoLoop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForceAutoLoopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServer).ForceAutoLoop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.Debug/ForceAutoLoop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServer).ForceAutoLoop(ctx, req.(*ForceAutoLoopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Debug_ServiceDesc is the grpc.ServiceDesc for Debug service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Debug_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "looprpc.Debug",
	HandlerType: (*DebugServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ForceAutoLoop",
			Handler:    _Debug_ForceAutoLoop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "looprpc/debug.proto",
}
