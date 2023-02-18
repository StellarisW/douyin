package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/internal/model/consts"
	"douyin/app/service/video/rpc/sys/internal/model/info"

	"douyin/app/service/video/rpc/sys/internal/svc"
	"douyin/app/service/video/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountLogic {
	return &GetFavoriteCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountLogic) GetFavoriteCount(in *pb.GetFavoriteCountReq) (*pb.GetFavoriteCountRes, error) {
	favoriteCnt, erx := l.svcCtx.InfoModel.GetFavoriteCount(l.ctx, in.UserId)
	if erx != nil {
		return &pb.GetFavoriteCountRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetFavoriteCount,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetFavoriteCountRes{
		StatusCode:    0,
		StatusMsg:     "get favorite count successfully",
		FavoriteCount: favoriteCnt,
	}, nil
}
