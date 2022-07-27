package logics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/exception"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 连接到 Frontend
func ConnectFrontend(num int, options ...bool) *websocket.Conn {
	proxyWebsocket, _, err := websocket.DefaultDialer.Dial(app.Config.Proxy.FrontendUrl + "?" + (url.Values{"SecretKey":{app.Config.Proxy.SecretKey}}).Encode(), nil)
	if err != nil {
		app.Log.Error(fmt.Sprintf("连接到 Frontend-%d 失败 %s", num, err.Error()))
		if len(options) >= 1 && options[0] {
			time.Sleep(3 * time.Second)
			return ConnectFrontend(num, options[0])
		} else {
			exception.Throw(exception.ConnectFrontendFailed)
			return nil
		}
	}
	app.Log.Info(fmt.Sprintf("连接到 Frontend-%d 成功", num))
	return proxyWebsocket
}

// 断开和 Frontend 的连接
func DisconnectFrontend() {
	for num, proxyWebsocket := range app.ProxyWebsockets {
		if proxyWebsocket != nil {
			err := proxyWebsocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				app.Log.Error(fmt.Sprintf("向 Frontend-%d 发送 CloseMessage 失败 %s", num, err.Error()))
			}
			err = proxyWebsocket.Close()
			if err == nil {
				app.Log.Error(fmt.Sprintf("已关闭与 Frontend-%d 的连接", num))
			} else {
				app.Log.Error(fmt.Sprintf("关闭与 Frontend-%d 的连接失败 %s", num, err.Error()))
			}
		}
	}
	app.ProxyWebsockets = []*websocket.Conn{}
}

// 发送心跳
func SendHeartbeat(num int) {
	for {
		time.Sleep(3 * time.Second)
		if !(len(app.ProxyWebsockets) >= num + 1 && app.ProxyWebsockets[num] != nil) {
			app.ProxyWebsockets[num] = ConnectFrontend(num, true)
		}
		app.ReplyMessageQueues[num] <- app.ReplyMessage{
			Type: websocket.PingMessage,
		}
	}
}

// 读取请求
func ReadRequest(num int) {
	for {
		if !(len(app.ProxyWebsockets) >= num + 1 && app.ProxyWebsockets[num] != nil) {
			app.Log.Error(fmt.Sprintf("与 Frontend-%d 的连接不存在", num))
			time.Sleep(3 * time.Second)
			continue
		}
		messageType, messageBody, messageErr := app.ProxyWebsockets[num].ReadMessage()
		if messageErr != nil {
			app.Log.Error(fmt.Sprintf("与 Frontend-%d 的连接异常 %s", num, messageErr.Error()))
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
			app.Log.Info(fmt.Sprintf("接收到来自 Frontend-%d 请求：%s %s %s", num, request.Id, request.Method, request.Uri))
			ProxyRequest(num, request)
		}
	}
}

// 代理请求
func ProxyRequest(num int, request app.Request) {
	var body io.Reader
	if len(request.Body) >= 1 {
		body = bytes.NewReader(request.Body)
	}
	newRequest, err := http.NewRequest(request.Method, app.Config.Proxy.BackendHost + request.Uri, body)
	if err != nil {
		ReplyError(num, request.Id, err)
		return
	}
	for key, values := range request.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}
	client := &http.Client{
		Timeout: time.Duration(request.Timeout) * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := client.Do(newRequest)
	if err != nil {
		ReplyError(num, request.Id, err)
		return
	}
	client.CloseIdleConnections()
	ReplyResponse(num, request.Id, response)
}

// 回复一个错误
func ReplyError(num int, requestId string, err error) {
	response := app.Response{
		RequestId: requestId,
		StatusCode: http.StatusBadGateway,
		Body: []byte(err.Error()),
	}
	app.Log.Info(fmt.Sprintf("向 Frontend-%d 回复响应：%s %d %s", num, requestId, response.StatusCode, err.Error()))
	responseJson, _ := json.Marshal(response)
	app.ReplyMessageQueues[num] <- app.ReplyMessage{
		Type: websocket.TextMessage,
		Data: responseJson,
	}
}

// 回复一个响应
func ReplyResponse(num int, requestId string, response *http.Response) {
	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		ReplyError(num, requestId, err)
		return
	}
	newResponse := app.Response{
		RequestId: requestId,
		StatusCode: response.StatusCode,
		Header: response.Header,
		Body: body,
	}
	app.Log.Info(fmt.Sprintf("向 Frontend-%d 回复响应：%s %d", num, requestId, response.StatusCode))
	newResponseJson, _ := json.Marshal(newResponse)
	app.ReplyMessageQueues[num] <- app.ReplyMessage{
		Type: websocket.TextMessage,
		Data: newResponseJson,
	}
}

// 消费消息队列：待回复的消息
func ConsumeReplyMessageQueue(num int) {
	for replyMessage := range app.ReplyMessageQueues[num] {
		if !(len(app.ProxyWebsockets) >= num + 1 && app.ProxyWebsockets[num] != nil) {
			app.Log.Error(fmt.Sprintf("与 Frontend-%d 的连接不存在", num))
			continue
		}
		err := app.ProxyWebsockets[num].WriteMessage(replyMessage.Type, replyMessage.Data)
		if err != nil {
			app.Log.Error(fmt.Sprintf("向 Frontend-%d 发送响应失败 %s", num, err.Error()))
			app.ProxyWebsockets[num] = nil
		}
	}
}