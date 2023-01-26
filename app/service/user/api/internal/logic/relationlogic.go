package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/common/middleware"
	"douyin/app/common/mq/nsq"
	"douyin/app/service/mq/nsq/producer/user"
	"douyin/app/service/user/api/internal/consts"
	"douyin/app/service/user/api/internal/consts/relation"
	"douyin/app/service/user/internal/sys"
	"go.uber.org/zap"
	"strconv"

	"douyin/app/service/user/api/internal/svc"
	"douyin/app/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationLogic {
	return &RelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationLogic) Relation(req *types.RelationReq) (resp *types.RelationRes, err error) {
	if req.ActionType != "1" && req.ActionType != "2" {
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdInvalidActionType,
			),
			StatusMsg: relation.ErrInvalidActionType,
		}, nil
	}

	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	dstUserId, err := strconv.ParseInt(req.ToUserId, 10, 64)
	if err != nil {
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	actionType, err := strconv.ParseUint(req.ActionType, 10, 32)
	if err != nil {
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdParseInt,
			),
			StatusMsg: relation.ErrParseInt,
		}, nil
	}

	if userId == dstUserId {
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdCannotFollowYourself,
			),
			StatusMsg: relation.ErrCannotFollowYourself,
		}, nil
	}

	producer, err := nsq.GetProducer()
	if err != nil {
		log.Logger.Error(errx.GetNsqProducer)
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdGetNsqProducer,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	err = user.Relation(producer, user.RelationMessage{
		SrcUserId:  userId,
		DstUserId:  dstUserId,
		ActionType: uint32(actionType),
	})
	if err != nil {
		log.Logger.Error(errx.NsqPublish, zap.Error(err))
		return &types.RelationRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicRelation,
				relation.ErrIdOprRelation,
				relation.ErrIdNsqPublish,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &types.RelationRes{
		StatusCode: 0,
		StatusMsg:  "relation successfully",
	}, nil
}
