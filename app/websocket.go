package app

import "github.com/gorilla/websocket"

// Frontend 和 Proxy 的连接
var WebsocketConns = []*WebsocketConn{}

// 消息队列：添加 Frontend 和 Proxy 的连接
var WebsocketConnChan = make(chan *WebsocketConn, 1)

// 消息队列：删除 Frontend 和 Proxy 的连接
var WebsocketConnDelChan = make(chan string, 1)

// Frontend 和 Proxy 的连接
type WebsocketConn struct {
	Id string
	Conn *websocket.Conn
}
