package crud

import (
	"bytes"
	"context"
	"douyin/app/common/douyin"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/video/internal/video"
	"douyin/app/service/video/rpc/sys/internal/model/dao/entity"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type (
	Model interface {
		Publish(ctx context.Context, userId int64, title string, data []byte) errx.Error
		Favorite(ctx context.Context, userId, videoId int64, actionType uint32) errx.Error
		Comment(ctx context.Context, userId, videoId int64, actionType uint32, commentText string, commentId int64) errx.Error
	}
	DefaultModel struct {
		idGenerator *idgen.DefaultIdGenerator
		db          *gorm.DB
		rdb         *redis.ClusterClient
		minioClient *minio.Client
	}
)

func NewModel(idGenerator *idgen.DefaultIdGenerator, db *gorm.DB, rdb *redis.ClusterClient, minioClient *minio.Client) *DefaultModel {
	return &DefaultModel{
		idGenerator: idGenerator,
		db:          db,
		rdb:         rdb,
		minioClient: minioClient,
	}
}

func (m *DefaultModel) Publish(ctx context.Context, userId int64, title string, data []byte) errx.Error {
	videoId := m.idGenerator.NewLong()
	if title == "" {
		title = defaultVideoTitle
	}
	playUrl := fmt.Sprintf("%s/%s/video/%d/video", douyin.MinioDomain, douyin.MinioBucket, videoId)
	coverUrl := defaultVideoCoverUrl

	err := m.db.WithContext(ctx).
		Create(&entity.VideoSubject{
			ID:       videoId,
			UserID:   userId,
			PlayURL:  playUrl,
			CoverURL: coverUrl,
			Title:    title,
		}).Error
	if err != nil {
		log.Logger.Error(errx.MysqlInsert, zap.Error(err))
		return errMysqlInsert
	}

	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "video") {
		return errInvalidContentType
	}

	buffer := bytes.NewBuffer(data)

	_, err = m.minioClient.PutObject(ctx,
		douyin.MinioBucket,
		fmt.Sprintf("video/%d/video", videoId),
		buffer,
		int64(buffer.Len()),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Logger.Error(errx.MinioPut, zap.Error(err))
		return errMinioPut
	}

	return nil
}

func (m *DefaultModel) Favorite(ctx context.Context, userId, videoId int64, actionType uint32) errx.Error {
	switch actionType {
	case 1:
		// 点赞

		now := time.Now()

		err := m.rdb.ZAdd(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
			redis.Z{
				Score:  float64(now.Unix()),
				Member: videoId,
			}).Err()
		if err != nil {
			log.Logger.Error(errx.RedisAdd, zap.Error(err))
			return errRedisAdd
		}

		err = m.rdb.Incr(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoId)).
			Err()
		if err != nil {
			log.Logger.Error(errx.RedisIncr, zap.Error(err))
			return errRedisIncr
		}

		return nil
	case 2:
		// 取消点赞

		err := m.rdb.ZRem(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId),
			videoId,
		).Err()
		if err != nil {
			log.Logger.Error(errx.RedisRem, zap.Error(err))
			return errRedisRem
		}

		err = m.rdb.Decr(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoId)).
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

func (m *DefaultModel) Comment(ctx context.Context, userId, videoId int64, actionType uint32, commentText string, commentId int64) errx.Error {
	switch actionType {
	case 1:
		// 发布评论

		commentId = m.idGenerator.NewLong()

		err := m.db.WithContext(ctx).
			Create(&entity.CommentSubject{
				ID:          commentId,
				UserID:      userId,
				VideoID:     videoId,
				CommentText: commentText,
			}).Error
		if err != nil {
			log.Logger.Error(errx.MysqlInsert, zap.Error(err))
			return errMysqlInsert
		}

		err = m.rdb.Incr(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoId)).
			Err()
		if err != nil {
			log.Logger.Error(errx.RedisIncr, zap.Error(err))
			return errRedisIncr
		}

		return nil
	case 2:
		// 删除评论

		err := m.db.WithContext(ctx).
			Where("`id` = ? `user_id` = ?", commentId, userId).
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