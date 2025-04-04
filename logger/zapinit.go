package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(topic string, level zapcore.Level) {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	kafkaCore := &KafkaCore{
		Level: level,
		Topic: topic,
	}

	Log = zap.New(zapcore.NewTee(consoleCore, kafkaCore))
}
