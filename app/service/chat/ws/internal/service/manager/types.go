package manager

import (
	"douyin/app/service/chat/ws/internal/client"
	"sync"
)

type ClientManager struct {
	Clients        map[*client.Client]struct{} // 全部的连接
	ClientsLock    sync.RWMutex                // 读写锁
	Users          map[int64]*client.Client    // 登录的用户 userId
	UserLock       sync.RWMutex                // 读写锁
	RegisterChan   chan *client.Client         // 连接连接处理
	LoginChan      chan *client.Login          // 用户登录处理
	UnregisterChan chan *client.Client         // 断开连接处理程序
	BroadcastChan  chan []byte                 // 广播 向全体用户发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:        make(map[*client.Client]struct{}),
		Users:          make(map[int64]*client.Client),
		RegisterChan:   make(chan *client.Client, 1000),
		LoginChan:      make(chan *client.Login, 1000),
		UnregisterChan: make(chan *client.Client, 1000),
		BroadcastChan:  make(chan []byte, 1000),
	}

	return
}
