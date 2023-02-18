package profile

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/user/rpc/sys/pb"
	videopb "douyin/app/service/video/rpc/sys/pb"
	videosys "douyin/app/service/video/rpc/sys/sys"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/zeromicro/go-zero/core/mr"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type (
	Model interface {
		GetProfile(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profile, errx.Error)
	}
	DefaultModel struct {
		db                *gorm.DB
		rdb               *redis.ClusterClient
		videoSysRpcClient videosys.Sys
	}
)

func NewModel(db *gorm.DB, rdb *redis.ClusterClient) *DefaultModel {
	return &DefaultModel{
		db:  db,
		rdb: rdb,
	}
}

func (m *DefaultModel) GetProfile(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profile, errx.Error) {
	userSubject := &entity.UserSubject{}
	var followCnt, followerCnt, workCnt, favoriteCnt, totalFavorited int64
	var isFollow bool

	err := mr.Finish(func() error {
		err := m.db.WithContext(ctx).
			Select("`id`, `username`").
			Where("`id` = ?", dstUserId).
			Take(userSubject).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errUserNotFound
			}
			log.Logger.Error(errx.MysqlGet, zap.Error(err))
			return errMysqlGet
		}

		return nil
	}, func() error {
		var err error
		followCnt, err = m.rdb.ZCard(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId)).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return errRedisGet
		}

		return nil
	}, func() error {
		var err error
		followerCnt, err = m.rdb.ZCard(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId)).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return errRedisGet
		}

		return nil
	}, func() error {
		var err error
		_, err = m.rdb.ZRank(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId), strconv.FormatInt(dstUserId, 10)).Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return errRedisGet
			}
		} else {
			isFollow = true
		}

		return nil
	}, func() error {
		rpcRes, _ := m.videoSysRpcClient.GetPublishCount(ctx, &videopb.GetPublishCountReq{
			UserId: dstUserId,
		})
		if rpcRes == nil {
			log.Logger.Error(errx.RequestRpcReceive)
			return errRequestRpcReceive
		}
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return errRequestRpcRes
		}
		workCnt = rpcRes.PublishCount

		return nil
	}, func() error {
		rpcRes, _ := m.videoSysRpcClient.GetFavoriteCount(ctx, &videopb.GetFavoriteCountReq{
			UserId: dstUserId,
		})
		if rpcRes == nil {
			log.Logger.Error(errx.RequestRpcReceive)
			return errRequestRpcReceive
		}
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return errRequestRpcRes
		}
		favoriteCnt = rpcRes.FavoriteCount

		return nil
	}, func() error {
		rpcRes, _ := m.videoSysRpcClient.GetTotalFavorited(ctx, &videopb.GetTotalFavoritedReq{
			UserId: dstUserId,
		})
		if rpcRes == nil {
			log.Logger.Error(errx.RequestRpcReceive)
			return errRequestRpcReceive
		}
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return errRequestRpcRes
		}
		totalFavorited = rpcRes.TotalFavorited

		return nil
	})
	if err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	profile := &pb.Profile{
		Id:             userSubject.ID,
		Name:           userSubject.Username,
		FollowCount:    followCnt,
		FollowerCount:  followerCnt,
		IsFollow:       isFollow,
		TotalFavorited: totalFavorited,
		WorkCount:      workCnt,
		FavoriteCount:  favoriteCnt,
	}

	return profile, nil
}
