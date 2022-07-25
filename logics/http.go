package logics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/njutsiang/web-hole/app"
	"log"
	"net/http"
	"strings"
	"time"
)

// 处理 HTTP 请求
func HttpHandler(c *gin.Context) {

	// 准备请求
	body, _ := c.GetRawData()
	request := app.Request{
		Id: uuid.NewString(),
		Method: c.Request.Method,
		Uri: c.Request.RequestURI,
		Header: c.Request.Header,
		Body: body,
		Timeout: app.Config.Frontend.HttpTimeout,
	}

	// 接收响应的通道
	responseChan := make(chan app.Response)
	defer func() {
		app.ResponseChanMapActionChan <- app.ResponseChanMapAction{
			Name: "del",
			RequestId: request.Id,
			ResponseChan: responseChan,
		}
	}()

	// 判断代理服务器是否存在
	if len(app.WebsocketConns) == 0 {
		c.JSON(http.StatusBadGateway, map[string]interface{}{
			"error": "代理服务器不存在",
		})
		return
	}

	// 将发送请求的动作写入队列
	app.ResponseChanMapActionChan <- app.ResponseChanMapAction{
		Name: "add",
		Request: request,
		ResponseChan: responseChan,
	}

	// 等待代理服务器回复的响应，超时则返回错误
	timer := time.NewTimer(time.Duration(request.Timeout) * time.Second)
	select {
	case <- timer.C:
		c.JSON(http.StatusBadGateway, map[string]interface{}{
			"error": "请求超时",
		})
	case response := <- responseChan:
		timer.Stop()
		c.Status(response.StatusCode)
		for key, values := range response.Header {
			for _, value := range values {
				if strings.ToLower(key) != "content-length" {
					c.Writer.Header().Add(key, value)
				}
			}
		}
		_, err := fmt.Fprint(c.Writer, string(response.Body))
		if err != nil {
			log.Println("响应错误", err)
		}
	}
}
