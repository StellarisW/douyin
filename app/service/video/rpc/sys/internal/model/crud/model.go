package crud

import (
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/video/internal/video"
	"douyin/app/service/video/rpc/sys/internal/model/dao/entity"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type (
	Model interface {
		Publish(ctx context.Context, userId int64, title string) (int64, errx.Error)
		Favorite(ctx context.Context, userId, videoId int64, actionType uint32) errx.Error
		Comment(ctx context.Context, userId, videoId int64, actionType uint32, commentText string, commentId int64) errx.Error
	}
	DefaultModel struct {
		idGenerator *idgen.DefaultIdGenerator
		db          *gorm.DB
		rdb         *redis.ClusterClient
	}
)

func NewModel(idGenerator *idgen.DefaultIdGenerator, db *gorm.DB, rdb *redis.ClusterClient) *DefaultModel {
	return &DefaultModel{
		idGenerator: idGenerator,
		db:          db,
		rdb:         rdb,
	}
}

func (m *DefaultModel) Publish(ctx context.Context, userId int64, title string) (int64, errx.Error) {
	videoId := m.idGenerator.NewLong()
	if title == "" {
		title = defaultVideoTitle
	}
	playUrl := fmt.Sprintf("%s/%s/video/%d/video", douyin.MinioDomain, douyin.MinioBucket, videoId)
	coverUrl := fmt.Sprintf("%s/%s/video/%d/image", douyin.MinioDomain, douyin.MinioBucket, videoId)

	videoSubject := &entity.VideoSubject{
		ID:       videoId,
		UserID:   userId,
		PlayURL:  playUrl,
		CoverURL: coverUrl,
		Title:    title,
	}

	err := m.db.WithContext(ctx).
		Create(videoSubject).Error
	if err != nil {
		log.Logger.Error(errx.MysqlInsert, zap.Error(err))
		return 0, errMysqlInsert
	}

	return videoId, nil
}

func (m *DefaultModel) Favorite(ctx context.Context, userId, videoId int64, actionType uint32) errx.Error {
	switch actionType {
	case 1:
		// 点赞

		_, err := m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
			strconv.FormatInt(videoId, 10)).Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return errRedisGet
			}
		} else {
			return errAlreadyLike
		}

		now := time.Now()

		cmds, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZAdd(ctx,
				fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
				redis.Z{
					Score:  float64(now.Unix()),
					Member: videoId,
				})

			pipe.Incr(ctx,
				fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoId),
			)

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
			return errRedisIncr
		}

		return nil
	case 2:
		// 取消点赞

		_, err := m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
			strconv.FormatInt(videoId, 10)).Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return errRedisGet
			}
			return errAlreadyDislike
		}

		cmds, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.ZRem(ctx,
				fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
				videoId,
			)

			pipe.Decr(ctx,
				fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoId),
			)

			return nil
		})
		if err != nil {
			log.Logger.Error(errx.RedisPipeExec, zap.Error(err))
			return errRedisPipeExec
		}
		if cmds[0].(*redis.IntCmd).Err() != nil {
			log.Logger.Error(errx.RedisRem, zap.Error(err))
			return errRedisRem
		}
		if cmds[1].(*redis.IntCmd).Err() != nil {
			log.Logger.Error(errx.RedisDecr, zap.Error(err))
			return errRedisDecr
		}

		return nil

	default:

		return errInvalidActionType
	}
}

func (m *DefaultModel) Comment(ctx context.Context, userId, videoId int64, actionType uint32, commentText string, commentId int64) errx.Error {
	switch actionType {
	case 1:
		// 发布评论

		commentSubject := &entity.CommentSubject{
			ID:          m.idGenerator.NewLong(),
			UserID:      userId,
			VideoID:     videoId,
			CommentText: commentText,
		}

		err := m.db.WithContext(ctx).
			Create(&commentSubject).Error
		if err != nil {
			log.Logger.Error(errx.MysqlInsert, zap.Error(err))
			return errMysqlInsert
		}

		err = m.rdb.Incr(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoId)).
			Err()
		if err != nil {
			m.db.WithContext(ctx).Delete(commentSubject)
			log.Logger.Error(errx.RedisIncr, zap.Error(err))
			return errRedisIncr
		}

		return nil
	case 2:
		// 删除评论

		err := m.db.WithContext(ctx).
			Where("`id` = ? AND `user_id` = ?", commentId, userId).
			Delete(&entity.CommentSubject{}).
			Error
		if err != nil {
			log.Logger.Error(errx.MysqlDelete, zap.Error(err))
			return errMysqlDelete
		}

		err = m.rdb.Decr(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoId)).
			Err()
		if err != nil {
			log.Logger.Error(errx.RedisDecr, zap.Error(err))
			return errRedisDecr
		}

		return nil
	default:
		return errInvalidActionType
	}
}
