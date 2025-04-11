package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log 提供简单的日志接口，替代原来的zap日志
var Log = &SimpleLogger{}

// SimpleLogger 简单的日志实现
type SimpleLogger struct{}

// Debug 输出调试日志
func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	l.log("DEBUG", msg, fields...)
}

// Info 输出信息日志
func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

// Warn 输出警告日志
func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	l.log("WARN", msg, fields...)
}

// Error 输出错误日志
func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

// Fatal 输出致命错误日志并退出程序
func (l *SimpleLogger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

// log 内部日志输出方法
func (l *SimpleLogger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("%s [%s] %s", timestamp, level, msg)

	// 添加字段信息
	if len(fields) > 0 {
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMsg += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}

	log.Println(logMsg)
}
