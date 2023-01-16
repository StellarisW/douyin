package chat

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/service/chat/rpc/sys/pb"
	"douyin/app/service/chat/rpc/sys/sys"
	"douyin/app/service/mq/nsq/producer/chat"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type Handler struct {
	ChatSysRpcClient sys.Sys
}

func (m *Handler) HandleMessage(nsqMsg *nsq.Message) error {
	msg := &chat.Message{}

	err := json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonUnmarshal, err)
	}

	if msg.ActionType == 1 {
		rpcRes, _ := m.ChatSysRpcClient.StoreMessage(context.Background(), &pb.StoreMessageReq{
			SrcUserId: msg.SrcUserId,
			DstUserId: msg.DstUserId,
			Content:   msg.Content,
		})
		if rpcRes.StatusCode != 0 {
			return fmt.Errorf("%s, code: %d", rpcRes.StatusMsg, rpcRes.StatusCode)
		}
	} else {
		return fmt.Errorf("action type not supported")
	}

	return nil
}
