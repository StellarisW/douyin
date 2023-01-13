package svc

import (
	"douyin/app/service/video/api/internal/config"
	"douyin/app/service/video/api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config            config.Config
	JWTAuthMiddleware rest.Middleware
	CORSMiddleware    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		JWTAuthMiddleware: middleware.NewJWTAuthMiddleware().Handle,
		CORSMiddleware:    middleware.NewCORSMiddleware().Handle,
	}
}
