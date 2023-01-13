package main

import (
	"go.uber.org/zap"
	apollo "douyin/app/common/config"
	"douyin/app/common/log"
	"douyin/app/common/errx"
	"douyin/utils"

	{{.imports}}

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: 填写 rpc 名称
const serviceName = ".rpc.{{.serviceName}}"

var c config.Config

func main() {
	serviceFullName := utils.GetServiceFullName(serviceName)

	// 初始化日志管理器
	err := log.InitLogger(serviceFullName + "-data")
	if err != nil {
		panic("initialize logger failed")
	}

	logWriter := log.GetLogxZapWriter(log.Logger)

	logx.MustSetup(log.GetLogxConfig(serviceFullName, "error"))
	logx.SetWriter(logWriter)

	// 初始化配置管理器
	err = apollo.InitClient()
	if err != nil {
		log.Logger.Fatal(errx.InitAgolloClient, zap.Error(err))
	}
	
	// 初始化微服务设置
	namespace, serviceType, serviceSingleName := utils.GetServiceDetails(serviceName)
	err = apollo.Common().UnmarshalServiceConfig(namespace, serviceType, serviceSingleName, &c)
	if err != nil {
		log.Logger.Fatal(errx.UnmarshalServiceConfig, zap.Error(err))
	}

	ctx := svc.NewServiceContext(c)
	
	// 启动微服务服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
{{range .serviceNames}}       {{.Pkg}}.Register{{.Service}}Server(grpcServer, {{.ServerPkg}}.New{{.Service}}Server(ctx))
{{end}}
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	
	// 注册服务到consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)
	
	defer s.Stop()

	log.Logger.Info("starting rpc server...", zap.String("addr", c.ListenOn))
	s.Start()
}
