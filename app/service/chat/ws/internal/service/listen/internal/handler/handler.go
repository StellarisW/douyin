package handler

import (
	"douyin/app/common/log"
	"douyin/app/service/chat/ws/internal/client"
	"douyin/app/service/chat/ws/internal/service/manager"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func WsPage(w http.ResponseWriter, r *http.Request) {
	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		log.Logger.Debug("upgrade protocol",
			zap.Strings("user-agent", r.Header["User-Agent"]),
			zap.Strings("referer", r.Header["Referer"]),
		)
		return true
	}}).Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Logger.Debug("build connection", zap.String("addr", conn.RemoteAddr().String()))

	currentTime := time.Now().Unix()

	c := client.NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go c.Read()
	go c.Write()

	// 用户连接事件
	manager.Manager.RegisterChan <- c
}
