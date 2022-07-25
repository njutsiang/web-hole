package log

import (
	"fmt"
)

// 日志输出到控制台
type ExportConsole struct {}

func (the *ExportConsole) Write(message string) {
	fmt.Println(message)
}