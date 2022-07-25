package logics

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"hole/app"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 发送心跳
func SendHeartbeat() {
	var err error
	for {
		time.Sleep(10 * time.Second)
		if app.ProxyConn == nil {
			app.ConnectFrontend(true)
		}
		err = app.ProxyConn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			log.Println("发送心跳失败", err)
			app.ConnectFrontend(true)
		}
	}
}

// 读取请求
func ReadRequest() {
	for {
		if app.ProxyConn == nil {
			log.Println("Websocket 连接不存在")
			time.Sleep(3 * time.Second)
			continue
		}
		messageType, messageBody, messageErr := app.ProxyConn.ReadMessage()
		if messageErr != nil {
			log.Println("读取消息失败", messageErr)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("收到消息", string(messageBody))
		if messageType == websocket.TextMessage {
			request := app.Request{}
			err := json.Unmarshal(messageBody, &request)
			if err != nil {
				log.Println("解析消息失败", err)
				continue
			}
			ProxyRequest(request)
		}
	}
}

// 代理请求
func ProxyRequest(request app.Request) {
	var body io.Reader
	if len(request.Body) >= 1 {
		body = bytes.NewReader(request.Body)
	}
	newRequest, err := http.NewRequest(request.Method, app.Config.Proxy.BackendHost + request.Uri, body)
	if err != nil {
		ReplyError(request.Id, err)
		return
	}
	for key, values := range request.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}
	response, err := (&http.Client{
		Timeout: time.Duration(request.Timeout) * time.Second,
	}).Do(newRequest)
	if err != nil {
		ReplyError(request.Id, err)
		return
	}
	ReplyResponse(request.Id, response)
}

// 回复一个错误
func ReplyError(requestId string, err error) {
	response := app.Response{
		RequestId: requestId,
		StatusCode: http.StatusBadGateway,
		Body: []byte(err.Error()),
	}
	responseJson, _ := json.Marshal(response)
	app.ReplyMessageChan <- responseJson
}

// 回复一个响应
func ReplyResponse(requestId string, response *http.Response) {
	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		ReplyError(requestId, err)
		return
	}
	newResponse := app.Response{
		RequestId: requestId,
		StatusCode: response.StatusCode,
		Header: response.Header,
		Body: body,
	}
	newResponseJson, _ := json.Marshal(newResponse)
	app.ReplyMessageChan <- newResponseJson
}

// 消费回复消息的队列
func ConsumeReplyMessageChan() {
	for message := range app.ReplyMessageChan {
		if app.ProxyConn == nil {
			log.Println("Websocket 连接不存在")
			continue
		}
		err := app.ProxyConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("回复消息失败", err)
		}
	}
}