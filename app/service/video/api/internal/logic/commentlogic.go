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
	"douyin/app/service/video/rpc/sys/pb"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"

	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentLogic) Comment(req *types.CommentReq) (resp *types.CommentRes, err error) {
	if req.ActionType != "1" && req.ActionType != "2" {
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdInvalidActionType,
			),
			StatusMsg: crud.ErrInvalidActionType,
		}, nil
	}

	userId, err := strconv.ParseInt(l.ctx.Value(middleware.KeyUserId).(string), 10, 64)
	if err != nil {
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	actionType, err := strconv.ParseUint(req.ActionType, 10, 32)
	if err != nil {
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	var videoId, commentId int64
	var commentContent string

	videoId, err = strconv.ParseInt(req.VideoId, 10, 64)
	if err != nil {
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdParseInt,
			),
			StatusMsg: crud.ErrParseInt,
		}, nil
	}

	_, m, d := time.Now().Date()
	createDate := fmt.Sprintf("%02d-%02d", m, d)

	var comment *pb.Comment

	if actionType == 1 {
		commentContent = l.svcCtx.Filter.GetFilter().Replace(req.CommentText, '*')

		rpcRes, _ := l.svcCtx.SysRpcClient.Comment(l.ctx, &pb.CommentReq{
			UserId:      userId,
			VideoId:     videoId,
			CommentText: commentContent,
			CreateDate:  createDate,
		})
		if rpcRes == nil {
			log.Logger.Error(errx.RequestRpcReceive)
			return &types.CommentRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Api,
					sys.ServiceIdApi,
					consts.ErrIdLogicCrud,
					crud.ErrIdOprComment,
					crud.ErrIdRequestRpcReceiveSys,
				),
				StatusMsg: errx.Internal,
			}, nil
		} else if rpcRes.StatusCode != 0 {
			return &types.CommentRes{
				StatusCode: rpcRes.StatusCode,
				StatusMsg:  rpcRes.StatusMsg,
			}, nil
		}

		comment = rpcRes.Comment
		commentId = comment.Id
	}
	if actionType == 2 {
		commentId, err = strconv.ParseInt(req.CommentId, 10, 64)
		if err != nil {
			return &types.CommentRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Api,
					sys.ServiceIdApi,
					consts.ErrIdLogicCrud,
					crud.ErrIdOprComment,
					crud.ErrIdParseInt,
				),
				StatusMsg: crud.ErrParseInt,
			}, nil
		}
	}

	producer, err := nsq.GetProducer()
	if err != nil {
		log.Logger.Error(errx.GetNsqProducer)
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdGetNsqProducer,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	err = video.Comment(producer, video.CommentMessage{
		UserId:      userId,
		VideoId:     videoId,
		ActionType:  uint32(actionType),
		CommentText: commentContent,
		CommentId:   commentId,
	})
	if err != nil {
		log.Logger.Error(errx.NsqPublish, zap.Error(err))
		return &types.CommentRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicCrud,
				crud.ErrIdOprComment,
				crud.ErrIdNsqPublish,
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	if actionType == 1 {
		return &types.CommentRes{
			StatusCode: 0,
			StatusMsg:  "comment successfully",
			Comment:    comment,
		}, nil
	}
	if actionType == 2 {
		return &types.CommentRes{
			StatusCode: 0,
			StatusMsg:  "delete comment successfully",
		}, nil
	}
	return &types.CommentRes{
		StatusCode: errx.Encode(
			errx.Logic,
			sys.SysId,
			douyin.Api,
			sys.ServiceIdApi,
			consts.ErrIdLogicCrud,
			crud.ErrIdOprComment,
			crud.ErrIdInvalidActionType,
		),
		StatusMsg: crud.ErrInvalidActionType,
	}, nil
}
