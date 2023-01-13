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
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenLogic) GetToken(in *pb.GetTokenReq) (*pb.GetTokenRes, error) {
	if in.Username == "" {
		return &pb.GetTokenRes{
			Code: errx.Encode(
				errx.Logic,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprGetToken,
				consts.ErrIdEmptyParam,
			),
			Msg:  consts.ErrEmptyParam,
			Data: nil,
		}, nil
	}

	// 从redis中获取令牌值
	tokenString, err := l.svcCtx.Rdb.Get(
		l.ctx,
		consts.RdbKeyToken+in.Username,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return &pb.GetTokenRes{
				Code: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Rpc,
					sys.ServiceIdRpcStore,
					consts.ErrIdLogic,
					consts.ErrIdOprGetToken,
					consts.ErrIdTokenNotFound,
				),
				Msg: consts.ErrTokenNotFound,
			}, nil
		}

		return &pb.GetTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprGetToken,
				consts.ErrIdRedisGetOauth2Token,
			),
			Msg:  errx.Internal,
			Data: nil,
		}, nil
	}

	// 解析令牌字符串
	token := &pb.Token{}
	err = json.Unmarshal([]byte(tokenString), token)
	if err != nil {
		log.Logger.Error(errx.JsonUnmarshal, zap.Error(err))
		return &pb.GetTokenRes{
			Code: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcStore,
				consts.ErrIdLogic,
				consts.ErrIdOprGetToken,
				consts.ErrIdJsonUnmarshalOauth2Token,
			),
			Msg:  errx.Internal,
			Data: nil,
		}, nil
	}

	return &pb.GetTokenRes{
		Code: 0,
		Msg:  "get token successfully",
		Data: &pb.GetTokenRes_Data{Token: token},
	}, nil
}
