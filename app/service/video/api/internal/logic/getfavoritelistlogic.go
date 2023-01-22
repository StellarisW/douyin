package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/info"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/pb"
	"strconv"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(req *types.GetFavoriteListReq) (resp *types.GetFavoriteListRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetFavoriteListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetFavoriteList,
				info.ErrIdParseInt,
			),
			StatusMsg: info.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return &types.GetFavoriteListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetFavoriteList,
				info.ErrIdParseInt,
			),
			StatusMsg: info.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.GetFavoriteList(l.ctx, &pb.GetFavoriteListReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.GetFavoriteListRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetFavoriteList,
				info.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.GetFavoriteListRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetFavoriteListRes{
		StatusCode: 0,
		StatusMsg:  "get favorite list successfully",
		VideoList:  rpcRes.Videos,
	}, nil
}
