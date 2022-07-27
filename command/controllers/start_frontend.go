package controllers

import "github.com/njutsiang/web-hole/logics"

// 启动 Frontend
func StartFrontend() {

	// 消费消息队列：对 ResponseChanMap 的操作
	go logics.ConsumeResponseChanMapActionChan()

	// 消费消息队列：发送到代理服务器的消息
	go logics.ConsumeRequestMessageChan()

	// 消费消息队列：添加 Frontend 和 Proxy 的连接
	go logics.ConsumeFrontendWebsocketChan()

	// 消费消息队列：删除 Frontend 和 Proxy 的连接
	go logics.ConsumeFrontendWebsocketDelChan()

	// 启动 HTTP 服务
	go StartHttp()

	// 启动 Websocket 服务
	StartWebsocket()
}
