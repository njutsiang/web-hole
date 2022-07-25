package exception

// 错误码
var (
	InitConfigFailed = 400000
	ConnectFrontendFailed = 400001
	OpenLogFileFailed = 400002

)

// 错误码描述
var messages = map[int]string{
	InitConfigFailed: "初始化配置失败",
	ConnectFrontendFailed: "连接到 Frontend 失败",
	OpenLogFileFailed: "打开日志文件失败",
}
