package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"douyin/app/service/chat/api/internal/consts"
	"douyin/app/service/chat/api/internal/consts/chat"
	"douyin/app/service/chat/internal/sys"
	"douyin/app/service/chat/rpc/sys/pb"
	"strconv"

	"douyin/app/service/chat/api/internal/svc"
	"douyin/app/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.GetMessageListReq) (resp *types.GetMessageListRes, err error) {
	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.GetMessageListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprGetMessageList,
				chat.ErrIdParseInt,
			),
			StatusMsg: chat.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.ToUserId, 10, 64)
	if err != nil {
		return &types.GetMessageListRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprGetMessageList,
				chat.ErrIdParseInt,
			),
			StatusMsg: chat.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.GetMessage(l.ctx, &pb.GetMessageReq{
		SrcUserId: userId,
		DstUserId: dstUserId,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.GetMessageListRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprGetMessageList,
				chat.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg:   errx.Internal,
			MessageList: nil,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.GetMessageListRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.GetMessageListRes{
		StatusCode:  0,
		StatusMsg:   "get message list successfully",
		MessageList: rpcRes.Messages,
	}, nil
}
