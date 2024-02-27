// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: internal/transport/grpc/protofiles/credentials_service.proto

package protofiles

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

const (
	CredSvc_CredAdd_FullMethodName  = "/CredSvc/CredAdd"
	CredSvc_CredGet_FullMethodName  = "/CredSvc/CredGet"
	CredSvc_CredUpd_FullMethodName  = "/CredSvc/CredUpd"
	CredSvc_CredDel_FullMethodName  = "/CredSvc/CredDel"
	CredSvc_CredList_FullMethodName = "/CredSvc/CredList"
)

// CredSvcClient is the client API for CredSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CredSvcClient interface {
	CredAdd(ctx context.Context, in *CredAddRequest, opts ...grpc.CallOption) (*CredAddResponse, error)
	CredGet(ctx context.Context, in *CredGetRequest, opts ...grpc.CallOption) (*CredGetResponse, error)
	CredUpd(ctx context.Context, in *CredUpdRequest, opts ...grpc.CallOption) (*CredUpdResponse, error)
	CredDel(ctx context.Context, in *CredDelRequest, opts ...grpc.CallOption) (*CredDelResponse, error)
	CredList(ctx context.Context, in *CredListRequest, opts ...grpc.CallOption) (*CredListResponse, error)
}

type credSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewCredSvcClient(cc grpc.ClientConnInterface) CredSvcClient {
	return &credSvcClient{cc}
}

func (c *credSvcClient) CredAdd(ctx context.Context, in *CredAddRequest, opts ...grpc.CallOption) (*CredAddResponse, error) {
	out := new(CredAddResponse)
	err := c.cc.Invoke(ctx, CredSvc_CredAdd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credSvcClient) CredGet(ctx context.Context, in *CredGetRequest, opts ...grpc.CallOption) (*CredGetResponse, error) {
	out := new(CredGetResponse)
	err := c.cc.Invoke(ctx, CredSvc_CredGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credSvcClient) CredUpd(ctx context.Context, in *CredUpdRequest, opts ...grpc.CallOption) (*CredUpdResponse, error) {
	out := new(CredUpdResponse)
	err := c.cc.Invoke(ctx, CredSvc_CredUpd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credSvcClient) CredDel(ctx context.Context, in *CredDelRequest, opts ...grpc.CallOption) (*CredDelResponse, error) {
	out := new(CredDelResponse)
	err := c.cc.Invoke(ctx, CredSvc_CredDel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credSvcClient) CredList(ctx context.Context, in *CredListRequest, opts ...grpc.CallOption) (*CredListResponse, error) {
	out := new(CredListResponse)
	err := c.cc.Invoke(ctx, CredSvc_CredList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CredSvcServer is the server API for CredSvc service.
// All implementations must embed UnimplementedCredSvcServer
// for forward compatibility
type CredSvcServer interface {
	CredAdd(context.Context, *CredAddRequest) (*CredAddResponse, error)
	CredGet(context.Context, *CredGetRequest) (*CredGetResponse, error)
	CredUpd(context.Context, *CredUpdRequest) (*CredUpdResponse, error)
	CredDel(context.Context, *CredDelRequest) (*CredDelResponse, error)
	CredList(context.Context, *CredListRequest) (*CredListResponse, error)
	mustEmbedUnimplementedCredSvcServer()
}

// UnimplementedCredSvcServer must be embedded to have forward compatible implementations.
type UnimplementedCredSvcServer struct {
}

func (UnimplementedCredSvcServer) CredAdd(context.Context, *CredAddRequest) (*CredAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CredAdd not implemented")
}
func (UnimplementedCredSvcServer) CredGet(context.Context, *CredGetRequest) (*CredGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CredGet not implemented")
}
func (UnimplementedCredSvcServer) CredUpd(context.Context, *CredUpdRequest) (*CredUpdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CredUpd not implemented")
}
func (UnimplementedCredSvcServer) CredDel(context.Context, *CredDelRequest) (*CredDelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CredDel not implemented")
}
func (UnimplementedCredSvcServer) CredList(context.Context, *CredListRequest) (*CredListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CredList not implemented")
}
func (UnimplementedCredSvcServer) mustEmbedUnimplementedCredSvcServer() {}

// UnsafeCredSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CredSvcServer will
// result in compilation errors.
type UnsafeCredSvcServer interface {
	mustEmbedUnimplementedCredSvcServer()
}

func RegisterCredSvcServer(s grpc.ServiceRegistrar, srv CredSvcServer) {
	s.RegisterService(&CredSvc_ServiceDesc, srv)
}

func _CredSvc_CredAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredSvcServer).CredAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CredSvc_CredAdd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredSvcServer).CredAdd(ctx, req.(*CredAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredSvc_CredGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredSvcServer).CredGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CredSvc_CredGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredSvcServer).CredGet(ctx, req.(*CredGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredSvc_CredUpd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredUpdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredSvcServer).CredUpd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CredSvc_CredUpd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredSvcServer).CredUpd(ctx, req.(*CredUpdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredSvc_CredDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredDelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredSvcServer).CredDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CredSvc_CredDel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredSvcServer).CredDel(ctx, req.(*CredDelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredSvc_CredList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredSvcServer).CredList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CredSvc_CredList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredSvcServer).CredList(ctx, req.(*CredListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CredSvc_ServiceDesc is the grpc.ServiceDesc for CredSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CredSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CredSvc",
	HandlerType: (*CredSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CredAdd",
			Handler:    _CredSvc_CredAdd_Handler,
		},
		{
			MethodName: "CredGet",
			Handler:    _CredSvc_CredGet_Handler,
		},
		{
			MethodName: "CredUpd",
			Handler:    _CredSvc_CredUpd_Handler,
		},
		{
			MethodName: "CredDel",
			Handler:    _CredSvc_CredDel_Handler,
		},
		{
			MethodName: "CredList",
			Handler:    _CredSvc_CredList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/transport/grpc/protofiles/credentials_service.proto",
}
