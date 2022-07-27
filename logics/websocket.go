package logics

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"math/rand"
)

// 添加 Websocket 连接到全局
func AddFrontendWebsocket(conn *app.FrontendWebsocket) {
	app.FrontendWebsocketChan <- conn
}

// 消费消息队列：添加 Frontend 和 Proxy 的连接
func ConsumeFrontendWebsocketChan() {
	for frontendWebsocket := range app.FrontendWebsocketChan {
		app.FrontendWebsockets = append(app.FrontendWebsockets, frontendWebsocket)
		app.Log.Info("将与 Proxy 的连接添加到全局 " + frontendWebsocket.Id)
	}
}

// 从全局移除 Websocket 连接
func DelFrontendWebsocket(connId string) {
	app.FrontendWebsocketDelChan <- connId
}

// 消费消息队列：删除 Frontend 和 Proxy 的连接
func ConsumeFrontendWebsocketDelChan() {
	for connId := range app.FrontendWebsocketDelChan {
		for i, frontendWebsocket := range app.FrontendWebsockets {
			if frontendWebsocket.Id == connId {
				if i == len(app.FrontendWebsockets) - 1 {
					app.FrontendWebsockets = app.FrontendWebsockets[0:i]
				} else {
					app.FrontendWebsockets = append(app.FrontendWebsockets[0:i], app.FrontendWebsockets[i+1:]...)
				}
				app.Log.Info("从全局移除与 Proxy 的连接 " + connId)
				break
			}
		}
	}
}

// 接收代理的响应
func ReceiveProxyResponse(message []byte) {
	response := app.Response{}
	err := json.Unmarshal(message, &response)
	if err != nil {
		app.Log.Error("解析响应失败 " + err.Error())
		return
	}
	if response.RequestId == "" {
		app.Log.Error("响应的 RequestId 为空")
		return
	}
	app.Log.Info(fmt.Sprintf("接收到 Proxy 的响应：%s %d", response.RequestId, response.StatusCode))
	app.ResponseChanMapActionChan <- app.ResponseChanMapAction{
		Name: "write",
		Response: response,
	}
}

// 消费消息队列：发送到代理服务器的消息
func ConsumeRequestMessageChan() {
	for message := range app.RequestMessageChan {
		frontendWebsocket := app.FrontendWebsockets[rand.Intn(len(app.FrontendWebsockets))]
		err := frontendWebsocket.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			app.Log.Error("将请求发送至 Proxy 失败 " + err.Error())
		}
	}
}
