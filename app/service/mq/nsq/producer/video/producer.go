package video

import (
	"douyin/app/common/errx"
	"douyin/app/service/mq/nsq/internal/consts"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
)

type FavoriteMessage struct {
	UserId     int64  `json:"user_id"`
	VideoId    int64  `json:"video_id"`
	ActionType uint32 `json:"action_type"`
}

type CommentMessage struct {
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id"`
	ActionType  uint32 `json:"action_type"`
	CommentText string `json:"comment_text"`
	CommentId   int64  `json:"comment_id"`
}

func Favorite(producer *nsq.Producer, rawMessage FavoriteMessage) error {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonMarshal, err)
	}

	err = producer.Publish(consts.TopicVideoFavorite, message)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.NsqPublish, err)
	}

	return nil
}

func Comment(producer *nsq.Producer, rawMessage CommentMessage) error {
	message, err := json.Marshal(rawMessage)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.JsonMarshal, err)
	}

	err = producer.Publish(consts.TopicVideoComment, message)
	if err != nil {
		return fmt.Errorf("%s, err: %v", errx.NsqPublish, err)
	}

	return nil
}
