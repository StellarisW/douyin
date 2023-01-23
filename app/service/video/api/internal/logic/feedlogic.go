package logic

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/video/api/internal/consts"
	"douyin/app/service/video/api/internal/consts/info"
	"douyin/app/service/video/api/internal/svc"
	"douyin/app/service/video/api/internal/types"
	"douyin/app/service/video/internal/sys"
	"douyin/app/service/video/rpc/sys/pb"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(requ *types.FeedReq) (resp *types.FeedRes, err error) {
	var latestTime int64
	if requ.LastestTime == "" {
		latestTime = 0
	} else {
		latestTime, err = strconv.ParseInt(requ.LastestTime, 10, 64)
		if err != nil {
			return &types.FeedRes{
				StatusCode: errx.Encode(
					errx.Logic,
					sys.SysId,
					douyin.Api,
					sys.ServiceIdApi,
					consts.ErrIdLogicInfo,
					info.ErrIdOprFeed,
					info.ErrIdParseInt,
				),
				StatusMsg: info.ErrParseInt,
			}, nil
		}
	}

	var userId int64

	readTokenRes, err := req.NewRequest().
		SetFormData(map[string]string{
			"token_value": requ.Token,
		}).
		Post("http://douyin-auth-api:11120/douyin/token/read")
	if err == nil {
		readTokenResJson := gjson.Parse(readTokenRes.String())
		if readTokenResJson.Get("code").Uint() == 0 {
			tokenPayload := readTokenResJson.Get("data.payload").String()
			tokenPayloadJson := gjson.Parse(tokenPayload)
			userId, _ = strconv.ParseInt(tokenPayloadJson.Get("sub").String(), 10, 64)
		}
	}

	rpcRes, _ := l.svcCtx.SysRpcClient.Feed(l.ctx, &pb.FeedReq{
		LatestTime: latestTime,
		UserId:     userId,
	})
	if rpcRes == nil {
		log.Logger.Error(errx.RequestRpcReceive)
		return &types.FeedRes{
			StatusCode: errx.Encode(
				errx.Sys,
				sys.SysId,
				douyin.Api,
				sys.ServiceIdApi,
				consts.ErrIdLogicInfo,
				info.ErrIdOprGetCommentList,
				info.ErrIdRequestRpcReceiveSys,
			),
			StatusMsg: errx.Internal,
		}, nil
	} else if rpcRes.StatusCode != 0 {
		return &types.FeedRes{
			StatusCode: rpcRes.StatusCode,
			StatusMsg:  rpcRes.StatusMsg,
		}, nil
	}

	return &types.FeedRes{
		StatusCode: 0,
		StatusMsg:  "get feed successfully",
		NextTime:   rpcRes.NextTime,
		VideoList:  rpcRes.Videos,
	}, nil
}
