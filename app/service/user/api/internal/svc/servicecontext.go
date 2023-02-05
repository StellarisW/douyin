package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/api/internal/config"
	"douyin/app/service/user/rpc/sys/sys"
	"github.com/dlclark/regexp2"
	"github.com/go-redis/redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	CORSMiddleware    rest.Middleware
	JWTAuthMiddleware rest.Middleware

	SysRpcClient sys.Sys

	Rdb *redis.ClusterClient

	Regexp *Regexp
}

type Regexp struct {
	UsernameReg *regexp2.Regexp
	PasswordReg *regexp2.Regexp
}

func NewServiceContext(c config.Config) *ServiceContext {
	corsMiddleware, err := apollo.Middleware().NewCORSMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize corsMiddleware failed.", zap.Error(err))
	}

	JWTAuthMiddleware, err := apollo.Middleware().NewJWTAuthMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize JWTAuthMiddleware failed.", zap.Error(err))
	}

	rdb, err := apollo.Database().GetRedisClusterClient()
	if err != nil {
		log.Logger.Fatal(errx.InitRedis, zap.Error(err))
	}

	// 4到32位(字母,数字,下划线,减号)
	usernameReg := regexp2.MustCompile(`^[a-zA-Z0-9_-]{4,32}$`, regexp2.None)
	// 5-32位(包括至少1个大写字母,1个小写字母,1个数字)
	passwordRes := regexp2.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{5,32}$`, regexp2.None)

	return &ServiceContext{
		Config: c,

		CORSMiddleware:    corsMiddleware.Handle,
		JWTAuthMiddleware: JWTAuthMiddleware.Handle,

		SysRpcClient: sys.NewSys(zrpc.MustNewClient(c.SysRpcClientConf)),

		Rdb: rdb,

		Regexp: &Regexp{
			UsernameReg: usernameReg,
			PasswordReg: passwordRes,
		},
	}
}
