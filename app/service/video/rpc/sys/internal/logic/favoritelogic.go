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

type FavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteLogic) Favorite(in *pb.FavoriteReq) (*pb.FavoriteRes, error) {
	erx := l.svcCtx.CrudModel.Favorite(l.ctx, in.UserId, in.VideoId, in.ActionType)
	if erx != nil {
		if erx.Code() == crud.ErrIdInvalidActionType || erx.Code() == crud.ErrIdAlreadyLike || erx.Code() == crud.ErrIdAlreadyDisLike {
			return &pb.FavoriteRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcSys,
					consts.ErrIdLogicCrud,
					crud.ErrIdOprFavorite,
					erx.Code(),
				),
				StatusMsg: erx.Error(),
			}, nil
		}
		return &pb.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.FavoriteRes{
		StatusCode: 0,
		StatusMsg:  "favorite successfully",
	}, nil
}
