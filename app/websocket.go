package app

import "github.com/gorilla/websocket"

// 全局 Websocket 连接
var WebsocketConns = []*WebsocketConn{}

// 添加 Websocket 连接的队列
var WebsocketConnChan = make(chan *WebsocketConn, 1)

// 删除 Websocket 连接的队列
var WebsocketConnDelChan = make(chan string, 1)

// Websocket 连接
type WebsocketConn struct {
	Id string
	Conn *websocket.Conn
}
