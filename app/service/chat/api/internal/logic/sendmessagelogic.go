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

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageReq) (resp *types.SendMessageRes, err error) {
	if req.ActionType != "1" {
		return &types.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprSendMessage,
				chat.ErrIdInvalidActionType,
			),
			StatusMsg: chat.ErrInvalidActionType,
		}, nil
	}

	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprSendMessage,
				chat.ErrIdParseInt,
			),
			StatusMsg: chat.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.ToUserId, 10, 64)
	if err != nil {
		return &types.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprSendMessage,
				chat.ErrIdParseInt,
			),
			StatusMsg: chat.ErrParseInt,
		}, nil
	}

	actionType, err := strconv.ParseInt(req.ActionType, 10, 64)
	if err != nil {
		return &types.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprSendMessage,
				chat.ErrIdParseInt,
			),
			StatusMsg: chat.ErrParseInt,
		}, nil
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.SendMessage(l.ctx, &pb.SendMessageReq{
		SrcUserId:  userId,
		DstUserId:  dstUserId,
		ActionType: uint32(actionType),
		Content:    l.svcCtx.Filter.GetFilter().Replace(req.Content, '*'),
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogic,
				chat.ErrIdOprSendMessage,
				chat.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.SendMessageRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.SendMessageRes{
		StatusCode: 0,
		StatusMsg:  "send message successfully",
	}, nil
}
