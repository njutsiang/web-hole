package logics

import (
	"encoding/json"
	"github.com/njutsiang/web-hole/app"
)

// 消费消息队列：对 ResponseChanMap 的操作
func ConsumeResponseChanMapActionChan() {
	for action := range app.ResponseChanMapActionChan {
		switch action.Name {
		case "add":
			app.ResponseChanMap[action.Request.Id] = action.ResponseChan
			requestJson, _ := json.Marshal(action.Request)
			app.RequestMessageChan <- requestJson
		case "del":
			close(action.ResponseChan)
			delete(app.ResponseChanMap, action.RequestId)
		case "write":
			responseChan, ok := app.ResponseChanMap[action.Response.RequestId]
			if !ok {
				app.Log.Error("响应的通道不存在")
				continue
			}
			responseChan <- action.Response
		}
	}
}
