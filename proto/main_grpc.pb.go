// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/main.proto

package proto

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

// WayServiceClient is the client API for WayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WayServiceClient interface {
	GetLines(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Lines, error)
}

type wayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWayServiceClient(cc grpc.ClientConnInterface) WayServiceClient {
	return &wayServiceClient{cc}
}

func (c *wayServiceClient) GetLines(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Lines, error) {
	out := new(Lines)
	err := c.cc.Invoke(ctx, "/main.WayService/GetLines", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WayServiceServer is the server API for WayService service.
// All implementations must embed UnimplementedWayServiceServer
// for forward compatibility
type WayServiceServer interface {
	GetLines(context.Context, *Params) (*Lines, error)
	mustEmbedUnimplementedWayServiceServer()
}

// UnimplementedWayServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWayServiceServer struct {
}

func (UnimplementedWayServiceServer) GetLines(context.Context, *Params) (*Lines, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLines not implemented")
}
func (UnimplementedWayServiceServer) mustEmbedUnimplementedWayServiceServer() {}

// UnsafeWayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WayServiceServer will
// result in compilation errors.
type UnsafeWayServiceServer interface {
	mustEmbedUnimplementedWayServiceServer()
}

func RegisterWayServiceServer(s grpc.ServiceRegistrar, srv WayServiceServer) {
	s.RegisterService(&WayService_ServiceDesc, srv)
}

func _WayService_GetLines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Params)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WayServiceServer).GetLines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.WayService/GetLines",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WayServiceServer).GetLines(ctx, req.(*Params))
	}
	return interceptor(ctx, in, info, handler)
}

// WayService_ServiceDesc is the grpc.ServiceDesc for WayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.WayService",
	HandlerType: (*WayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLines",
			Handler:    _WayService_GetLines_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/main.proto",
}