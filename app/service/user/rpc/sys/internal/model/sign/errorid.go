package sign

const (
	// OprId

	ErrIdOprRegister = iota
	ErrIdOprLogin
)

const (
	// 通用错误

	ErrIdMysqlGet = iota
	ErrIdMysqlInsert
	ErrIdRedisSet
	ErrIdRedisGet
	ErrIdRedisDel
	ErrIdRedisIncr
	ErrIdRequestHttpSend
	ErrIdRequestHttpStatusCode
	ErrIdCommon
)

// Register 操作

const (
	ErrIdUsernameExists = iota + ErrIdCommon
)

const (
	ErrUsernameExists = "username already exists"
)

// Login 操作

const (
	ErrIdWrongUsernameOrPassword = iota + ErrIdCommon
)

const (
	ErrWrongUsernameOrPassword = "wrong username or password"
)
