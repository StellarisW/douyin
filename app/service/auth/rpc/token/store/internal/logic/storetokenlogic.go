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
	"encoding/json"
	"go.uber.org/zap"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type StoreTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStoreTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StoreTokenLogic {
	return &StoreTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StoreTokenLogic) StoreToken(in *pb.StoreTokenReq) (*pb.StoreTokenRes, error) {
	if in.Username == "" || in.Token == nil {
		return &pb.StoreTokenRes{
			Code: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprStoreToken,
				consts.ErrIdEmptyParam,
			),
			Msg: consts.ErrEmptyParam,
		}, nil
	}

	// 序列化令牌值
	tokenBytes, err := json.Marshal(in.Token)
	if err != nil {
		log.Logger.Error(errx.JsonMarshal, zap.Error(err))
		return &pb.StoreTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprStoreToken,
				consts.ErrIdJsonMarshalOauth2Token,
			),
			Msg: errx.Internal,
		}, nil
	}

	// 存储令牌值
	err = l.svcCtx.Rdb.Set(
		l.ctx,
		consts.RdbKeyToken+in.Username,
		string(tokenBytes),
		time.Until(time.Unix(in.Token.ExpiresAt, 0)),
	).Err()
	if err != nil {
		log.Logger.Debug(errx.RedisSet, zap.Error(err))
		return &pb.StoreTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprStoreToken,
				consts.ErrIdRedisSetOauth2Token,
			),
			Msg: errx.Internal,
		}, nil
	}

	return &pb.StoreTokenRes{
		Code: 0,
		Msg:  "store token successfully",
	}, nil
}
