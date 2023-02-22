# 一、概述

本项目中使用到多种组件来保障微服务的稳定和可靠。主要包括

- 中间件
    - JWT 中间件
    - Sentinel 中间件
- NSQ 消息队列
- 统一的日志系统

# 二、服务中间件设计

## 2.1 JWT 中间件

本项目特意设计了 auth 认证模块用于提供用户 token 的查询、颁布和存储，JWT 中间件通过调用 auth 认证模块提供的服务来进行用户 token 的鉴权。

### 2.1.1 令牌设计

本项目 jwt 通过 **[jwx](https://github.com/lestrrat-go/jwx)** 实现。

- Header 部分是一个 JSON 对象，描述 JWT 的元数据

```JSON
{
  "typ": "JWT", // 令牌类型
  "alg": "HS256" // 加密算法
}
```

- Payload 部分也是一个 JSON 对象，用来存放实际需要传递的数据。
    - 因为 JWT 默认是不加密的，所以本项目使用 `PBES2_HS512_A256KW` 加密算法将 payload 加密后再进行签名操作

```JSON
{
  "iss": "douyin.xxx.com", // 签发人
  "iat": 1516239022, // 签发时间
  "sub": "StellarisW", // 主题
  "aud": "douyin.xxx.com", // 受众
  "nbf": 1516239022, // 生效时间
  "exp": 1516339022 // 过期时间
}
```

- Signature 部分是对前两部分的签名，防止数据篡改。
    - 首先，需要指定一个密钥（secret）。这个密钥只有服务器才知道，不能泄露给用户。
    - 然后，使用 Header 里面指定的签名算法（本项目使用 HMAC SHA256），按照下面的公式产生签名
    - 算出签名以后，把 Header、Payload、Signature 三个部分拼成一个字符串，每个部分之间用"点"（`.`）分隔，就可以返回给用户

```JSON
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload), 
)
```

### 2.1.2 令牌结构体设计

```Go
type Token struct {
   TokenValue string `json:"token_value,omitempty"` // 令牌值
   ExpiresAt  int64  `json:"expires_at,omitempty"`  // 过期时间 (unix)
}
```

### 2.1.3 令牌生成和解析

```Go
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

### 2.1.4 令牌的存储

本项目中令牌使用缓存进行存储。

|            key            |     value     | TTL  |
| :-----------------------: | :-----------: | :--: |
| `auth:oauth2_token:{obj}` | `token_value` |  7d  |

颁发令牌时先检测 Redis 中有无现存的令牌，这样可以避免重复颁发令牌。此外，后续可以根据该功能实现的令牌的黑名单功能。

## 2.2 熔断、限流中间件

为更好地保障微服务的各个服务之间的稳定性，本项目对微服务的熔断、限流进行了配置。

### 2.2.1 熔断中间件

本项目使用的 go-zero 微服务框架内置提供了微服务自适应熔断的机制。

熔断器主要是用来保护调用端，调用端在发起请求的时候需要先经过熔断器，而客户端拦截器正好兼具了这个这个功能，所以 go-zero 框架内熔断器是实现在客户端拦截器内。

```Go
func BreakerInterceptor(ctx context.Context, method string, req, reply interface{},
    cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    // 基于请求方法进行熔断
    breakerName := path.Join(cc.Target(), method)
    return breaker.DoWithAcceptable(breakerName, func() error {
        // 真正发起调用
        return invoker(ctx, method, req, reply, cc, opts...)
    // codes.Acceptable判断哪种错误需要加入熔断错误计数
    }, codes.Acceptable)
}
```

其内部实现的逻辑是参考了 [Google Sre过载保护算法](https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101)，该算法的原理如下：

- 请求数量（requests）：调用方发起请求的数量总和
- 请求接受数量（accepts）：被调用方正常处理的请求数量

在正常情况下，这两个值是相等的，随着被调用方服务出现异常开始拒绝请求，accepts 的值开始逐渐小于请求数量requests，这个时候调用方可以继续发送请求，直到 requests = K * accepts。一旦超过这个限制，熔断器就回打开，新的请求会在本地以一定的概率被抛弃直接返回错误。通过修改算法中的 K 值，可以调节熔断器的敏感度。

### 2.2.2 限流中间件

在微服务中，由于任意时间到来的请求往往是随机不可控的，而系统的处理能力是有限的。所以需要根据系统的处理能力对流量进行控制。

本项目中使用阿里巴巴开源的 Sentinel 限流组件设计限流中间件，为每个服务提供了限流配置，保障服务的可靠和稳定。

#### 2.2.2.1 限流规则

Sentinel 是基于限流规则实现的流量控制，一条限流规则主要由下面几个因素组成，可以组合这些元素来实现不同的限流效果：

- `Resource`：资源埋点名，即限流规则的作用对象
- `TokenCalculateStrategy`：限流策略，Direct 表示直接使用`Threshold`作为阈值进行限流，WarmUp 表示根据设定的预热曲线进行限流
- `ControlBehavior`：控制策略，拒绝或匀速排队
- `Threshold`： 周期内流控阈值
- `StatIntervalInMs`：数据统计周期，单位 ms。设为 1000 就是基于 QPS

本项目实际使用中，将每个服务的限流规则定义在服务的配置文件中（由 Apollo 统一管理）。当服务启动时，通过解析获取定义的限流规则来实现限流组件的自动加载。

```Go
func (g *Group) newFlowRules(namespace string) ([]*flow.Rule, error) {
   var rules []*flow.Rule

   err := common.GetGroup().UnmarshalKey(namespace, "Api.Sentinel.Flow", &rules)
   if err != nil {
      return nil, err
   }

   return rules, nil
}
```

#### 2.2.2.2 资源埋点

每个服务都会提供多个应用服务接口，用于接收请求并处理数据，服务的配置文件中针对各个接口设定的限流规则是对当前资源埋点字段的数据进行的捕获。当埋点数据超过阈值，则触发限流。

本项目中统一使用 `req.Method+req.RequestURI` 作为资源埋点信息，在配置文件中的 `resource` 字段中进行同样的定义。

#### 2.2.2.3 限流中间件的加载

本项目在 `common` 目录中提供了限流中间件加载方法，每个服务在启动时调用该方法，即可完成限流中间件加载。

```Go
func NewSentinelMiddleware(entity *config.Entity, rules []*flow.Rule) *SentinelMiddleware {
   err := sentinel.InitWithConfig(entity)
   if err != nil {
      panic("invalid config")
   }

   _, err = flow.LoadRules(rules)
   if err != nil {
      panic("invalid flow rule")
   }

   return &SentinelMiddleware{}
}
```

# 三、消息队列设计

本项目基于 NSQ 搭建了消息队列服务，并针对**用户**、**视频**和**聊天**三个模块设计了对应模块需要的消息生产者和消费者。生产者对数据序列化成消息体，并推送至消息队列。消费者不断从消息队列中取出消息，反序列化之后调用需要的 rpc 服务完成处理任务。

具体来讲，本项目中主要服务使用到消息队列的功能有

- **用户服务**
    - 关注/取消关注
- **视频服务**
    - 点赞/取消点赞
    - 评论操作
- **聊天服务**
    - 消息存储

使用消息队列可以有效提升系统稳定性，同时节约了响应时间，一定程度上提升了用户体验。

# 四、错误信息设计

## 4.1 编码设计

本项目使用全局统一格式的错误码，方便定位问题，便于错误的排除和处理。

错误码总计 32 位，按照如下规则进行设计

| 位置 | 1        | 2-6     | 7        | 8-12    | 13-18   | 19-24   | 25-32   |
| ---- | -------- | ------- | -------- | ------- | ------- | ------- | ------- |
| 说明 | 错误类型 | 系统 id | 服务类型 | 服务 id | 业务 id | 操作 id | 错误 id |

针对以上规则，本项目在 `common/errx` 包中设计了相应的 `Encode` 和 `Decode` 函数实现错误码的编码和解码。

## 4.2 错误信息定义

错误信息数据结构包括错误编码和错误内容。在 `common/errx` 包中提供了如下的接口

```Go
type (
   Error interface {
      Code() uint32
      Error() string
   }
)
```

在实际使用中，搭配日志系统记录错误信息。

# 五、日志系统设计

本项目使用 **[zap](https://pkg.go.dev/go.uber.org/zap)** 作为日志系统。

```Go
func InitLogger(name string) error {
   options := Options{
      Mode:         os.Getenv(douyin.ModeEnvName),
      SavePath:     name + "/log",
      EncoderType:  ConsoleEncoder,
      EncodeLevel:  CapitalLevelEncoder,
      EncodeCaller: FullCallerEncoder,
   }
   return NewLoggerWithOptions(options)
}
```

针对不同的环境为日志系统添加不同的配置。

```Go
func NewLoggerWithOptions(options Options) (err error) {
   // 创建日志保存的文件夹
   err = file.IsNotExistMkDir(options.SavePath)
   if err != nil {
      return err
   }

   dynamicLevel := zap.NewAtomicLevel()
   encoder := getEncoder(options)

   switch options.Mode {
   case douyin.DevMode:
      // 将当前日志等级设置为 Debug
      // 注意日志等级低于设置的等级，日志文件也不分记录
      dynamicLevel.SetLevel(zap.DebugLevel)

      // 调试级别
      debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
         return lev == zap.DebugLevel
      })
      // 日志级别
      infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
         return lev == zap.InfoLevel
      })
      // 警告级别
      warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
         return lev == zap.WarnLevel
      })
      // 错误级别(包含error,panic,fatal级别的日志)
      errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
         return lev >= zap.ErrorLevel
      })
      cores := [...]zapcore.Core{
         zapcore.NewCore(encoder, os.Stdout, dynamicLevel), //控制台输出
         //日志文件输出,按等级归档
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/all/server_all.log", options.SavePath)), zapcore.DebugLevel),
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/debug/server_debug.log", options.SavePath)), debugPriority),
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/info/server_info.log", options.SavePath)), infoPriority),
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/warn/server_warn.log", options.SavePath)), warnPriority),
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/error/server_error.log", options.SavePath)), errorPriority),
      }
      Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
      defer func(zapLogger *zap.Logger) {
         _ = zapLogger.Sync()
      }(Logger)
   case douyin.FatMode, douyin.UatMode:
      // 将当前日志等级设置为 Info
      // 注意日志等级低于设置的等级，日志文件也不分记录
      dynamicLevel.SetLevel(zap.InfoLevel)

      cores := [...]zapcore.Core{
         zapcore.NewCore(encoder, os.Stdout, dynamicLevel), //控制台输出
         zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/server.log", options.SavePath)), dynamicLevel),
      }

      Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
      defer func(zapLogger *zap.Logger) {
         _ = zapLogger.Sync()
      }(Logger)
   case douyin.ProMode:
      dynamicLevel.SetLevel(zap.ErrorLevel)

      Logger = zap.New(zapcore.NewCore(
         encoder,
         getWriteSyncer(fmt.Sprintf("./%s/server.log", options.SavePath)), dynamicLevel),
         zap.AddCaller(),
      )
      defer func(zapLogger *zap.Logger) {
         _ = zapLogger.Sync()
      }(Logger)
   default:
      panic(errx.EmptyProjectMode)
   }

   //设置全局logger
   Logger.Info("initialize logger successfully")
   return nil
```