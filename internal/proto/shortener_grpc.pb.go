// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: internal/proto/shortener.proto

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

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error)
	AddBatchURL(ctx context.Context, in *AddBatchURLRequest, opts ...grpc.CallOption) (*AddBatchURLResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
}

// autogenerated protobuf
type shortenerClient struct {
	cc grpc.ClientConnInterface
}

// autogenerated protobuf
func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

// autogenerated protobuf
func (c *shortenerClient) AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error) {
	out := new(AddURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/AddURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// autogenerated protobuf
func (c *shortenerClient) AddBatchURL(ctx context.Context, in *AddBatchURLRequest, opts ...grpc.CallOption) (*AddBatchURLResponse, error) {
	out := new(AddBatchURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/AddBatchURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// autogenerated protobuf
func (c *shortenerClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/GetURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// autogenerated protobuf
func (c *shortenerClient) Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/shortener.Shortener/Stats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error)
	AddBatchURL(context.Context, *AddBatchURLRequest) (*AddBatchURLResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	Stats(context.Context, *StatsRequest) (*StatsResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

// autogenerated protobuf
func (UnimplementedShortenerServer) AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddURL not implemented")
}

// autogenerated protobuf
func (UnimplementedShortenerServer) AddBatchURL(context.Context, *AddBatchURLRequest) (*AddBatchURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBatchURL not implemented")
}

// autogenerated protobuf
func (UnimplementedShortenerServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}

// autogenerated protobuf
func (UnimplementedShortenerServer) Stats(context.Context, *StatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}

// autogenerated protobuf
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

// autogenerated protobuf
func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

// autogenerated protobuf
func _Shortener_AddURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).AddURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/AddURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).AddURL(ctx, req.(*AddURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// autogenerated protobuf
func _Shortener_AddBatchURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBatchURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).AddBatchURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/AddBatchURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).AddBatchURL(ctx, req.(*AddBatchURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// autogenerated protobuf
func _Shortener_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/GetURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// autogenerated protobuf
func _Shortener_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.Shortener/Stats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Stats(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddURL",
			Handler:    _Shortener_AddURL_Handler,
		},
		{
			MethodName: "AddBatchURL",
			Handler:    _Shortener_AddBatchURL_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _Shortener_GetURL_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Shortener_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/shortener.proto",
}
