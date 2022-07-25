package log

import (
	"encoding/json"
	"github.com/njutsiang/web-hole/utils"
	"time"
)

// 日志级别
const (
	LevelInfo = 1
	LevelWarning = 2
	LevelError = 3
)

// 日志组件
type Logger struct {
	Level int
	Exports []Export
}

// 日志输出接口
type Export interface {
	Write(string)
}

// 日志数据格式
type Message struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

// 信息日志
func (the *Logger) Info(message string) {
	if the.Level == LevelInfo {
		for _, export := range the.Exports {
			export.Write(Format("info", message))
		}
	}
}

// 警告日志
func (the *Logger) Warning(message string) {
	if the.Level <= LevelWarning {
		for _, export := range the.Exports {
			export.Write(Format("warning", message))
		}
	}
}

// 错误日志
func (the *Logger) Error(message string) {
	for _, export := range the.Exports {
		export.Write(Format("error", message))
	}
}

// 格式化日志内容
func Format(level string, message string) string {
	messageJson, _ := json.Marshal(Message{
		Time: utils.StrPadRight(time.Now().Format("2006-01-02 15:04:05.999"), "0", 23),
		Level: level,
		Message: message,
	})
	return string(messageJson)
}