package exception

// 抛出异常
func Throw(code int, messages... string) {
	message := ""
	if len(messages) >= 1 && messages[0] != "" {
		message = messages[0]
	} else {
		message = GetMessage(code)
	}
	panic(&Exception{
		code: code,
		message: message,
	})
}

// 获取错误提示语
func GetMessage(code int) string {
	message := messages[code]
	if message == "" {
		message = "未知错误"
	}
	return message
}