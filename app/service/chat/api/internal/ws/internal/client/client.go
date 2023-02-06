package client

import (
	"douyin/app/common/log"
	"douyin/app/service/chat/api/internal/ws/internal/service/manager"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
)

const (
	heartbeatExpirationTime = 6 * 60
)

func (c *Client) GetUserId() int64 {
	return c.UserId
}

func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Debug("write stop",
				zap.String("stack", string(debug.Stack())),
				zap.Reflect("recover", r),
			)
		}
	}()

	defer func() {
		log.Logger.Debug("close read channel")
		close(c.Send)
	}()

	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			log.Logger.Error("read client data failed", zap.Error(err), zap.String("addr", c.Addr))
			return
		}

		log.Logger.Debug("process client data", zap.String("msg", string(msg)))

		manager.ProcessData(c, msg)
	}
}

func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Debug("write stop",
				zap.String("stack", string(debug.Stack())),
				zap.Reflect("recover", r),
			)
		}
	}()

	defer func() {
		manager.Manager.UnregisterChan <- c
		_ = c.Socket.Close()
		log.Logger.Debug("client send data defer")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Logger.Error("client send data failed", zap.String("addr", c.Addr))
				return
			}

			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) SendMessage(msg []byte) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Debug("send msg stop",
				zap.String("stack", string(debug.Stack())),
				zap.Reflect("recover", r),
			)
		}
	}()

	c.Send <- msg
}

func (c *Client) close() {
	close(c.Send)
}

func (c *Client) Login(userId int64, loginTime int64) {
	c.UserId = userId
	c.LoginTime = loginTime
	c.HeartBeat(loginTime)
}

func (c *Client) HeartBeat(currentTime int64) {
	c.HeartbeatTime = currentTime
}

func (c *Client) IsHeartBeatTimeout(currentTime int64) bool {
	return c.HeartbeatTime+heartbeatExpirationTime <= currentTime
}

func (c *Client) IsLogin() bool {
	return c.UserId != 0
}
