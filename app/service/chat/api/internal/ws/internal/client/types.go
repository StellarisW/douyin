package client

import "github.com/gorilla/websocket"

// Login 用户登录
type Login struct {
	UserId int64
	Client *Client
}

// Client 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	UserId        int64           // 用户Id
	FirstTime     int64           // 首次连接事件
	HeartbeatTime int64           // 用户上次心跳时间
	LoginTime     int64           // 登录时间 登录以后才有
}

// NewClient 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime int64) (client *Client) {
	return &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
}
