package mq

import (
	"log"

	"github.com/IBM/sarama"
)

var producer sarama.AsyncProducer

func InitKafka(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500

	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return err
	}
	producer = p

	go func() {
		for err := range producer.Errors() {
			log.Println("Kafka 异步日志发送失败:", err)
		}
	}()

	return nil
}

func SendLog(topic, message string) {
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
}
