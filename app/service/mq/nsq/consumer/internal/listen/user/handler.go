package user

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/producer/user"
	"douyin/app/service/user/rpc/sys/pb"
	"douyin/app/service/user/rpc/sys/sys"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type RelationHandler struct {
	UserSysRpcClient sys.Sys
}

func (m *RelationHandler) HandleMessage(nsqMsg *nsq.Message) error {
	msg := &user.RelationMessage{}

	err := json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonUnmarshal, err)
	}

	rpcRes, _ := m.UserSysRpcClient.Relation(context.Background(), &pb.RelationReq{
		SrcUserId:  msg.SrcUserId,
		DstUserId:  msg.DstUserId,
		ActionType: msg.ActionType,
	})
	if rpcRes == nil {
		return fmt.Errorf(errx.RequestRpcReceive)
	}
	if rpcRes.StatusCode != 0 {
		return fmt.Errorf("%s, code: %d", rpcRes.StatusMsg, rpcRes.StatusCode)
	}

	return nil
}
