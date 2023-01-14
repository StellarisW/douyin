package profile

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/user/rpc/sys/pb"
	"fmt"
	"github.com/go-redis/redis/v9"
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

	err := m.db.WithContext(ctx).
		Select("`id`, `username`").
		Where("`id` = ?").
		Take(userSubject).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errUserNotFound
		}

		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, errMysqlGet
	}

	followCnt, err := m.rdb.ZCard(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId)).Result()
	if err != nil {
		log.Logger.Error(errx.RedisGet, zap.Error(err))
		return nil, errRedisGet
	}

	followerCnt, err := m.rdb.ZCard(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId)).Result()
	if err != nil {
		log.Logger.Error(errx.RedisGet, zap.Error(err))
		return nil, errRedisGet
	}

	profile := &pb.Profile{
		Id:            userSubject.ID,
		Name:          userSubject.Username,
		FollowCount:   followCnt,
		FollowerCount: followerCnt,
		IsFollow:      false,
	}

	_, err = m.rdb.ZRank(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId), strconv.FormatInt(dstUserId, 10)).Result()
	if err != nil {
		if err != redis.Nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}
	} else {
		profile.IsFollow = true
	}

	return profile, nil
}
