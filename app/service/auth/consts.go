package auth

const (
	ErrTokenNotSupportedType = "token type not supported"
	ErrInvalidSignature      = "invalid signature"
	ErrInvalidKey            = "invalid key"

	ErrTokenExpired     = "token is expired"
	ErrTokenNotValidYet = "token not active yet"
	ErrTokenInvalidType = "invalid token type"
)
