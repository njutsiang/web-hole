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

	// 消费消息队列：待回复的消息
	go logics.ConsumeReplyMessageChan()

	// 连接到 Frontend
	app.ConnectFrontend()

	// 发送心跳
	go logics.SendHeartbeat()

	// 读取请求
	go logics.ReadRequest()

	// 当进程停止时，关闭连接
	<- interrupt
	app.Log.Info("进程停止")
	app.DisconnectFrontend()
}