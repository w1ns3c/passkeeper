// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: internal/transport/grpc/protofiles/sessionkey.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SessionKeySvc_KeyExchange_FullMethodName = "/SessionKeySvc/KeyExchange"
)

// SessionKeySvcClient is the client API for SessionKeySvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionKeySvcClient interface {
	KeyExchange(ctx context.Context, in *SessionKeyReq, opts ...grpc.CallOption) (*SessionKeyResp, error)
}

type sessionKeySvcClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionKeySvcClient(cc grpc.ClientConnInterface) SessionKeySvcClient {
	return &sessionKeySvcClient{cc}
}

func (c *sessionKeySvcClient) KeyExchange(ctx context.Context, in *SessionKeyReq, opts ...grpc.CallOption) (*SessionKeyResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SessionKeyResp)
	err := c.cc.Invoke(ctx, SessionKeySvc_KeyExchange_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionKeySvcServer is the server API for SessionKeySvc service.
// All implementations must embed UnimplementedSessionKeySvcServer
// for forward compatibility.
type SessionKeySvcServer interface {
	KeyExchange(context.Context, *SessionKeyReq) (*SessionKeyResp, error)
	mustEmbedUnimplementedSessionKeySvcServer()
}

// UnimplementedSessionKeySvcServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSessionKeySvcServer struct{}

func (UnimplementedSessionKeySvcServer) KeyExchange(context.Context, *SessionKeyReq) (*SessionKeyResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeyExchange not implemented")
}
func (UnimplementedSessionKeySvcServer) mustEmbedUnimplementedSessionKeySvcServer() {}
func (UnimplementedSessionKeySvcServer) testEmbeddedByValue()                       {}

// UnsafeSessionKeySvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionKeySvcServer will
// result in compilation myerrors.
type UnsafeSessionKeySvcServer interface {
	mustEmbedUnimplementedSessionKeySvcServer()
}

func RegisterSessionKeySvcServer(s grpc.ServiceRegistrar, srv SessionKeySvcServer) {
	// If the following call pancis, it indicates UnimplementedSessionKeySvcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SessionKeySvc_ServiceDesc, srv)
}

func _SessionKeySvc_KeyExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionKeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionKeySvcServer).KeyExchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SessionKeySvc_KeyExchange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionKeySvcServer).KeyExchange(ctx, req.(*SessionKeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SessionKeySvc_ServiceDesc is the grpc.ServiceDesc for SessionKeySvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SessionKeySvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SessionKeySvc",
	HandlerType: (*SessionKeySvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "KeyExchange",
			Handler:    _SessionKeySvc_KeyExchange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/transport/grpc/protofiles/sessionkey.proto",
}
