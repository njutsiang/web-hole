package logics

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"log"
	"math/rand"
)

// 添加 Websocket 连接到全局
func AddWebsocketConn(conn *app.WebsocketConn) {
	app.WebsocketConnChan <- conn
}

// 添加 Websocket 连接到全局
func ConsumeWebsocketConnChan() {
	for websocketConn := range app.WebsocketConnChan {
		app.WebsocketConns = append(app.WebsocketConns, websocketConn)
		log.Println("添加 Websocket 连接到全局", websocketConn.Id)
	}
}

// 从全局移除 Websocket 连接
func DelWebsocketConn(connId string) {
	app.WebsocketConnDelChan <- connId
}

// 从全局移除 Websocket 连接
func ConsumeWebsocketConnDelChan() {
	for connId := range app.WebsocketConnDelChan {
		for i, websocketConn := range app.WebsocketConns {
			if websocketConn.Id == connId {
				if i == len(app.WebsocketConns) - 1 {
					app.WebsocketConns = app.WebsocketConns[0:i]
				} else {
					app.WebsocketConns = append(app.WebsocketConns[0:i], app.WebsocketConns[i+1:]...)
				}
				log.Println("从全局移除 Websocket 连接", connId)
				break
			}
		}
	}
}

// 接收代理的响应
func ReceiveProxyResponse(message []byte) {
	log.Println("接收代理的响应", string(message))
	response := app.Response{}
	err := json.Unmarshal(message, &response)
	if err != nil {
		log.Println("解析代理的响应失败", err)
		return
	}
	if response.RequestId == "" {
		log.Println("RequestId 为空")
		return
	}
	app.ResponseChanMapActionChan <- app.ResponseChanMapAction{
		Name: "write",
		Response: response,
	}
}

// 消费将请求发送至代理服务的通道
func ConsumeRequestMessageChan() {
	for message := range app.RequestMessageChan {
		websocketConn := app.WebsocketConns[rand.Intn(len(app.WebsocketConns))]
		err := websocketConn.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("将请求发送至代理服务失败")
		}
	}
}
