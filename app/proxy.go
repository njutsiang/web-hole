package app

import (
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/exception"
	"log"
	"net/url"
	"time"
)

var ProxyConn *websocket.Conn

// 回复消息的队列
var ReplyMessageChan = make(chan []byte, 1000)

// 连接到 Frontend
func ConnectFrontend(options ...bool) {
	var err error
	log.Println("建立 Websocket 连接")
	ProxyConn, _, err = websocket.DefaultDialer.Dial(Config.Proxy.FrontendUrl + "?" + (url.Values{"SecretKey":{Config.Proxy.SecretKey}}).Encode(), nil)
	if err != nil {
		log.Println("建立 Websocket 连接失败", err)
		if len(options) >= 1 && options[0] {
			time.Sleep(3 * time.Second)
			ConnectFrontend(options[0])
		} else {
			exception.Throw(exception.ConnectFrontendFailed)
		}
		return
	}
	log.Println("建立 Websocket 连接成功")
}

// 断开和 Frontend 的连接
func DisconnectFrontend() {
	if ProxyConn != nil {
		err := ProxyConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("发送 CloseMessage 失败", err)
		}
		err = ProxyConn.Close()
		if err != nil {
			log.Println("关闭 Websocket 连接时报错", err)
		}
	}
}