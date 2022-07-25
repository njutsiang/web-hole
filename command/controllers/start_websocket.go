package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/logics"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.URL.Query().Get("SecretKey") == app.Config.Frontend.SecretKey
	},
}

// 启动 Websocket 服务
func StartWebsocket() {
	app.Log.Info("启动 Websocket 服务")
	http.HandleFunc(app.Config.Frontend.WebsocketPath, func(writer http.ResponseWriter, request *http.Request) {
		app.Log.Info("与 Proxy 建立连接")
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			app.Log.Error("与 Proxy 建立连接失败 " + err.Error())
			return
		}
		connId := uuid.NewString()
		logics.AddWebsocketConn(&app.WebsocketConn{
			Id: connId,
			Conn: conn,
		})
		defer func() {
			app.Log.Info("关闭与 Proxy 的连接")
			logics.DelWebsocketConn(connId)
			err = conn.Close()
			if err != nil {
				app.Log.Error("关闭与 Proxy 的连接失败 " + err.Error())
			}
		}()
		for {
			messageType, messageBody, messageErr := conn.ReadMessage()
			if messageErr != nil {
				app.Log.Error("与 Proxy 的连接异常 " + messageErr.Error())
				break
			}
			if messageType == websocket.TextMessage {
				go logics.ReceiveProxyResponse(messageBody)
			}
		}
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Frontend.WebsocketPort), nil)
	if err != nil {
		app.Log.Error("启动 Websocket 服务失败 " + err.Error())
	}
}