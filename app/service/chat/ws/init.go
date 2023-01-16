package ws

import (
	apollo "douyin/app/common/config"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/utils"
	"go.uber.org/zap"
)

const serviceName = "chat.websocket"

func main() {
	serviceFullName := utils.GetServiceFullName(serviceName)

	// 初始化日志管理器
	err := log.InitLogger(serviceFullName + "-data")
	if err != nil {
		panic("initialize logger failed")
	}

	// 初始化配置管理器
	err = apollo.InitClient()
	if err != nil {
		log.Logger.Fatal(errx.InitAgolloClient, zap.Error(err))
	}

}
