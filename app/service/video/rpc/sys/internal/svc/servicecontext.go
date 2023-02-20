package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	usersys "douyin/app/service/user/rpc/sys/sys"
	"douyin/app/service/video/rpc/sys/internal/config"
	"douyin/app/service/video/rpc/sys/internal/model/crud"
	"douyin/app/service/video/rpc/sys/internal/model/info"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config config.Config

	CrudModel crud.Model
	InfoModel info.Model
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

	userSysRpcClient := usersys.NewSys(zrpc.MustNewClient(c.UserSysRpcClientConf))

	return &ServiceContext{
		Config: c,

		CrudModel: crud.NewModel(idGenerator, db, rdb, userSysRpcClient),
		InfoModel: info.NewModel(db, rdb, userSysRpcClient),
	}
}
