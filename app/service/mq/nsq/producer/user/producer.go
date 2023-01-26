package user

import (
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/internal/consts"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type RelationMessage struct {
	SrcUserId  int64  `json:"src_user_id"`
	DstUserId  int64  `json:"dst_user_id"`
	ActionType uint32 `json:"action_type"`
}

func Relation(producer *nsq.Producer, rawMessage RelationMessage) error {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonMarshal, err)
	}

	err = producer.Publish(consts.TopicUserRelation, message)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.NsqPublish, err)
	}

	return nil
}
