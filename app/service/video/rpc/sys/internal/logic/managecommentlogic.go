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

type ManageCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageCommentLogic {
	return &ManageCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ManageCommentLogic) ManageComment(in *pb.ManageCommentReq) (*pb.ManageCommentRes, error) {
	erx := l.svcCtx.CrudModel.ManageComment(l.ctx, in.UserId, in.VideoId, in.ActionType, in.CommentText, in.CommentId)
	if erx != nil {
		if erx.Code() == crud.ErrIdInvalidActionType {
			return &pb.ManageCommentRes{
				StatusCode: errx.Encode(
					errx.Sys,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicCrud,
					crud.ErrIdInvalidActionType,
					erx.Code(),
				),
				StatusMsg: errx.Internal,
			}, nil
		}
		return &pb.ManageCommentRes{
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

	return &pb.ManageCommentRes{
		StatusCode: 0,
		StatusMsg:  "manage comment successfully",
	}, nil
}
