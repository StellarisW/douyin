## 需求分析

该项目的接口除了用户注册、登录，及视频流接口外皆需要在query或form表单里面提供用户鉴权token，

所以该模块需要有两个功能：

- 在用户登录后给用户 **提供** 鉴权token，

- 根据在其他接口的JWT中间件中发送的请求中提供的token，**校验** token有效性及 **解析** 其中的用户信息

## 架构设计

<div style="text-align: center">
<img src="image/auth-module.png">
</div>

### 令牌的颁发与校验

业务逻辑主要通过 **令牌颁发器接口** 和 **令牌服务接口** 实现

#### 令牌颁发器

```go
type (
	// TokenGranter 令牌颁发器接口
	TokenGranter interface {
		Grant(ctx context.Context, grantType string, auth string, obj string) (*auth.Token, errx.Error)
	}

	// ComposeTokenGranter 集成令牌颁发器
	ComposeTokenGranter struct {
		TokenGrantDict map[string]TokenGranter
	}

	// AuthorizationTokenGranter 通过认证令牌的颁发器结构体
	AuthorizationTokenGranter struct {
		SupportGrantType string
		ClientSecret     map[string]string
		TokenService     TokenService
	}
)
```

该接口只有一个方法 `Grant()` ，用来颁发令牌，

目前只有 **Authorization** 一种认证方式，该接口的好处是具有良好的扩展性，方便以后拓展 **Oauth2 登录**等多种登录认证方式

参数说明：

- **granType**：颁发类型

    目前有以下类型：

    - **Authorization**：根据Basic Auth头颁发令牌

- **auth**：认证参数

- **obj**：颁发对象

#### 令牌服务

```go
type (
	// TokenService 令牌服务接口
	TokenService interface {
		GenerateToken(ctx context.Context, subject string, audience string) (*auth.Token, errx.Error)
		ReadToken(ctx context.Context, tokenValue string) (string, errx.Error)
	}

	// DefaultTokenService 默认令牌服务结构体
	DefaultTokenService struct {
	}

	// RpcTokenService rpc 令牌服务结构体
	RpcTokenService struct {
		TokenEnhancerClient tokenenhancer.TokenEnhancer
	}
)
```

该接口共有两种方法：`GenerateToken()` 用来颁发令牌，`ReadToken()`用来读取令牌

该接口通过调用 `TokenEnhancerClient` Rpc服务实现业务逻辑

### 令牌设计

#### jwt 设计

本项目 jwt 通过 [jwx](https://github.com/lestrrat-go/jwx) 实现

- **header**

    Header 部分是一个 JSON 对象，描述 JWT 的元数据

    ```json
    {
      "typ": "JWT", // 令牌类型
      "alg": "HS256" // 加密算法
    }
    ```

- **payload**

    Payload 部分也是一个 JSON 对象，用来存放实际需要传递的数据

    因为 JWT 默认是不加密的，所以本项目使用 `PBES2_HS512_A256KW` 加密算法将 payload 加密后再进行签名操作

    ```json
    {
      "iss": "douyin.xxx.com", // 签发人
      "iat": 1516239022, // 签发时间
      "sub": "StellarisW", // 主题
      "aud": "douyin.xxx.com", // 受众
      "nbf": 1516239022, // 生效时间
      "exp": 1516339022 // 过期时间
    }
    ```

- **signature**

    Signature 部分是对前两部分的签名，防止数据篡改。

    首先，需要指定一个密钥（secret）。这个密钥只有服务器才知道，不能泄露给用户。

    然后，使用 Header 里面指定的签名算法（本项目使用 HMAC SHA256），按照下面的公式产生签名

    算出签名以后，把 Header、Payload、Signature 三个部分拼成一个字符串，每个部分之间用"点"（`.`）分隔，就可以返回给用户

    ```
    HMACSHA256(
      base64UrlEncode(header) + "." +
      base64UrlEncode(payload), 
    )
    ```

#### 结构体 设计

```go
type Token struct {
	TokenValue string `json:"token_value,omitempty"` // 令牌值
	ExpiresAt  int64  `json:"expires_at,omitempty"`  // 过期时间 (unix)
}
```

#### 生成令牌方法

```GO
// GenerateToken 生成授权令牌
func GenerateToken(subject string, audience string, scope string) (*auth.Token, errx.Error) {
	// 解析授权令牌有效时间
	accessTokenValidityTime, _ := time.ParseDuration(cast.ToString(auth.ClientDetails[audience].AccessTokenValidityTime))

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

	return &auth.Token{
		TokenValue: accessTokenValue,
		ExpiresAt:  accessToken.Expiration().Unix(),
	}, nil
}
```

#### 解析令牌方法

```go
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
```

### 令牌存储

颁发令牌时先检测redis中有无现存的令牌，这样可以避免重复颁发令牌，

此外，后续可以根据该功能实现的令牌的黑名单功能

#### 缓存设计

- key：`auth:oauth2_token:{obj}`

- value：

    ```json
    {
        "token_value": "eyJhbGciOiJIUzI1NiIsImN0eSI6IkpXVCJ9.ZXlKaGJHY2lPaUpRUWtWVE1pMUlVelV4TWl0Qk1qVTJTMWNpTENKbGJtTWlPaUpCTWpVMlIwTk5JaXdpY0RKaklqb3hNREF3TUN3aWNESnpJam9pZUVGclYwcEhkR3BwTVU1TlZHOURVMk16VW05c2FtWkZTSEF3UkVnNU9EZFVTekJFVVRZeGVWZ3ljeUlzSW5SNWNDSTZJa3BYVkNKOS5WMURrR2daWGhzbTJrTDg5bVFpNmZpTk1qeUdpUEVLMG90RU82a2xCQ3BaLWNrMkJFRTZlNUEuNXBHLW1rSXlmQ0tQSTh5ai4wUU8xSkRja19ka2I4YVVSUVU3aE5fV1NlUmJuSmJscldPSERTWWVXLTRka2JrcTFobHZZSFhlVFhDNEhPUUVDak1FSEdwYmM2aXRKT2ZJWDhtQTBZLWVQQVJwX2ZuVXVseDhlSl9EOVo0NnlXOVRlcG94WGtlVHpIQ0VZMm9JWmRmRGVEbUNMY0FtT1FYVjZMd3oxZTY1WmhURmhPdTNSVDJJcExnaEczaTBwU0t1djBkazFySEJVMGgxQ3BPWWRvaGtVMEEuV1N1UDVlaXhES0lYRXdaQlNHMEc4Zw.XPNcH7ItFYahj6cd7YegdzYGThvZ3aiqpcE-m74Y2Ls",
        "expires_at": 1674708117
    }
    ```

- TTL：7d