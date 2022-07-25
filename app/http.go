package app

// 接收响应的通道集合
var ResponseChanMap = map[string]chan Response{}

// 待发送到代理服务器的消息队列
var RequestMessageChan = make(chan []byte, 1000)

// 请求
type Request struct {
	Id      string
	Method  string
	Uri     string
	Header  map[string][]string
	Body    []byte
	Timeout int
}

// 响应
type Response struct {
	RequestId  string
	StatusCode int
	Header     map[string][]string
	Body       []byte
}
