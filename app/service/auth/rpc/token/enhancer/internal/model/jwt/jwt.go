package jwt

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	auth2 "douyin/app/service/auth/internal/auth"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"time"
)

// InitJWT 初始化JWT令牌参数
func InitJWT() error {
	var err error

	v, err := apollo.Common().GetViper("douyin.yaml")
	if err != nil {
		return err
	}

	// 初始化域名信息
	domain, err = apollo.Common().GetDomain()
	if err != nil {
		return err
	}

	// 初始化JWT设置
	jwt.Settings(jwt.WithFlattenAudience(true))

	// 初始化签发者信息
	issuer = v.GetString("JWT.Issuer")

	// 初始化签名密钥
	signingKey, err = jwk.FromRaw([]byte(v.GetString("JWT.SigningKey")))
	if err != nil {
		return err
	}

	// 初始化加密密钥
	encryptKey, err = jwk.FromRaw([]byte(v.GetString("JWT.EncryptKey")))
	if err != nil {
		return err
	}

	return nil
}

// GenerateToken 生成授权令牌
func GenerateToken(subject string, audience string, scope string) (*auth2.Token, errx.Error) {
	// 解析授权令牌有效时间
	accessTokenValidityTime, _ := time.ParseDuration(cast.ToString(auth2.ClientDetails[audience].AccessTokenValidityTime))

	now := time.Now()

	// 构建授权令牌
	accessToken, err := jwt.NewBuilder().
		Issuer(issuer).
		IssuedAt(now).
		Subject(subject).
		Audience([]string{audience}).
		NotBefore(now.Truncate(time.Second)).
		Expiration(now.Add(accessTokenValidityTime)).
		Claim("scope", scope).
		Build()
	if err != nil {
		return nil, errTokenBuild
	}

	// 将授权令牌进行加密签名
	serializedAccessToken, err := jwt.NewSerializer().
		Encrypt(jwt.WithKey(jwa.PBES2_HS512_A256KW, encryptKey)).
		Sign(jwt.WithKey(jwa.HS256, signingKey)).
		Serialize(accessToken)
	if err != nil {
		return nil, errTokenSerialize
	}

	accessTokenValue := string(serializedAccessToken)

	return &auth2.Token{
		TokenValue: accessTokenValue,
		ExpiresAt:  accessToken.Expiration().Unix(),
	}, nil
}

// ParseToken 解析令牌
func ParseToken(tokenValue string) (string, errx.Error) {
	var tokenBytes []byte

	// 验证令牌是否有效
	payload, err := jws.Verify([]byte(tokenValue), jws.WithKey(jwa.HS256, signingKey))
	if err != nil {
		return "", errInvalidSignature
	}

	// 将令牌载荷解密
	decrypted, err := jwe.Decrypt(payload, jwe.WithKey(jwa.PBES2_HS512_A256KW, encryptKey))
	if err != nil {
		return "", errInvalidKey
	}

	tokenBytes = decrypted

	payloadJson := gjson.ParseBytes(tokenBytes)

	if time.Now().Unix() <= payloadJson.Get("nbf").Int() {
		return "", errTokenNotValidYet
	}

	if time.Now().Unix() > payloadJson.Get("exp").Int() {
		return "", errTokenExpired
	}

	return string(tokenBytes), nil
}
