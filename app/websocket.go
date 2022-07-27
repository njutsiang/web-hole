package app

import "github.com/gorilla/websocket"

// Frontend 和 Proxy 的连接
var FrontendWebsockets = []*FrontendWebsocket{}

// 消息队列：添加 Frontend 和 Proxy 的连接
var FrontendWebsocketChan = make(chan *FrontendWebsocket, 1)

// 消息队列：删除 Frontend 和 Proxy 的连接
var FrontendWebsocketDelChan = make(chan string, 1)

// Frontend 和 Proxy 的连接
type FrontendWebsocket struct {
	Id string
	Conn *websocket.Conn
}
