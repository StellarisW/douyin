package profile

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlGet          = errx.New(ErrIdMysqlGet, errx.MysqlGet)
	errRequestRpcRes     = errx.New(ErrIdRequestRpcRes, errx.RequestRpcRes)
	errRequestRpcReceive = errx.New(ErrIdRequestRpcReceive, errx.RequestRpcReceive)
	errRedisGet          = errx.New(ErrIdRedisGet, errx.RedisGet)
)

var (
	// GetProfile 操作

	errUserNotFound = errx.New(ErrIdUserNotFound, ErrUserNotFound)
)
