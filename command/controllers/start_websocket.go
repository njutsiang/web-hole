package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/logics"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.URL.Query().Get("SecretKey") == app.Config.Frontend.SecretKey
	},
}

// 启动 Websocket 服务
func StartWebsocket() {
	log.Println("启动 Websocket 服务")
	go logics.ConsumeWebsocketConnChan()
	go logics.ConsumeWebsocketConnDelChan()
	http.HandleFunc(app.Config.Frontend.WebsocketPath, func(writer http.ResponseWriter, request *http.Request) {
		log.Println("建立 Websocket 连接")
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println("建立 Websocket 连接失败", err)
			return
		}
		connId := uuid.NewString()
		logics.AddWebsocketConn(&app.WebsocketConn{
			Id: connId,
			Conn: conn,
		})
		defer func() {
			log.Println("关闭 Websocket 连接")
			logics.DelWebsocketConn(connId)
			err = conn.Close()
			if err != nil {
				log.Println("关闭 Websocket 连接时报错", err)
			}
		}()
		for {
			messageType, messageBody, messageErr := conn.ReadMessage()
			if messageErr != nil {
				log.Println("读取消息失败", messageErr)
				break
			}
			log.Println("收到消息", string(messageBody))
			if messageType == websocket.TextMessage {
				go logics.ReceiveProxyResponse(messageBody)
			}
		}
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Frontend.WebsocketPort), nil)
	if err != nil {
		log.Println("启动 Websocket 服务失败", err)
	}
}