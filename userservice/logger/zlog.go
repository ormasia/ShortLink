package logger

import (
	"shortLink/userservice/mq"

	"go.uber.org/zap/zapcore"
)

type KafkaCore struct {
	Level zapcore.Level
	Topic string
}

func (k *KafkaCore) Enabled(lvl zapcore.Level) bool {
	return lvl >= k.Level
}

func (k *KafkaCore) With(fields []zapcore.Field) zapcore.Core {
	return k
}

func (k *KafkaCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if k.Enabled(entry.Level) {
		return ce.AddCore(entry, k)
	}
	return ce
}

func (k *KafkaCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	mq.SendLog(k.Topic, entry.Message)
	return nil
}

func (k *KafkaCore) Sync() error {
	return nil
}
