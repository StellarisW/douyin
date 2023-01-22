package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/rpc/sys/internal/config"
	"douyin/app/service/user/rpc/sys/internal/model/profile"
	"douyin/app/service/user/rpc/sys/internal/model/relation"
	"douyin/app/service/user/rpc/sys/internal/model/sign"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	SignModel     sign.Model
	ProfileModel  profile.Model
	RelationModel relation.Model

	Rdb *redis.ClusterClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	idGenerator, err := apollo.Database().NewIdGenerator("user.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetIdGenerator, zap.Error(err))
	}

	db, err := apollo.Database().GetMysqlGormDB()
	if err != nil {
		log.Logger.Fatal(errx.InitMysql, zap.Error(err))
	}

	rdb, err := apollo.Database().GetRedisClusterClient()
	if err != nil {
		log.Logger.Fatal(errx.InitRedis, zap.Error(err))
	}

	authViper, err := apollo.Common().GetViper("auth.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetViper, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		SignModel: sign.NewModel(
			authViper,
			idGenerator,
			db,
			rdb,
		),
		ProfileModel: profile.NewModel(
			db,
			rdb,
		),
		RelationModel: relation.NewModel(
			db,
			rdb,
		),

		Rdb: rdb,
	}
}
