package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/chat/internal/sys"
	"douyin/app/service/chat/rpc/sys/internal/consts"
	"douyin/app/service/chat/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/chat/rpc/sys/internal/svc"
	"douyin/app/service/chat/rpc/sys/pb"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
	return &GetMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageLogic) GetMessage(in *pb.GetMessageReq) (*pb.GetMessageRes, error) {
	chatMessages := make([]*entity.ChatMessage, 0)

	err := l.svcCtx.Db.WithContext(l.ctx).
		Select("`id`, `src_user_id`, `dst_user_id`, `content`, `update_time`").
		Where("`src_user_id` = ? AND `dst_user_id` = ?", in.DstUserId, in.SrcUserId).
		Find(&chatMessages).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return &pb.GetMessageRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogic,
				consts.ErrIdOprGetMessage,
				consts.ErrIdMysqlGet,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	messages := make([]*pb.Message, 0, len(chatMessages))

	for _, message := range chatMessages {
		messages = append(messages, &pb.Message{
			Id:         message.ID,
			FromUserId: message.SrcUserID,
			ToUserId:   message.DstUserID,
			Content:    message.Content,
			CreateTime: message.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.GetMessageRes{
		StatusCode: 0,
		StatusMsg:  "get message successfully",
		Messages:   messages,
	}, nil
}
