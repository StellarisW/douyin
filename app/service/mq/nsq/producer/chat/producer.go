package chat

import (
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/internal/consts"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	SrcUserId  int64  `json:"src_user_id"`
	DstUserId  int64  `json:"dst_user_id"`
	ActionType uint32 `json:"action_type"`
	Content    string `json:"content"`
}

func Chat(producer *nsq.Producer, rawMessage Message) error {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonMarshal, err)
	}

	err = producer.Publish(consts.ChannelChat, message)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.NsqPublish, err)
	}

	return nil
}
