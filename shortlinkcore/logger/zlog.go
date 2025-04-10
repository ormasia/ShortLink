package logger

import (
	"encoding/json"
	"shortLink/shortlinkcore/mq"
	"time"

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
	data := make(map[string]interface{})
	data["timestamp"] = entry.Time.Format(time.RFC3339)
	data["level"] = entry.Level.String()
	data["message"] = entry.Message
	data["logger"] = entry.LoggerName
	data["caller"] = entry.Caller.TrimmedPath()
	data["app"] = "shortlink"
	data["env"] = "prod"

	// 把 zap.Field 中的动态字段加入日志结构
	for _, f := range fields {
		data[f.Key] = f.Interface
	}

	// 转为 JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	mq.SendLog(k.Topic, string(jsonBytes))
	return nil
	// mq.SendLog(k.Topic, entry.Message)
	// return nil
}

func (k *KafkaCore) Sync() error {
	return nil
}
