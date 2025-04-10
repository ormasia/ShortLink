package mq

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type LogConsumer struct{}

func (LogConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (LogConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (LogConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf(" 收到日志: [%s] %s", msg.Topic, string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

func StartLogConsumer(brokers []string, groupID string, topics []string) error {
	config := sarama.NewConfig()
	//仅在该 消费者组是首次消费某个 topic 的时候，从最早的消息（offset = 0）开始消费。
	//一旦消费者组有了历史记录（offset 已存在），这个设置就不再生效。
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return fmt.Errorf("无法创建消费者组: %w", err)
	}

	go func() {
		for err := range consumer.Errors() {
			log.Println("消费出错：", err)
		}
	}()

	log.Println(" 启动 Kafka 日志消费者")
	for {
		err := consumer.Consume(context.TODO(), topics, LogConsumer{})
		if err != nil {
			log.Println("消费循环错误：", err)
		}
	}
}

/*
1,看这个 groupID 有没有记录（在 __consumer_offsets topic 中）

2,没有 → 就按你设置的 OffsetOldest 从头读

3,你 session.MarkMessage(...)，Kafka 会记录 offset

4,下次你用同一个 groupID 来消费，它会从你上次消费完的地方继续
*/
