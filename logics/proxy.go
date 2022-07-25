package logics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"io"
	"io/ioutil"
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
			app.Log.Error("发送心跳失败 " + err.Error())
			app.ConnectFrontend(true)
		}
	}
}

// 读取请求
func ReadRequest() {
	for {
		if app.ProxyConn == nil {
			app.Log.Error("与 Frontend 的连接不存在")
			time.Sleep(3 * time.Second)
			continue
		}
		messageType, messageBody, messageErr := app.ProxyConn.ReadMessage()
		if messageErr != nil {
			app.Log.Error("与 Frontend 的连接异常 " + messageErr.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		if messageType == websocket.TextMessage {
			request := app.Request{}
			err := json.Unmarshal(messageBody, &request)
			if err != nil {
				app.Log.Error("解析请求失败 " + err.Error())
				continue
			}
			app.Log.Info("接收到请求：" + request.Id + " " + request.Method + " " + request.Uri)
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
	app.Log.Info(fmt.Sprintf("向 Frontend 回复响应：%s %d %s", requestId, response.StatusCode, err.Error()))
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
	app.Log.Info(fmt.Sprintf("向 Frontend 回复响应：%s %d", requestId, response.StatusCode))
	newResponseJson, _ := json.Marshal(newResponse)
	app.ReplyMessageChan <- newResponseJson
}

// 消费消息队列：待回复的消息
func ConsumeReplyMessageChan() {
	for message := range app.ReplyMessageChan {
		if app.ProxyConn == nil {
			app.Log.Error("与 Frontend 的连接不存在")
			continue
		}
		err := app.ProxyConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			app.Log.Error("向 Frontend 发送响应失败" + err.Error())
		}
	}
}