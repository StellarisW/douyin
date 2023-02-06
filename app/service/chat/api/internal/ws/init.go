package ws

import (
	"douyin/app/service/chat/api/internal/ws/internal/service/manager"
	"douyin/app/service/chat/rpc/sys/sys"
)

func Init(RpcClient sys.Sys) {
	manager.Manager = manager.NewClientManager(RpcClient)
}
