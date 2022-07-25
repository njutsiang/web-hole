package app

import (
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/exception"
	"net/url"
	"time"
)

// Proxy 和 Frontend 的连接
var ProxyConn *websocket.Conn

// 消息队列：待回复的消息
var ReplyMessageChan = make(chan []byte, 1000)

// 连接到 Frontend
func ConnectFrontend(options ...bool) {
	var err error
	Log.Info("连接到 Frontend")
	ProxyConn, _, err = websocket.DefaultDialer.Dial(Config.Proxy.FrontendUrl + "?" + (url.Values{"SecretKey":{Config.Proxy.SecretKey}}).Encode(), nil)
	if err != nil {
		Log.Error("连接到 Frontend 失败 " + err.Error())
		if len(options) >= 1 && options[0] {
			time.Sleep(3 * time.Second)
			ConnectFrontend(options[0])
		} else {
			exception.Throw(exception.ConnectFrontendFailed)
		}
		return
	}
	Log.Info("连接到 Frontend 成功")
}

// 断开和 Frontend 的连接
func DisconnectFrontend() {
	if ProxyConn != nil {
		err := ProxyConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			Log.Error("发送 CloseMessage 失败 " + err.Error())
		}
		err = ProxyConn.Close()
		if err != nil {
			Log.Error("关闭与 Frontend 的连接失败 " + err.Error())
		}
	}
}