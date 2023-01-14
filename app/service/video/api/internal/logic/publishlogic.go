package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/middleware"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/crud"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/pb"
	"strconv"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishReq, data []byte) (resp *types.PublishRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.PublishRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprPublish,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.Publish(l.ctx, &pb.PublishReq{
		UserId: userId,
		Title:  req.Title,
		Data:   data,
	})
	if rpcRes.StatusCode != 0 {
		return &types.PublishRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.PublishRes{
		StatusCode: 0,
		StatusMsg:  "publish successfully",
	}, nil
}
