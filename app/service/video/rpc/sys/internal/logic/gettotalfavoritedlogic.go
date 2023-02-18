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

type GetTotalFavoritedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTotalFavoritedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTotalFavoritedLogic {
	return &GetTotalFavoritedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTotalFavoritedLogic) GetTotalFavorited(in *pb.GetTotalFavoritedReq) (*pb.GetTotalFavoritedRes, error) {
	totalFavorited, erx := l.svcCtx.InfoModel.GetTotalFavorited(l.ctx, in.UserId)
	if erx != nil {
		return &pb.GetTotalFavoritedRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetTotalFavorited,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetTotalFavoritedRes{
		StatusCode:     0,
		StatusMsg:      "get total favorited successfully",
		TotalFavorited: totalFavorited,
	}, nil
}
