package controllers

import "github.com/njutsiang/web-hole/logics"

// 启动 Frontend
func StartFrontend() {

	go logics.ConsumeResponseChanMapActionChan()

	// 启动 HTTP 服务
	go StartHttp()

	// 启动 Websocket 服务
	StartWebsocket()
}
