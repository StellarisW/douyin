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
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"sync"
	"time"
)

type (
	Model interface {
		Feed(ctx context.Context, latestTime int64, srcUserId int64) ([]*pb.Video, int64, errx.Error)
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

func (m *DefaultModel) Feed(ctx context.Context, latestTime int64, srcUserId int64) ([]*pb.Video, int64, errx.Error) {
	videoSubjects := make([]*entity.VideoSubject, 0)

	var filterTime time.Time

	if latestTime == 0 {
		filterTime = time.Now()
	} else {
		filterTime = time.UnixMilli(latestTime)
	}

	err := m.db.WithContext(ctx).
		Where("`update_time` < ?", filterTime).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: "update_time"},
			Desc:   true,
		}).Limit(30).
		Find(&videoSubjects).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, 0, errMysqlGet
	}

	videos, erx := m.getVideosInfo(ctx, srcUserId, videoSubjects)
	if erx != nil {
		return nil, 0, erx
	}

	var lateTime int64
	if len(videos) > 0 {
		latestTime = videoSubjects[len(videoSubjects)-1].UpdateTime.Unix()
	} else {
		lateTime = latestTime
	}

	return videos, lateTime, nil
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

	videos, erx := m.getVideosInfo(ctx, srcUserId, videoSubjects)
	if erx != nil {
		return nil, erx
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
		Table(entity.TableNameVideoSubject).
		Where(" `id` IN ?", videoIds).
		Find(&videoSubjects).
		Error
	if err != nil {
		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return nil, errMysqlGet
	}

	videos, erx := m.getVideosInfo(ctx, srcUserId, videoSubjects)
	if erx != nil {
		return nil, erx
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

	comments := make([]*pb.Comment, len(commentSubjects))

	size := len(commentSubjects)

	eg := new(errgroup.Group)

	for i := 0; i < size; i++ {
		i := i
		eg.Go(func() error {
			rpcRes, _ := m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
				SrcUserId: userId,
				DstUserId: commentSubjects[i].UserID,
			})
			if rpcRes == nil {
				log.Logger.Error(errx.RequestRpcReceive)
				return errRequestRpcReceive
			}
			if rpcRes.StatusCode != 0 {
				log.Logger.Error(errx.RequestRpcRes, zap.Uint32("code", rpcRes.StatusCode))
				return errRequestRpcRes
			}

			comments[i] = &pb.Comment{
				Id: commentSubjects[i].ID,
				User: &pb.Profile{
					Id:              commentSubjects[i].UserID,
					Name:            rpcRes.User.Name,
					FollowCount:     rpcRes.User.FollowCount,
					FollowerCount:   rpcRes.User.FollowerCount,
					IsFollow:        rpcRes.User.IsFollow,
					Avatar:          rpcRes.User.Avatar,
					BackgroundImage: rpcRes.User.BackgroundImage,
					Signature:       rpcRes.User.Signature,
					TotalFavorited:  rpcRes.User.TotalFavorited,
					WorkCount:       rpcRes.User.WorkCount,
					FavoriteCount:   rpcRes.User.FavoriteCount,
				},
				Content:    commentSubjects[i].CommentText,
				CreateDate: commentSubjects[i].CreateTime.Format("01-02"),
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	return comments, nil
}

func (m *DefaultModel) getVideosInfo(ctx context.Context, srcUserId int64, videoSubjects []*entity.VideoSubject) ([]*pb.Video, errx.Error) {
	videos := make([]*pb.Video, len(videoSubjects))

	size := len(videoSubjects)

	eg := new(errgroup.Group)

	for i := 0; i < size; i++ {
		i := i
		eg.Go(func() error {
			wg := sync.WaitGroup{}

			rpcRes := &userpb.GetProfileRes{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				rpcRes, _ = m.userSysRpcClient.GetProfile(ctx, &userpb.GetProfileReq{
					SrcUserId: srcUserId,
					DstUserId: videoSubjects[i].UserID,
				})
			}()

			var favoriteCount, commentCount int64
			var isFavorite bool
			var erx errx.Error

			wg.Add(1)
			go func() {
				defer wg.Done()
				cmds, err := m.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
					pipe.Get(ctx,
						fmt.Sprintf("%s%d", video.RdbKeyFavoriteCnt, videoSubjects[i].ID),
					)

					pipe.Get(ctx,
						fmt.Sprintf("%s%d", video.RdbKeyCommentCnt, videoSubjects[i].ID),
					)

					pipe.ZRank(ctx,
						fmt.Sprintf("%s%d", video.RdbKeyFavorite, srcUserId), strconv.FormatInt(videoSubjects[i].ID, 10),
					)

					return nil
				})
				if err != nil {
					if err != redis.Nil {
						log.Logger.Error(errx.RedisPipeExec, zap.Error(err))
						erx = errRedisPipeExec
						return
					}
				}

				favoriteCount, err = cmds[0].(*redis.StringCmd).Int64()
				if err != nil {
					if err != redis.Nil {
						log.Logger.Error(errx.RedisGet, zap.Error(err))
						erx = errRedisGet
						return
					}
				}

				commentCount, err = cmds[1].(*redis.StringCmd).Int64()
				if err != nil {
					if err != redis.Nil {
						log.Logger.Error(errx.RedisGet, zap.Error(err))
						erx = errRedisGet
						return
					}
				}

				_, err = cmds[2].(*redis.IntCmd).Result()
				if err != nil {
					if err != redis.Nil {
						log.Logger.Error(errx.RedisGet, zap.Error(err))
						erx = errRedisGet
						return
					}
				} else {
					isFavorite = true
				}
			}()

			wg.Wait()

			if rpcRes == nil {
				log.Logger.Error(errx.RequestRpcReceive)
				return errRequestRpcReceive
			}
			if rpcRes.StatusCode != 0 {
				log.Logger.Error(errx.RequestRpcRes)
				return errRequestRpcRes
			}

			if erx != nil {
				return erx
			}

			videos[i] = &pb.Video{
				Id: videoSubjects[i].ID,
				Author: &pb.Profile{
					Id:              videoSubjects[i].UserID,
					Name:            rpcRes.User.Name,
					FollowCount:     rpcRes.User.FollowCount,
					FollowerCount:   rpcRes.User.FollowerCount,
					IsFollow:        rpcRes.User.IsFollow,
					Avatar:          rpcRes.User.Avatar,
					BackgroundImage: rpcRes.User.BackgroundImage,
					Signature:       rpcRes.User.Signature,
					TotalFavorited:  rpcRes.User.TotalFavorited,
					WorkCount:       rpcRes.User.WorkCount,
					FavoriteCount:   rpcRes.User.FavoriteCount,
				},
				PlayUrl:       videoSubjects[i].PlayURL,
				CoverUrl:      videoSubjects[i].CoverURL,
				FavoriteCount: favoriteCount,
				CommentCount:  commentCount,
				IsFavorite:    isFavorite,
				Title:         videoSubjects[i].Title,
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errx.New(errx.GetCode(err), err.Error())
	}

	return videos, nil
}
