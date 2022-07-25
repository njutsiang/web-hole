package app

// 消息队列：对 ResponseChanMap 的操作
var ResponseChanMapActionChan = make(chan ResponseChanMapAction, 1000)

// 对 ResponseChanMap 的一次操作
type ResponseChanMapAction struct {
	Name         string // add、del、write
	Request      Request
	RequestId    string
	ResponseChan chan Response
	Response     Response
}