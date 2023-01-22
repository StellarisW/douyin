package info

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlGet          = errx.New(ErrIdMysqlGet, errx.MysqlGet)
	errRequestRpcRes     = errx.New(ErrIdRequestRpcRes, errx.RequestRpcRes)
	errRequestRpcReceive = errx.New(ErrIdRequestRpcReceive, errx.RequestRpcReceive)
	errRedisGet          = errx.New(ErrIdRedisGet, errx.RedisGet)
	errRedisPipeExec     = errx.New(ErrIdRedisPipeExec, errx.RedisPipeExec)
	errRedisRange        = errx.New(ErrIdRedisRange, errx.RedisRange)
)
