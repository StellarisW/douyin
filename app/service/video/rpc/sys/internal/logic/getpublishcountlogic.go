package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/internal/model/consts"
	"douyin/app/service/video/rpc/sys/internal/model/info"

	"douyin/app/service/video/rpc/sys/internal/svc"
	"douyin/app/service/video/rpc/sys/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublishCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishCountLogic {
	return &GetPublishCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPublishCountLogic) GetPublishCount(in *pb.GetPublishCountReq) (*pb.GetPublishCountRes, error) {
	publishCnt, erx := l.svcCtx.InfoModel.GetPublishCount(l.ctx, in.UserId)
	if erx != nil {
		return &pb.GetPublishCountRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Rpc,
				sys.ServiceIdRpcSys,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetPublishCount,
				erx.Code(),
			),
			StatusMsg: errx.Internal,
		}, nil
	}

	return &pb.GetPublishCountRes{
		StatusCode:   0,
		StatusMsg:    "get publish count successfully",
		PublishCount: publishCnt,
	}, nil
}
