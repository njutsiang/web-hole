package exception

// 错误码
var (
	ConnectFrontendFailed = 400000
)

// 错误码描述
var messages = map[int]string{
	ConnectFrontendFailed: "连接到 Frontend 失败",
}
