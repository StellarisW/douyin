package jwt

import (
	"douyin/app/common/errx"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var (
	domain string

	issuer     string
	signingKey jwk.Key
	encryptKey jwk.Key

	errTokenBuild     = errx.New(ErrIdTokenBuild, ErrTokenBuild)
	errTokenSerialize = errx.New(ErrIdTokenSerialize, ErrTokenSerialize)

	errInvalidSignature = errx.New(ErrIdInvalidSignature, ErrInvalidSignature)
	errInvalidKey       = errx.New(ErrIdInvalidKey, ErrInvalidKey)

	errTokenNotValidYet = errx.New(ErrIdTokenNotValidYet, ErrTokenNotValidYet)
	errTokenExpired     = errx.New(ErrIdTokenExpired, ErrTokenExpired)
)
