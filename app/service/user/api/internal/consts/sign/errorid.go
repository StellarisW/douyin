package sign

const (
	ErrIdOprRegister = iota
	ErrIdOprLogin
)

// 通用错误

const (
	ErrIdRequestRpcReceiveSys = iota + 1
	ErrIdRedisGet
	ErrIdInvalidUsername
	ErrIdInvalidPassword
	ErrIdCommon
)

const (
	ErrInvalidUsername = "invalid username"
	ErrInvalidPassword = "invalid password"
)

// Register

const (
	ErrIdUsernameExist = iota + ErrIdCommon
)

const (
	ErrUsernameExist = "username already exist"
)

// Login

const (
	ErrIdUsernameNotExist = iota + ErrIdCommon
	ErrIdLoginFrozen
)

const (
	ErrUsernameNotExist = "username not exist"
	ErrLoginFrozen      = "user is frozen"
)
