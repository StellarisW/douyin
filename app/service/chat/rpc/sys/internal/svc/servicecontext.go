package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/chat/rpc/sys/internal/config"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	IdGenerator *idgen.DefaultIdGenerator
	Db          *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	idGenerator, err := apollo.Database().NewIdGenerator("chat.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetIdGenerator, zap.Error(err))
	}

	db, err := apollo.Database().GetMysqlGormDB()
	if err != nil {
		log.Logger.Fatal(errx.InitMysql, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		IdGenerator: idGenerator,
		Db:          db,
	}
}
