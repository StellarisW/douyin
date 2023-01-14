package info

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	userpb "douyin/app/service/user/rpc/sys/pb"
	usersys "douyin/app/service/user/rpc/sys/sys"
	"douyin/app/service/video/internal/video"
	"douyin/app/service/video/rpc/sys/internal/model/dao/entity"
	"douyin/app/service/video/rpc/sys/pb"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type (
	Model interface {
		Feed(ctx context.Context, latestTime int64, userId int64) ([]*pb.Video, int64, errx.Error)
		GetPublishList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Video, errx.Error)
		GetFavoriteList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Video, errx.Error)
		GetCommentList(ctx context.Context, userId, videoId int64) ([]*pb.Comment, errx.Error)
	}
	DefaultModel struct {
		db               *gorm.DB
		rdb              *redis.ClusterClient
		userSysRpcClient usersys.Sys
	}
)

func NewModel(db *gorm.DB, rdb *redis.ClusterClient, userSysRpcClient usersys.Sys) *DefaultModel {
	return &DefaultModel{
		db:               db,
		rdb:              rdb,
		userSysRpcClient: userSysRpcClient,
	}
}

func (m *DefaultModel) Feed(ctx context.Context, latestTime int64, userId int64) ([]*pb.Video, int64, errx.Error) {
	videoSubjects := make([]*entity.VideoSubject, 0)

	var filterTime time.Time

	if latestTime == 0 {
		filterTime = time.Now()
	} else {
		filterTime = time.Unix(latestTime, 0)
	}

	err := m.db.WithContext(ctx).
		Where("`update_time` < ?", filterTime).
		Find(&videoSubjects).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "update_time"},
			Desc:   true,
		}).Limit(30).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, 0, errMysqlGet
	}

	videos := make([]*pb.Video, 0, len(videoSubjects))

	for _, videoSubject := range videoSubjects {
		rpcRes, _ := m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
			SrcUserId: userId,
			DstUserId: videoSubject.UserID,
		})
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return nil, 0, errRequestRpcRes
		}

		var favoriteCount, commentCount int64

		favoriteCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, 0, errRedisGet
		}

		commentCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, 0, errRedisGet
		}

		var isFavorite bool

		_, err = m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, userId), strconv.FormatInt(videoSubject.ID, 10)).
			Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return nil, 0, errRedisGet
			}
		} else {
			isFavorite = true
		}

		videos = append(videos, &pb.Video{
			Id: videoSubject.ID,
			User: &pb.Profile{
				Id:            videoSubject.UserID,
				Name:          rpcRes.User.Name,
				FollowCount:   rpcRes.User.FollowCount,
				FollowerCount: rpcRes.User.FollowerCount,
				IsFollow:      rpcRes.User.IsFollow,
			},
			PlayUrl:       videoSubject.PlayURL,
			CoverUrl:      videoSubject.CoverURL,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoSubject.Title,
		})
	}

	return videos, videoSubjects[len(videoSubjects)-1].UpdateTime.Unix(), nil
}

func (m *DefaultModel) GetPublishList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Video, errx.Error) {
	videoSubjects := make([]*entity.VideoSubject, 0)

	err := m.db.WithContext(ctx).
		Where("`user_id` = ?", dstUserId).
		Find(&videoSubjects).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "update_time"},
			Desc:   true,
		}).Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, errMysqlGet
	}

	videos := make([]*pb.Video, 0, len(videoSubjects))

	for _, videoSubject := range videoSubjects {
		rpcRes, _ := m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
			SrcUserId: srcUserId,
			DstUserId: videoSubject.UserID,
		})
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return nil, errRequestRpcRes
		}

		var favoriteCount, commentCount int64

		favoriteCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		commentCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		var isFavorite bool

		_, err = m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, srcUserId), strconv.FormatInt(videoSubject.ID, 10)).
			Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return nil, errRedisGet
			}
		} else {
			isFavorite = true
		}

		videos = append(videos, &pb.Video{
			Id: videoSubject.ID,
			User: &pb.Profile{
				Id:            videoSubject.UserID,
				Name:          rpcRes.User.Name,
				FollowCount:   rpcRes.User.FollowCount,
				FollowerCount: rpcRes.User.FollowerCount,
				IsFollow:      rpcRes.User.IsFollow,
			},
			PlayUrl:       videoSubject.PlayURL,
			CoverUrl:      videoSubject.CoverURL,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoSubject.Title,
		})
	}

	return videos, nil
}

func (m *DefaultModel) GetFavoriteList(ctx context.Context, srcUserId, dstUserId int64) ([]*pb.Video, errx.Error) {
	videoIds, err := m.rdb.ZRange(ctx, fmt.Sprintf("%s%d", video.RdbKeyFavorite, dstUserId), 0, -1).Result()
	if err != nil {
		log.Logger.Error(errx.RedisRange, zap.Error(err))
		return nil, errRedisRange
	}

	videoSubjects := make([]*entity.VideoSubject, 0)

	err = m.db.WithContext(ctx).
		Where(" `id` IN ?", videoIds).
		Find(&videoIds).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, errMysqlGet
	}

	videos := make([]*pb.Video, 0, len(videoSubjects))

	for _, videoSubject := range videoSubjects {
		rpcRes, _ := m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
			SrcUserId: srcUserId,
			DstUserId: videoSubject.UserID,
		})
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return nil, errRequestRpcRes
		}

		var favoriteCount, commentCount int64

		favoriteCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		commentCount, err = m.rdb.Get(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoSubject.ID),
		).Int64()
		if err != nil {
			log.Logger.Error(errx.RedisGet, zap.Error(err))
			return nil, errRedisGet
		}

		var isFavorite bool

		_, err = m.rdb.ZRank(ctx,
			fmt.Sprintf("%s%d", video.RdbKeyFavorite, srcUserId), strconv.FormatInt(videoSubject.ID, 10)).
			Result()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return nil, errRedisGet
			}
		} else {
			isFavorite = true
		}

		videos = append(videos, &pb.Video{
			Id: videoSubject.ID,
			User: &pb.Profile{
				Id:            videoSubject.UserID,
				Name:          rpcRes.User.Name,
				FollowCount:   rpcRes.User.FollowCount,
				FollowerCount: rpcRes.User.FollowerCount,
				IsFollow:      rpcRes.User.IsFollow,
			},
			PlayUrl:       videoSubject.PlayURL,
			CoverUrl:      videoSubject.CoverURL,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoSubject.Title,
		})
	}

	return videos, nil
}

func (m *DefaultModel) GetCommentList(ctx context.Context, userId, videoId int64) ([]*pb.Comment, errx.Error) {
	commentSubjects := make([]*entity.CommentSubject, 0)

	err := m.db.WithContext(ctx).
		Where(" `video_id` = ?", videoId).
		Find(&commentSubjects).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, errMysqlGet
	}

	comments := make([]*pb.Comment, 0, len(commentSubjects))

	for _, commentSubject := range commentSubjects {
		rpcRes, _ := m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
			SrcUserId: userId,
			DstUserId: commentSubject.UserID,
		})
		if rpcRes.StatusCode != 0 {
			log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
			return nil, errRequestRpcRes
		}

		comments = append(comments, &pb.Comment{
			Id: commentSubject.ID,
			User: &pb.Profile{
				Id:            commentSubject.UserID,
				Name:          rpcRes.User.Name,
				FollowCount:   rpcRes.User.FollowCount,
				FollowerCount: rpcRes.User.FollowerCount,
				IsFollow:      rpcRes.User.IsFollow,
			},
			Content:    commentSubject.CommentText,
			CreateDate: commentSubject.CreateTime.Format("06-01"),
		})
	}

	return comments, nil
}
