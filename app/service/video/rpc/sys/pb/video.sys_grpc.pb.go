// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: video.sys.proto

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

// SysClient is the client API for Sys service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SysClient interface {
	Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishRes, error)
	GetPublishList(ctx context.Context, in *GetPublishListReq, opts ...grpc.CallOption) (*GetPublishListRes, error)
	Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedRes, error)
	Favorite(ctx context.Context, in *FavoriteReq, opts ...grpc.CallOption) (*FavoriteRes, error)
	GetFavoriteList(ctx context.Context, in *GetFavoriteListReq, opts ...grpc.CallOption) (*GetFavoriteListRes, error)
	Comment(ctx context.Context, in *CommentReq, opts ...grpc.CallOption) (*CommentRes, error)
	GetCommentList(ctx context.Context, in *GetCommentListReq, opts ...grpc.CallOption) (*GetCommentListRes, error)
}

type sysClient struct {
	cc grpc.ClientConnInterface
}

func NewSysClient(cc grpc.ClientConnInterface) SysClient {
	return &sysClient{cc}
}

func (c *sysClient) Publish(ctx context.Context, in *PublishReq, opts ...grpc.CallOption) (*PublishRes, error) {
	out := new(PublishRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) GetPublishList(ctx context.Context, in *GetPublishListReq, opts ...grpc.CallOption) (*GetPublishListRes, error) {
	out := new(GetPublishListRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/GetPublishList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) Feed(ctx context.Context, in *FeedReq, opts ...grpc.CallOption) (*FeedRes, error) {
	out := new(FeedRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/Feed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) Favorite(ctx context.Context, in *FavoriteReq, opts ...grpc.CallOption) (*FavoriteRes, error) {
	out := new(FavoriteRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/Favorite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) GetFavoriteList(ctx context.Context, in *GetFavoriteListReq, opts ...grpc.CallOption) (*GetFavoriteListRes, error) {
	out := new(GetFavoriteListRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/GetFavoriteList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) Comment(ctx context.Context, in *CommentReq, opts ...grpc.CallOption) (*CommentRes, error) {
	out := new(CommentRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/Comment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sysClient) GetCommentList(ctx context.Context, in *GetCommentListReq, opts ...grpc.CallOption) (*GetCommentListRes, error) {
	out := new(GetCommentListRes)
	err := c.cc.Invoke(ctx, "/video.sys.Sys/GetCommentList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SysServer is the server API for Sys service.
// All implementations must embed UnimplementedSysServer
// for forward compatibility
type SysServer interface {
	Publish(context.Context, *PublishReq) (*PublishRes, error)
	GetPublishList(context.Context, *GetPublishListReq) (*GetPublishListRes, error)
	Feed(context.Context, *FeedReq) (*FeedRes, error)
	Favorite(context.Context, *FavoriteReq) (*FavoriteRes, error)
	GetFavoriteList(context.Context, *GetFavoriteListReq) (*GetFavoriteListRes, error)
	Comment(context.Context, *CommentReq) (*CommentRes, error)
	GetCommentList(context.Context, *GetCommentListReq) (*GetCommentListRes, error)
	mustEmbedUnimplementedSysServer()
}

// UnimplementedSysServer must be embedded to have forward compatible implementations.
type UnimplementedSysServer struct {
}

func (UnimplementedSysServer) Publish(context.Context, *PublishReq) (*PublishRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedSysServer) GetPublishList(context.Context, *GetPublishListReq) (*GetPublishListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublishList not implemented")
}
func (UnimplementedSysServer) Feed(context.Context, *FeedReq) (*FeedRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Feed not implemented")
}
func (UnimplementedSysServer) Favorite(context.Context, *FavoriteReq) (*FavoriteRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Favorite not implemented")
}
func (UnimplementedSysServer) GetFavoriteList(context.Context, *GetFavoriteListReq) (*GetFavoriteListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavoriteList not implemented")
}
func (UnimplementedSysServer) Comment(context.Context, *CommentReq) (*CommentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Comment not implemented")
}
func (UnimplementedSysServer) GetCommentList(context.Context, *GetCommentListReq) (*GetCommentListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentList not implemented")
}
func (UnimplementedSysServer) mustEmbedUnimplementedSysServer() {}

// UnsafeSysServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SysServer will
// result in compilation errors.
type UnsafeSysServer interface {
	mustEmbedUnimplementedSysServer()
}

func RegisterSysServer(s grpc.ServiceRegistrar, srv SysServer) {
	s.RegisterService(&Sys_ServiceDesc, srv)
}

func _Sys_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).Publish(ctx, req.(*PublishReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_GetPublishList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublishListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).GetPublishList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/GetPublishList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).GetPublishList(ctx, req.(*GetPublishListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_Feed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).Feed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/Feed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).Feed(ctx, req.(*FeedReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_Favorite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FavoriteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).Favorite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/Favorite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).Favorite(ctx, req.(*FavoriteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_GetFavoriteList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFavoriteListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).GetFavoriteList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/GetFavoriteList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).GetFavoriteList(ctx, req.(*GetFavoriteListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_Comment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).Comment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/Comment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).Comment(ctx, req.(*CommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sys_GetCommentList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SysServer).GetCommentList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.sys.Sys/GetCommentList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SysServer).GetCommentList(ctx, req.(*GetCommentListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Sys_ServiceDesc is the grpc.ServiceDesc for Sys service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sys_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "video.sys.Sys",
	HandlerType: (*SysServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Sys_Publish_Handler,
		},
		{
			MethodName: "GetPublishList",
			Handler:    _Sys_GetPublishList_Handler,
		},
		{
			MethodName: "Feed",
			Handler:    _Sys_Feed_Handler,
		},
		{
			MethodName: "Favorite",
			Handler:    _Sys_Favorite_Handler,
		},
		{
			MethodName: "GetFavoriteList",
			Handler:    _Sys_GetFavoriteList_Handler,
		},
		{
			MethodName: "Comment",
			Handler:    _Sys_Comment_Handler,
		},
		{
			MethodName: "GetCommentList",
			Handler:    _Sys_GetCommentList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "video.sys.proto",
}
