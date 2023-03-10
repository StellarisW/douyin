package profile

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlGet      = errx.New(ErrIdMysqlGet, errx.MysqlGet)
	errRedisGet      = errx.New(ErrIdRedisGet, errx.RedisGet)
	errRedisPipeExec = errx.New(ErrIdRedisPipeExec, errx.RedisPipeExec)
)

var (
	// GetProfile 操作

	errUserNotFound = errx.New(ErrIdUserNotFound, ErrUserNotFound)
)
