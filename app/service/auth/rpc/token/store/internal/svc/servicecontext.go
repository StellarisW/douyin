package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/auth/rpc/token/store/internal/config"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	Rdb *redis.ClusterClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, err := apollo.Database().GetRedisClusterClient()
	if err != nil {
		log.Logger.Fatal(errx.InitRedis, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		Rdb: rdb,
	}
}
