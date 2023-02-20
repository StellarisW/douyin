package crud

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlInsert       = errx.New(ErrIdMysqlInsert, errx.MysqlInsert)
	errMinioPut          = errx.New(ErrIdMinioPut, errx.MinioPut)
	errRedisAdd          = errx.New(ErrIdRedisAdd, errx.RedisAdd)
	errRedisGet          = errx.New(ErrIdRedisGet, errx.RedisGet)
	errRedisIncr         = errx.New(ErrIdRedisIncr, errx.RedisIncr)
	errRedisRem          = errx.New(ErrIdRedisRem, errx.RedisRem)
	errRedisDecr         = errx.New(ErrIdRedisDecr, errx.RedisDecr)
	errRedisPipeExec     = errx.New(ErrIdRedisPipeExec, errx.RedisPipeExec)
	errInvalidActionType = errx.New(ErrIdInvalidActionType, ErrInvalidActionType)
	errMysqlDelete       = errx.New(ErrIdMysqlDelete, errx.MysqlDelete)
	errRequestRpcRes     = errx.New(ErrIdRequestRpcRes, errx.RequestRpcRes)
	errRequestRpcReceive = errx.New(ErrIdRequestRpcReceive, errx.RequestRpcReceive)
)

var (
	// Publish 操作

	errInvalidContentType = errx.New(ErrIdInvalidContentType, ErrInvalidContentType)
)

var (
	// Favorite 操作
	errAlreadyLike    = errx.New(ErrIdAlreadyLike, ErrAlreadyLike)
	errAlreadyDislike = errx.New(ErrIdAlreadyDisLike, ErrAlreadyDisLike)
)
