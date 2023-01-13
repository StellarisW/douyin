package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/log"
	"douyin/app/service/user/api/internal/config"
	"github.com/zeromicro/go-zero/rest"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	CORSMiddleware    rest.Middleware
	JWTAuthMiddleware rest.Middleware
	CasAuthMiddleware rest.Middleware
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

	return &ServiceContext{
		Config: c,

		CORSMiddleware:    corsMiddleware.Handle,
		JWTAuthMiddleware: JWTAuthMiddleware.Handle,
	}
}
