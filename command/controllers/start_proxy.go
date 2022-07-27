package controllers

import (
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/logics"
	"os"
	"os/signal"
)

// 启动 Proxy：Websocket 客户端、HTTP 代理
func StartProxy() {
	app.Log.Info("监听进程信号")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for num := 0; num < app.Config.Proxy.WebsocketNum; num++ {
		// 初始化消息队列：待回复的消息
		app.ReplyMessageQueues = append(app.ReplyMessageQueues, make(chan app.ReplyMessage, 1000))

		// 连接到 Frontend
		app.ProxyWebsockets = append(app.ProxyWebsockets, logics.ConnectFrontend(num))

		// 发送心跳
		go logics.SendHeartbeat(num)

		// 读取请求
		go logics.ReadRequest(num)

		// 消费消息队列：待回复的消息
		go logics.ConsumeReplyMessageQueue(num)
	}

	// 当进程停止时，关闭连接
	<- interrupt
	app.Log.Info("进程停止")
	logics.DisconnectFrontend()
}