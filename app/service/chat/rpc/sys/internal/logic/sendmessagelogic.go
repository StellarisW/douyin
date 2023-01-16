package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/mq/nsq"
	"douyin/app/service/chat/internal/sys"
	"douyin/app/service/chat/rpc/sys/internal/consts"
	"douyin/app/service/mq/nsq/producer/chat"
	"go.uber.org/zap"

	"douyin/app/service/chat/rpc/sys/internal/svc"
	"douyin/app/service/chat/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageReq) (*pb.SendMessageRes, error) {
	produer, err := nsq.GetProducer()
	if err != nil {
		log.Logger.Error(errx.GetNsqProducer)
		return &pb.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogic,
				consts.ErrIdOprSendMessage,
				consts.ErrIdGetNsqProducer,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	err = chat.Chat(produer, chat.Message{
		SrcUserId:  in.SrcUserId,
		DstUserId:  in.DstUserId,
		ActionType: in.ActionType,
		Content:    in.Content,
	})
	if err != nil {
		log.Logger.Error(errx.NsqPublish, zap.Error(err))
		return &pb.SendMessageRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogic,
				consts.ErrIdOprSendMessage,
				consts.ErrIdNsqPublish,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.SendMessageRes{
		StatusCode: 0,
		StatusMsg:  "send message successfully",
	}, nil
}
