// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: internal/transport/grpc/protofiles/credentials_service.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BlobSvc_BlobAdd_FullMethodName  = "/BlobSvc/BlobAdd"
	BlobSvc_BlobGet_FullMethodName  = "/BlobSvc/BlobGet"
	BlobSvc_BlobUpd_FullMethodName  = "/BlobSvc/BlobUpd"
	BlobSvc_BlobDel_FullMethodName  = "/BlobSvc/BlobDel"
	BlobSvc_BlobList_FullMethodName = "/BlobSvc/BlobList"
)

// BlobSvcClient is the client API for BlobSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlobSvcClient interface {
	BlobAdd(ctx context.Context, in *BlobAddRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	BlobGet(ctx context.Context, in *BlobGetRequest, opts ...grpc.CallOption) (*BlobGetResponse, error)
	BlobUpd(ctx context.Context, in *BlobUpdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	BlobDel(ctx context.Context, in *BlobDelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	BlobList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BlobListResponse, error)
}

type blobSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewBlobSvcClient(cc grpc.ClientConnInterface) BlobSvcClient {
	return &blobSvcClient{cc}
}

func (c *blobSvcClient) BlobAdd(ctx context.Context, in *BlobAddRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlobSvc_BlobAdd_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobSvcClient) BlobGet(ctx context.Context, in *BlobGetRequest, opts ...grpc.CallOption) (*BlobGetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BlobGetResponse)
	err := c.cc.Invoke(ctx, BlobSvc_BlobGet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobSvcClient) BlobUpd(ctx context.Context, in *BlobUpdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlobSvc_BlobUpd_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobSvcClient) BlobDel(ctx context.Context, in *BlobDelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlobSvc_BlobDel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobSvcClient) BlobList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BlobListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BlobListResponse)
	err := c.cc.Invoke(ctx, BlobSvc_BlobList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlobSvcServer is the server API for BlobSvc service.
// All implementations must embed UnimplementedBlobSvcServer
// for forward compatibility.
type BlobSvcServer interface {
	BlobAdd(context.Context, *BlobAddRequest) (*emptypb.Empty, error)
	BlobGet(context.Context, *BlobGetRequest) (*BlobGetResponse, error)
	BlobUpd(context.Context, *BlobUpdRequest) (*emptypb.Empty, error)
	BlobDel(context.Context, *BlobDelRequest) (*emptypb.Empty, error)
	BlobList(context.Context, *emptypb.Empty) (*BlobListResponse, error)
	mustEmbedUnimplementedBlobSvcServer()
}

// UnimplementedBlobSvcServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBlobSvcServer struct{}

func (UnimplementedBlobSvcServer) BlobAdd(context.Context, *BlobAddRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlobAdd not implemented")
}
func (UnimplementedBlobSvcServer) BlobGet(context.Context, *BlobGetRequest) (*BlobGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlobGet not implemented")
}
func (UnimplementedBlobSvcServer) BlobUpd(context.Context, *BlobUpdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlobUpd not implemented")
}
func (UnimplementedBlobSvcServer) BlobDel(context.Context, *BlobDelRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlobDel not implemented")
}
func (UnimplementedBlobSvcServer) BlobList(context.Context, *emptypb.Empty) (*BlobListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlobList not implemented")
}
func (UnimplementedBlobSvcServer) mustEmbedUnimplementedBlobSvcServer() {}
func (UnimplementedBlobSvcServer) testEmbeddedByValue()                 {}

// UnsafeBlobSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlobSvcServer will
// result in compilation myerrors.
type UnsafeBlobSvcServer interface {
	mustEmbedUnimplementedBlobSvcServer()
}

func RegisterBlobSvcServer(s grpc.ServiceRegistrar, srv BlobSvcServer) {
	// If the following call pancis, it indicates UnimplementedBlobSvcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BlobSvc_ServiceDesc, srv)
}

func _BlobSvc_BlobAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobSvcServer).BlobAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlobSvc_BlobAdd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobSvcServer).BlobAdd(ctx, req.(*BlobAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlobSvc_BlobGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobSvcServer).BlobGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlobSvc_BlobGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobSvcServer).BlobGet(ctx, req.(*BlobGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlobSvc_BlobUpd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobUpdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobSvcServer).BlobUpd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlobSvc_BlobUpd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobSvcServer).BlobUpd(ctx, req.(*BlobUpdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlobSvc_BlobDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobDelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobSvcServer).BlobDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlobSvc_BlobDel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobSvcServer).BlobDel(ctx, req.(*BlobDelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlobSvc_BlobList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobSvcServer).BlobList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlobSvc_BlobList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobSvcServer).BlobList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BlobSvc_ServiceDesc is the grpc.ServiceDesc for BlobSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlobSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BlobSvc",
	HandlerType: (*BlobSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BlobAdd",
			Handler:    _BlobSvc_BlobAdd_Handler,
		},
		{
			MethodName: "BlobGet",
			Handler:    _BlobSvc_BlobGet_Handler,
		},
		{
			MethodName: "BlobUpd",
			Handler:    _BlobSvc_BlobUpd_Handler,
		},
		{
			MethodName: "BlobDel",
			Handler:    _BlobSvc_BlobDel_Handler,
		},
		{
			MethodName: "BlobList",
			Handler:    _BlobSvc_BlobList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/transport/grpc/protofiles/credentials_service.proto",
}
