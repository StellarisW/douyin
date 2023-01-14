package sign

import "douyin/app/common/errx"

var (
	// 通用错误

	errMysqlGet              = errx.New(ErrIdMysqlGet, errx.MysqlGet)
	errMysqlInsert           = errx.New(ErrIdMysqlInsert, errx.MysqlInsert)
	errRequestHttpSend       = errx.New(ErrIdRequestHttpSend, errx.RequestHttpSend)
	errRequestHttpStatusCode = errx.New(ErrIdRequestHttpStatusCode, errx.RequestHttpStatusCode)
)

var (
	// Register 操作

	errUsernameExists = errx.New(ErrIdUsernameExists, ErrUsernameExists)
)

var (
	// Login 操作

	errWrongUsernameOrPassword = errx.New(ErrIdWrongUsernameOrPassword, ErrWrongUsernameOrPassword)
)
