package app

import (
	"github.com/njutsiang/web-hole/exception"
	"github.com/njutsiang/web-hole/library/log"
	"os"
)

// 日志组件
var Log = &log.Logger{}

// 初始化日志组件
func InitLog() {
	switch Config.Log.Level {
	case "info":
		Log.Level = log.LevelInfo
	case "warning":
		Log.Level = log.LevelWarning
	case "error":
		Log.Level = log.LevelError
	}
	if Config.Log.ExportConsole == 1 {
		Log.Exports = append(Log.Exports, &log.ExportConsole{})
	}
	if Config.Log.ExportFile.Path != "" {
		exportFile := &log.ExportFile{
			Path: Config.Log.ExportFile.Path,
		}
		var err error
		exportFile.File, err = os.OpenFile(exportFile.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			exception.Throw(exception.OpenLogFileFailed, "打开日志文件失败" + err.Error())
		}
		Log.Exports = append(Log.Exports, exportFile)
		go log.ConsumeExportFileFlushChan()
	}
}