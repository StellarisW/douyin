package manager

import (
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/chat/api/internal/ws/internal/client"
	"douyin/app/service/chat/api/internal/ws/internal/model/request"
	"douyin/app/service/chat/api/internal/ws/internal/model/response"
	"encoding/json"
	"go.uber.org/zap"
)

func ProcessData(client *client.Client, msg []byte) {
	log.Logger.Debug("processing data...", zap.String("addr", client.Addr))

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Debug("process data complete", zap.Reflect("recover", r))
		}
	}()

	req := &request.Request{}

	err := json.Unmarshal(msg, req)
	if err != nil {
		log.Logger.Error(errx.JsonUnmarshal, zap.Error(err))
		return
	}

	dstClient := Manager.GetUserClient(req.ToUserId)
	if dstClient == nil {
		client.SendMessage([]byte("用户不在线"))
		return
	}

	responseBytes, err := json.Marshal(&response.Response{
		FromUserId: req.UserId,
		MsgContent: req.MsgContent,
	})

	dstClient.SendMessage(responseBytes)

	Manager.StoreMessage(req.UserId, req.ToUserId, req.MsgContent)

	return
}
