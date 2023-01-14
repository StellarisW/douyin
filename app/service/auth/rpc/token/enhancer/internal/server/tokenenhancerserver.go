// Code generated by goctl. DO NOT EDIT!
// Source: auth.token_enhancer.proto

package server

import (
	"context"
	"douyin/app/common/log"
	"go.uber.org/zap"

	"douyin/app/service/auth/rpc/token/enhancer/internal/logic"
	"douyin/app/service/auth/rpc/token/enhancer/internal/svc"
	"douyin/app/service/auth/rpc/token/enhancer/pb"
)

type TokenEnhancerServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedTokenEnhancerServer
}

func NewTokenEnhancerServer(svcCtx *svc.ServiceContext) *TokenEnhancerServer {
	return &TokenEnhancerServer{
		svcCtx: svcCtx,
	}
}

func (s *TokenEnhancerServer) GenerateToken(ctx context.Context, in *pb.GenerateTokenReq) (*pb.GenerateTokenRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewGenerateTokenLogic(ctx, s.svcCtx)
	res, err := l.GenerateToken(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *TokenEnhancerServer) ReadToken(ctx context.Context, in *pb.ReadTokenReq) (*pb.ReadTokenRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewReadTokenLogic(ctx, s.svcCtx)
	res, err := l.ReadToken(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}