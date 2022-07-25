package controllers

import (
	"hole/app"
	"hole/logics"
	"log"
	"os"
	"os/signal"
)

// 启动 Proxy：Websocket 客户端、HTTP 代理
func StartProxy() {
	log.Println("监听进程信号")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go logics.ConsumeReplyMessageChan()

	// 连接到 Frontend
	app.ConnectFrontend()

	// 发送心跳
	go logics.SendHeartbeat()

	// 读取请求
	go logics.ReadRequest()

	// 当进程停止时，关闭连接
	<- interrupt
	log.Println("进程停止")
	app.DisconnectFrontend()
	log.Println("END")
}