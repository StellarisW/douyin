package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"douyin/app/common/mq/nsq"
	"douyin/app/service/mq/nsq/producer/video"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/crud"
	"douyin/app/service/video/internal/sys"
	"go.uber.org/zap"
	"strconv"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteLogic) Favorite(req *types.FavoriteReq) (resp *types.FavoriteRes, err error) {
	if req.ActionType != "1" && req.ActionType != "2" {
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdInvalidActionType,
			),
			StatusMsg: crud.ErrInvalidActionType,
		}, nil
	}

	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	videoId, err := strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	actionType, err := strconv.ParseUint(req.ActionType, 10, 32)
	if err != nil {
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	producer, err := nsq.GetProducer()
	if err != nil {
		log.Logger.Error(errx.GetNsqProducer)
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdGetNsqProducer,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	err = video.Favorite(producer, video.FavoriteMessage{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: uint32(actionType),
	})
	if err != nil {
		log.Logger.Error(errx.NsqPublish, zap.Error(err))
		return &types.FavoriteRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprFavorite,
				crud.ErrIdNsqPublish,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &types.FavoriteRes{
		StatusCode: 0,
		StatusMsg:  "favorite successfully",
	}, nil
}
