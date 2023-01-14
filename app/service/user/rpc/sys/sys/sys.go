// Code generated by goctl. DO NOT EDIT!
// Source: user.sys.proto

package sys

import (
	"context"

	"douyin/app/service/user/rpc/sys/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginReq    = pb.LoginReq
	LoginRes    = pb.LoginRes
	RegisterReq = pb.RegisterReq
	RegisterRes = pb.RegisterRes

	Sys interface {
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error)
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error)
	}

	defaultSys struct {
		cli zrpc.Client
	}
)

func NewSys(cli zrpc.Client) Sys {
	return &defaultSys{
		cli: cli,
	}
}

func (m *defaultSys) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error) {
	client := pb.NewSysClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultSys) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error) {
	client := pb.NewSysClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}
