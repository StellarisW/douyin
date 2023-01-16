package nsq

import (
	"douyin/app/common/log"
	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func NewProducer() (*nsq.Producer, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	producer, err = nsq.NewProducer(mustGetNSQDAddr(), config)
	if err != nil {
		return nil, err
	}

	logger := NewLogger(log.Logger.Sugar())
	for i := 0; i <= nsq.LogLevelMax; i++ {
		producer.SetLogger(logger, nsq.LogLevel(i))
	}
	producer.SetLoggerLevel(nsq.LogLevelInfo)

	err = producer.Ping()
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func GetProducer() (*nsq.Producer, error) {
	if producer == nil {
		return NewProducer()
	}
	return producer, nil
}
