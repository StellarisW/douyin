package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/internal/model/consts"
	"douyin/app/service/video/rpc/sys/internal/model/crud"

	"douyin/app/service/video/rpc/sys/internal/svc"
	"douyin/app/service/video/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentLogic) Comment(in *pb.CommentReq) (*pb.CommentRes, error) {
	erx := l.svcCtx.CrudModel.Comment(l.ctx, in.UserId, in.VideoId, in.ActionType, in.CommentText, in.CommentId)
	if erx != nil {
		if erx.Code() == crud.ErrIdInvalidActionType {
			return &pb.CommentRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicCrud,
					crud.ErrIdOprComment,
					crud.ErrIdInvalidActionType,
				),
				StatusMsg: crud.ErrInvalidContentType,
			}, nil
		}
		return &pb.CommentRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.CommentRes{
		StatusCode: 0,
		StatusMsg:  "comment successfully",
	}, nil
}
