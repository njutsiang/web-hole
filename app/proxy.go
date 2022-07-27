package app

import (
	"github.com/gorilla/websocket"
)

// Proxy 和 Frontend 的连接
var ProxyWebsockets = []*websocket.Conn{}

// 消息队列：待回复的消息
var ReplyMessageQueues = []chan ReplyMessage{}

// 回复的消息
type ReplyMessage struct {
	Type int
	Data []byte
}