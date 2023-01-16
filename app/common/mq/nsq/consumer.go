package nsq

import (
	"douyin/app/common/log"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ConsumerService struct {
	consumer *nsq.Consumer
	channel  chan os.Signal

	logger *zap.SugaredLogger
}

func (m *ConsumerService) Start() {
	err := m.consumer.ConnectToNSQLookupds(mustGetNSQLookupAddrs())
	if err != nil {
		m.logger.Fatalf("start nsq consumer service failed, err: %v", err)
	}
	signal.Notify(m.channel, syscall.SIGINT)
	<-m.channel
}

func (m *ConsumerService) Stop() {
	close(m.channel)
}

func NewConsumerService(topic string, channel string, handler nsq.Handler) (service *ConsumerService, err error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	config.LookupdPollInterval = 15 * time.Second
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	logger := NewLogger(log.Logger.Sugar())
	for i := 0; i <= nsq.LogLevelMax; i++ {
		consumer.SetLogger(logger, nsq.LogLevel(i))
	}

	consumer.AddHandler(handler)

	return &ConsumerService{consumer: consumer, channel: make(chan os.Signal), logger: log.Logger.Sugar()}, nil
}
