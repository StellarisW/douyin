package relation

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/user/rpc/sys/pb"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type (
	Model interface {
		Relation(ctx context.Context, srcUserId, dstUserId int64, actionType uint32) errx.Error
		GetFollowList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error)
		GetFollowerList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error)
		GetFriendList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error)
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

func (m *DefaultModel) Relation(ctx context.Context, srcUserId, dstUserId int64, actionType uint32) errx.Error {
	switch actionType {
	case 1:
		// 关注

		now := time.Now()

		err := m.rdb.ZAdd(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			redis.Z{
				Score:  float64(now.Unix()),
				Member: dstUserId,
			}).Err()
		if err != nil {
			log.Logger.Error(errx.RedisAdd, zap.Error(err))
			return errRedisAdd
		}

		err = m.rdb.ZAdd(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
			redis.Z{
				Score:  float64(now.Unix()),
				Member: srcUserId,
			}).Err()
		if err != nil {
			log.Logger.Error(errx.RedisAdd, zap.Error(err))
			return errRedisAdd
		}

		return nil
	case 2:
		// 取消关注

		err := m.rdb.ZRem(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			dstUserId,
		).Err()
		if err != nil {
			log.Logger.Error(errx.RedisRem, zap.Error(err))
			return errRedisRem
		}

		err = m.rdb.ZRem(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
			srcUserId,
		).Err()
		if err != nil {
			log.Logger.Error(errx.RedisRem, zap.Error(err))
			return errRedisRem
		}

		return nil
	default:
		return errInvalidActionType
	}
}

func (m *DefaultModel) GetFollowList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error) {
	ids, err := m.rdb.ZRange(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId), 0, -1).Result()
	if err != nil {
		log.Logger.Error(errx.RedisRange, zap.Error(err))
		return nil, errRedisRange
	}

	interIds, err := m.rdb.ZInter(ctx, &redis.ZStore{
		Keys: []string{
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId),
		},
	}).Result()
	if err != nil {
		log.Logger.Error(errx.RedisInter, zap.Error(err))
		return nil, errRedisInter
	}

	profiles := make([]*pb.Profile, 0, len(ids))

	interIndex := 0
	for _, id := range ids {
		var username string

		err = m.db.WithContext(ctx).
			Table(entity.TableNameUserSubject).
			Select("`username`").
			Where("`id` = ?", id).
			Take(&username).
			Error
		if err != nil {
			log.Logger.Error(errx.MysqlGet, zap.Error(err))
			return nil, errMysqlGet
		}

		followCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		followerCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		if id == interIds[interIndex] {
			profiles = append(profiles, &pb.Profile{
				Id:            cast.ToInt64(id),
				Name:          username,
				FollowCount:   followCnt,
				FollowerCount: followerCnt,
				IsFollow:      true,
			})
			interIndex++
		} else {
			profiles = append(profiles, &pb.Profile{
				Id:            cast.ToInt64(id),
				Name:          username,
				FollowCount:   followCnt,
				FollowerCount: followerCnt,
				IsFollow:      false,
			})
		}
	}

	return &pb.Profiles{
		Profiles: profiles,
	}, nil
}

func (m *DefaultModel) GetFollowerList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error) {
	ids, err := m.rdb.ZRange(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollower, srcUserId), 0, -1).Result()
	if err != nil {
		log.Logger.Error(errx.RedisRange, zap.Error(err))
		return nil, errRedisRange
	}

	interIds, err := m.rdb.ZInter(ctx, &redis.ZStore{
		Keys: []string{
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
		},
	}).Result()
	if err != nil {
		log.Logger.Error(errx.RedisInter, zap.Error(err))
		return nil, errRedisInter
	}

	profiles := make([]*pb.Profile, 0, len(ids))

	interIndex := 0
	for _, id := range ids {
		var username string

		err = m.db.WithContext(ctx).
			Table(entity.TableNameUserSubject).
			Select("`username`").
			Where("`id` = ?", id).
			Take(&username).
			Error
		if err != nil {
			log.Logger.Error(errx.MysqlGet, zap.Error(err))
			return nil, errMysqlGet
		}

		followCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		followerCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		if id == interIds[interIndex] {
			profiles = append(profiles, &pb.Profile{
				Id:            cast.ToInt64(id),
				Name:          username,
				FollowCount:   followCnt,
				FollowerCount: followerCnt,
				IsFollow:      true,
			})
			interIndex++
		} else {
			profiles = append(profiles, &pb.Profile{
				Id:            cast.ToInt64(id),
				Name:          username,
				FollowCount:   followCnt,
				FollowerCount: followerCnt,
				IsFollow:      false,
			})
		}
	}

	return &pb.Profiles{
		Profiles: profiles,
	}, nil
}

func (m *DefaultModel) GetFriendList(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profiles, errx.Error) {
	ids, err := m.rdb.ZRange(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId), 0, -1).Result()
	if err != nil {
		log.Logger.Error(errx.RedisRange, zap.Error(err))
		return nil, errRedisRange
	}

	friendIds, err := m.rdb.ZInter(ctx, &redis.ZStore{
		Keys: []string{
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
		},
	}).Result()
	if err != nil {
		log.Logger.Error(errx.RedisInter, zap.Error(err))
		return nil, errRedisInter
	}

	profiles := make([]*pb.Profile, 0, len(friendIds))

	idsLen := len(ids)
	idsIndex := 0
	isFollow := false
	for _, id := range friendIds {
		for i := idsIndex; i <= idsLen; i++ {
			if id == ids[i] {
				isFollow = true
				idsIndex = i + 1
			}
		}

		var username string

		err = m.db.WithContext(ctx).
			Table(entity.TableNameUserSubject).
			Select("`username`").
			Where("`id` = ?", id).
			Take(&username).
			Error
		if err != nil {
			log.Logger.Error(errx.MysqlGet, zap.Error(err))
			return nil, errMysqlGet
		}

		followCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		followerCnt, err := m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		profiles = append(profiles, &pb.Profile{
			Id:            cast.ToInt64(id),
			Name:          username,
			FollowCount:   followCnt,
			FollowerCount: followerCnt,
			IsFollow:      isFollow,
		})

		isFollow = false
	}

	return &pb.Profiles{
		Profiles: profiles,
	}, nil
}
