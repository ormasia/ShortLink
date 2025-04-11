package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局日志对象
var Log *zap.Logger

func init() {
	var err error

	// 创建基本配置
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建Logger
	Log, err = config.Build()
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}

	// 替换全局Logger
	zap.ReplaceGlobals(Log)
}
