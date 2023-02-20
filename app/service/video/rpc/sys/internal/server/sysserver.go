// Code generated by goctl. DO NOT EDIT!
// Source: video.sys.proto

package server

import (
	"context"
	"douyin/app/common/log"
	"go.uber.org/zap"

	"douyin/app/service/video/rpc/sys/internal/logic"
	"douyin/app/service/video/rpc/sys/internal/svc"
	"douyin/app/service/video/rpc/sys/pb"
)

type SysServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedSysServer
}

func NewSysServer(svcCtx *svc.ServiceContext) *SysServer {
	return &SysServer{
		svcCtx: svcCtx,
	}
}

func (s *SysServer) Publish(ctx context.Context, in *pb.PublishReq) (*pb.PublishRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewPublishLogic(ctx, s.svcCtx)
	res, err := l.Publish(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) GetPublishList(ctx context.Context, in *pb.GetPublishListReq) (*pb.GetPublishListRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewGetPublishListLogic(ctx, s.svcCtx)
	res, err := l.GetPublishList(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) Feed(ctx context.Context, in *pb.FeedReq) (*pb.FeedRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewFeedLogic(ctx, s.svcCtx)
	res, err := l.Feed(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) Favorite(ctx context.Context, in *pb.FavoriteReq) (*pb.FavoriteRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewFavoriteLogic(ctx, s.svcCtx)
	res, err := l.Favorite(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) GetFavoriteList(ctx context.Context, in *pb.GetFavoriteListReq) (*pb.GetFavoriteListRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewGetFavoriteListLogic(ctx, s.svcCtx)
	res, err := l.GetFavoriteList(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) Comment(ctx context.Context, in *pb.CommentReq) (*pb.CommentRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewCommentLogic(ctx, s.svcCtx)
	res, err := l.Comment(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) ManageComment(ctx context.Context, in *pb.ManageCommentReq) (*pb.ManageCommentRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewManageCommentLogic(ctx, s.svcCtx)
	res, err := l.ManageComment(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}

func (s *SysServer) GetCommentList(ctx context.Context, in *pb.GetCommentListReq) (*pb.GetCommentListRes, error) {
	log.Logger.Debug("recv:", zap.String("msg", in.String()))
	l := logic.NewGetCommentListLogic(ctx, s.svcCtx)
	res, err := l.GetCommentList(in)
	log.Logger.Debug("send:", zap.String("msg", res.String()))
	return res, err
}
