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

type StoreMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStoreMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreMessageLogic {
	return &StoreMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StoreMessageLogic) StoreMessage(in *pb.StoreMessageReq) (*pb.StoreMessageRes, error) {
	messageId := l.svcCtx.IdGenerator.NewLong()

	err := l.svcCtx.Db.WithContext(l.ctx).
		Create(&entity.ChatMessage{
			ID:        messageId,
			SrcUserID: in.SrcUserId,
			DstUserID: in.DstUserId,
			Content:   in.Content,
		}).Error
	if err != nil {
		log.Logger.Error(errx.MysqlInsert, zap.Error(err))
		return &pb.StoreMessageRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogic,
				consts.ErrIdOprStoreMessage,
				consts.ErrIdMysqlInsert,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.StoreMessageRes{
		StatusCode: 0,
		StatusMsg:  "store message successfully",
	}, nil
}
