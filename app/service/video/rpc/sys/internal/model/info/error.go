package info

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlGet      = errx.New(ErrIdMysqlGet, errx.MysqlGet)
	errRequestRpcRes = errx.New(ErrIdRequestRpcRes, errx.RequestRpcRes)
	errRedisGet      = errx.New(ErrIdRedisGet, errx.RedisGet)
	errRedisRange    = errx.New(ErrIdRedisRange, errx.RedisRange)
)
