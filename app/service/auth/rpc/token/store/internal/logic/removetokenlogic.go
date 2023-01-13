package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/auth/internal/sys"
	"douyin/app/service/auth/rpc/token/store/internal/consts"
	"douyin/app/service/auth/rpc/token/store/internal/svc"
	"douyin/app/service/auth/rpc/token/store/pb"
	"go.uber.org/zap"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveTokenLogic {
	return &RemoveTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveTokenLogic) RemoveToken(in *pb.RemoveTokenReq) (*pb.RemoveTokenRes, error) {
	if in.Username == "" {
		return &pb.RemoveTokenRes{
			Code: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprRemoveToken,
				consts.ErrIdEmptyParam,
			),
			Msg: consts.ErrEmptyParam,
		}, nil
	}

	// 删除令牌
	err := l.svcCtx.Rdb.Del(
		l.ctx,
		consts.RdbKeyToken+in.Username,
	).Err()
	if err != nil {
		log.Logger.Error(errx.RedisDel, zap.Error(err))
		return &pb.RemoveTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprRemoveToken,
				consts.ErrIdRedisDelOauth2Token,
			),
			Msg: errx.Internal,
		}, nil
	}

	return &pb.RemoveTokenRes{
		Code: http.StatusOK,
		Msg:  "remove token successfully",
	}, nil
}
