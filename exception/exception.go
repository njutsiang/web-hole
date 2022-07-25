package exception

// 定义异常
type Exception struct {
	code int
	message string
}

// 异常错误码
func (the *Exception) GetCode() int {
	return the.code
}

// 异常提示语
func (the *Exception) GetMessage() string {
	return the.message
}

// 实现 error
func (the *Exception) Error() string {
	return the.message
}