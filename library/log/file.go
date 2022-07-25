package log

import (
	"bufio"
	"fmt"
	"os"
)

// 消息队列：日志文件刷盘
var ExportFileFlushChan = make(chan ExportFileFlush, 2000)

// 日志输出到文件
type ExportFile struct {
	Path string
	File *os.File
}

// 日志文件刷盘操作
type ExportFileFlush struct {
	ExportFile *ExportFile
	Message string
}

// 写入消息队列
func (the *ExportFile) Write(message string) {
	ExportFileFlushChan <- ExportFileFlush{
		ExportFile: the,
		Message: message,
	}
}

// 执行日志刷盘
func (the *ExportFile) Flush(message string) {
	writer := bufio.NewWriter(the.File)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		fmt.Println(Format("error", "写入日志文件失败 " + err.Error()))
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(Format("error", "日志文件刷盘失败 " + err.Error()))
	}
}

// 消费消息队列：日志文件刷盘
func ConsumeExportFileFlushChan() {
	for exportFileFlush := range ExportFileFlushChan {
		exportFileFlush.ExportFile.Flush(exportFileFlush.Message)
	}
}
