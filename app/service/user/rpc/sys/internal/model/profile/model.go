package profile

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/internal/video"
	"douyin/app/service/user/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/user/rpc/sys/pb"
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
		db  *gorm.DB
		rdb *redis.ClusterClient
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
	var followCnt, followerCnt, totalFavorited, workCnt, favoriteCnt int64
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
		var err error
		videoSubjects := make([]*entity.VideoSubject, 0)

		err = m.db.WithContext(ctx).
			Table(entity.TableNameVideoSubject).
			Select("`id`").
			Where("`user_id` = ?", dstUserId).
			Find(&videoSubjects).Error
		if err != nil {
			log.Logger.Error(errx.MysqlGet, zap.Error(err))
			return errMysqlGet
		}

		size := len(videoSubjects)

		cmds, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			for i := 0; i < size; i++ {
				pipe.Get(ctx,
					fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoSubjects[i].ID),
				)
			}
			return nil
		})
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisPipeExec, zap.Error(err))
				return errRedisPipeExec
			}
		}

		for i := 0; i < size; i++ {
			value, err := cmds[i].(*redis.StringCmd).Int64()
			if err != nil {
				if err != redis.Nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}
			}
			totalFavorited += value
		}

		return nil
	}, func() error {
		var err error
		err = m.db.WithContext(ctx).
			Model(&entity.VideoSubject{}).
			Where("`user_id` = ?", dstUserId).
			Count(&workCnt).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Logger.Error(errx.MysqlGet, zap.Error(err))
				return errMysqlGet
			}
			workCnt = 0
		}

		return nil
	}, func() error {
		var err error
		favoriteCnt, err = m.rdb.ZCard(ctx, fmt.Sprintf("%s%d", video.RdbKeyFavorite, dstUserId)).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return errRedisGet
		}

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
