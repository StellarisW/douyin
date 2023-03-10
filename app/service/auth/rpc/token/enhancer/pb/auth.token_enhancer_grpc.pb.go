// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: auth.token_enhancer.proto

package pb

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

// TokenEnhancerClient is the client API for TokenEnhancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenEnhancerClient interface {
	GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenRes, error)
	ReadToken(ctx context.Context, in *ReadTokenReq, opts ...grpc.CallOption) (*ReadTokenRes, error)
}

type tokenEnhancerClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenEnhancerClient(cc grpc.ClientConnInterface) TokenEnhancerClient {
	return &tokenEnhancerClient{cc}
}

func (c *tokenEnhancerClient) GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenRes, error) {
	out := new(GenerateTokenRes)
	err := c.cc.Invoke(ctx, "/auth.token_enhancer.TokenEnhancer/GenerateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenEnhancerClient) ReadToken(ctx context.Context, in *ReadTokenReq, opts ...grpc.CallOption) (*ReadTokenRes, error) {
	out := new(ReadTokenRes)
	err := c.cc.Invoke(ctx, "/auth.token_enhancer.TokenEnhancer/ReadToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenEnhancerServer is the server API for TokenEnhancer service.
// All implementations must embed UnimplementedTokenEnhancerServer
// for forward compatibility
type TokenEnhancerServer interface {
	GenerateToken(context.Context, *GenerateTokenReq) (*GenerateTokenRes, error)
	ReadToken(context.Context, *ReadTokenReq) (*ReadTokenRes, error)
	mustEmbedUnimplementedTokenEnhancerServer()
}

// UnimplementedTokenEnhancerServer must be embedded to have forward compatible implementations.
type UnimplementedTokenEnhancerServer struct {
}

func (UnimplementedTokenEnhancerServer) GenerateToken(context.Context, *GenerateTokenReq) (*GenerateTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateToken not implemented")
}
func (UnimplementedTokenEnhancerServer) ReadToken(context.Context, *ReadTokenReq) (*ReadTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadToken not implemented")
}
func (UnimplementedTokenEnhancerServer) mustEmbedUnimplementedTokenEnhancerServer() {}

// UnsafeTokenEnhancerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenEnhancerServer will
// result in compilation errors.
type UnsafeTokenEnhancerServer interface {
	mustEmbedUnimplementedTokenEnhancerServer()
}

func RegisterTokenEnhancerServer(s grpc.ServiceRegistrar, srv TokenEnhancerServer) {
	s.RegisterService(&TokenEnhancer_ServiceDesc, srv)
}

func _TokenEnhancer_GenerateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenEnhancerServer).GenerateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.token_enhancer.TokenEnhancer/GenerateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenEnhancerServer).GenerateToken(ctx, req.(*GenerateTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenEnhancer_ReadToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenEnhancerServer).ReadToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.token_enhancer.TokenEnhancer/ReadToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenEnhancerServer).ReadToken(ctx, req.(*ReadTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenEnhancer_ServiceDesc is the grpc.ServiceDesc for TokenEnhancer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenEnhancer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.token_enhancer.TokenEnhancer",
	HandlerType: (*TokenEnhancerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateToken",
			Handler:    _TokenEnhancer_GenerateToken_Handler,
		},
		{
			MethodName: "ReadToken",
			Handler:    _TokenEnhancer_ReadToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.token_enhancer.proto",
}
