package svc

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/video/api/internal/config"
	"douyin/app/service/video/rpc/sys/sys"
	"github.com/StellarisW/go-sensitive"
	"github.com/minio/minio-go/v7"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config config.Config

	JWTAuthMiddleware rest.Middleware
	CORSMiddleware    rest.Middleware

	SysRpcClient sys.Sys

	Filter *sensitive.Manager

	IdGenerator *idgen.DefaultIdGenerator

	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	filterManager := sensitive.NewFilter(
		sensitive.StoreOption{
			Type: sensitive.StoreMemory,
		},
		sensitive.FilterOption{
			Type: sensitive.FilterDfa,
		})

	err := filterManager.GetStore().LoadDictPath("./manifest/config/dict/default_dict.txt")
	if err != nil {
		log.Logger.Fatal("initialize filter failed", zap.Error(err))
	}

	idGenerator, err := apollo.Database().NewIdGenerator("video.yaml")
	if err != nil {
		log.Logger.Fatal(errx.GetIdGenerator, zap.Error(err))
	}

	corsMiddleware, err := apollo.Middleware().NewCORSMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize corsMiddleware failed.", zap.Error(err))
	}

	JWTAuthMiddleware, err := apollo.Middleware().NewJWTAuthMiddleware()
	if err != nil {
		log.Logger.Fatal("initialize JWTAuthMiddleware failed.", zap.Error(err))
	}

	minioClient, err := apollo.Database().NewMinioClient()
	if err != nil {
		log.Logger.Fatal(errx.InitMinio, zap.Error(err))
	}

	return &ServiceContext{
		Config: c,

		CORSMiddleware:    corsMiddleware.Handle,
		JWTAuthMiddleware: JWTAuthMiddleware.Handle,

		SysRpcClient: sys.NewSys(
			zrpc.MustNewClient(
				c.SysRpcClientConf,
				zrpc.WithDialOption(
					grpc.WithDefaultCallOptions(
						grpc.MaxCallRecvMsgSize(256<<32), // 最大上传 256MB 视频
					),
				),
			),
		),

		Filter: filterManager,

		IdGenerator: idGenerator,

		MinioClient: minioClient,
	}
}
