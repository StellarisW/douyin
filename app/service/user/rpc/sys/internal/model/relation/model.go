package relation

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
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
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

		_, err := m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			strconv.FormatInt(dstUserId, 10)).Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return errRedisGet
			}
		} else {
			return errAlreadyFollow
		}

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

		_, err := m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
			strconv.FormatInt(dstUserId, 10)).Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return errRedisGet
			}
			return errAlreadyUnfollow
		}

		_, err = m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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

	if len(interIds) == 0 {
		interIds = append(interIds, "")
	}

	profiles := make([]*pb.Profile, len(ids))

	size := len(ids)

	eg := new(errgroup.Group)

	for i := 0; i < size; i++ {
		i := i

		eg.Go(func() error {
			profiles[i], err = m.getProfile(ctx, srcUserId, dstUserId)
			if err != nil {
				return errRedisGet
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	return profiles, nil
}

func (m *DefaultModel) GetFollowerList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error) {
	ids, err := m.rdb.ZRange(ctx, fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId), 0, -1).Result()
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

	if len(interIds) == 0 {
		interIds = append(interIds, "")
	}

	profiles := make([]*pb.Profile, len(ids))

	size := len(ids)

	eg := new(errgroup.Group)

	for i := 0; i < size; i++ {
		i := i

		eg.Go(func() error {
			profiles[i], err = m.getProfile(ctx, srcUserId, dstUserId)
			if err != nil {
				return errRedisGet
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	return profiles, nil
}

func (m *DefaultModel) GetFriendList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Profile, errx.Error) {
	ids, err := m.rdb.ZInter(ctx, &redis.ZStore{
		Keys: []string{
			fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
		},
	}).Result()
	if err != nil {
		log.Logger.Error(errx.RedisInter, zap.Error(err))
		return nil, errRedisInter
	}

	interIds, err := m.rdb.ZInter(ctx, &redis.ZStore{
		Keys: []string{
			fmt.Sprintf("%s%d", user.RdbKeyFollow, dstUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollower, dstUserId),
			fmt.Sprintf("%s%d", user.RdbKeyFollow, srcUserId),
		},
	}).Result()
	if err != nil {
		log.Logger.Error(errx.RedisInter, zap.Error(err))
		return nil, errRedisInter
	}

	if len(interIds) == 0 {
		interIds = append(interIds, "")
	}

	profiles := make([]*pb.Profile, len(ids))

	size := len(ids)

	eg := new(errgroup.Group)

	for i := 0; i < size; i++ {
		i := i

		eg.Go(func() error {
			profiles[i], err = m.getProfile(ctx, srcUserId, dstUserId)
			if err != nil {
				return errRedisGet
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	return profiles, nil
}

func (m *DefaultModel) getProfile(ctx context.Context, srcUserId, dstUserId int64) (*pb.Profile, errx.Error) {
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
		Id:              userSubject.ID,
		Name:            userSubject.Username,
		FollowCount:     followCnt,
		FollowerCount:   followerCnt,
		IsFollow:        isFollow,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
		TotalFavorited:  totalFavorited,
		WorkCount:       workCnt,
		FavoriteCount:   favoriteCnt,
	}

	return profile, nil
}
