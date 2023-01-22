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
	"github.com/zeromicro/go-zero/core/mr"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type (
	Model interface {
		Relation(ctx context.Context, srcUserId, dstUserId int64, actionType uint32) errx.Error
		GetFollowList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error)
		GetFollowerList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error)
		GetFriendList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error)
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

		cmds, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZAdd(ctx,
				fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
				redis.Z{
					Score:  float64(now.Unix()),
					Member: dstUserId,
				})

			pipe.ZAdd(ctx,
				fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
				redis.Z{
					Score:  float64(now.Unix()),
					Member: srcUserId,
				})

			return nil
		})
		if err != nil {
			log.Logger.Error(errx.RedisPipeExec, zap.Error(err))
			return errRedisPipeExec
		}
		if cmds[0].(*redis.IntCmd).Err() != nil {
			log.Logger.Error(errx.RedisAdd, zap.Error(cmds[0].(*redis.IntCmd).Err()))
			return errRedisAdd
		}
		if cmds[1].(*redis.IntCmd).Err() != nil {
			log.Logger.Error(errx.RedisAdd, zap.Error(cmds[0].(*redis.IntCmd).Err()))
			return errRedisAdd
		}

		return nil
	case 2:
		// 取消关注

		_, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZRem(ctx,
				fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
				dstUserId,
			)

			pipe.ZRem(ctx,
				fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
				srcUserId,
			)

			return nil
		})
		if err != nil {
			log.Logger.Error(errx.RedisRem, zap.Error(err))
			return errRedisRem
		}

		return nil
	default:
		return errInvalidActionType
	}
}

func (m *DefaultModel) GetFollowList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error) {
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
		var followCnt, followerCnt int64

		err = mr.Finish(
			func() error {
				err = m.db.WithContext(ctx).
					Table(entity.TableNameUserSubject).
					Select("`username`").
					Where("`id` = ?", id).
					Take(&username).
					Error
				if err != nil {
					log.Logger.Error(errx.MysqlGet, zap.Error(err))
					return errMysqlGet
				}

				return nil
			},
			func() error {
				followCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			},
			func() error {
				followerCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			})
		if err != nil {
			return nil, errx.New(errx.GetCode(err), err.Error())
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

	return profiles, nil
}

func (m *DefaultModel) GetFollowerList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error) {
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
		var followCnt, followerCnt int64

		err = mr.Finish(
			func() error {
				err = m.db.WithContext(ctx).
					Table(entity.TableNameUserSubject).
					Select("`username`").
					Where("`id` = ?", id).
					Take(&username).
					Error
				if err != nil {
					log.Logger.Error(errx.MysqlGet, zap.Error(err))
					return errMysqlGet
				}

				return nil
			},
			func() error {
				followCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			},
			func() error {
				followerCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			})
		if err != nil {
			return nil, errx.New(errx.GetCode(err), err.Error())
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

	return profiles, nil
}

func (m *DefaultModel) GetFriendList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error) {
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
		var followCnt, followerCnt int64

		err = mr.Finish(
			func() error {
				err = m.db.WithContext(ctx).
					Table(entity.TableNameUserSubject).
					Select("`username`").
					Where("`id` = ?", id).
					Take(&username).
					Error
				if err != nil {
					log.Logger.Error(errx.MysqlGet, zap.Error(err))
					return errMysqlGet
				}

				return nil
			},
			func() error {
				followCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollow+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			},
			func() error {
				followerCnt, err = m.rdb.ZCard(ctx, user.RdbKeyFollower+id).Result()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return errRedisGet
				}

				return nil
			})
		if err != nil {
			return nil, errx.New(errx.GetCode(err), err.Error())
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

	return profiles, nil
}
