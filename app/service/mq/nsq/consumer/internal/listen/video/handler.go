package video

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/producer/video"
	"douyin/app/service/video/rpc/sys/pb"
	"douyin/app/service/video/rpc/sys/sys"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type FavoriteHandler struct {
	VideoSysRpcClient sys.Sys
}

func (m *FavoriteHandler) HandleMessage(nsqMsg *nsq.Message) error {
	msg := &video.FavoriteMessage{}

	err := json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonUnmarshal, err)
	}

	rpcRes, _ := m.VideoSysRpcClient.Favorite(context.Background(), &pb.FavoriteReq{
		UserId:     msg.UserId,
		VideoId:    msg.VideoId,
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

type CommentHandler struct {
	VideoSysRpcClient sys.Sys
}

func (m *CommentHandler) HandleMessage(nsqMsg *nsq.Message) error {
	msg := &video.CommentMessage{}

	err := json.Unmarshal(nsqMsg.Body, &msg)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonUnmarshal, err)
	}

	rpcRes, _ := m.VideoSysRpcClient.ManageComment(context.Background(), &pb.ManageCommentReq{
		UserId:      msg.UserId,
		VideoId:     msg.VideoId,
		ActionType:  msg.ActionType,
		CommentText: msg.CommentText,
		CommentId:   msg.CommentId,
	})
	if rpcRes == nil {
		return fmt.Errorf(errx.RequestRpcReceive)
	}
	if rpcRes.StatusCode != 0 {
		return fmt.Errorf("%s, code: %d", rpcRes.StatusMsg, rpcRes.StatusCode)
	}

	return nil
}
