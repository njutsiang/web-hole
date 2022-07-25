package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/logics"
	"log"
)

// 启动 HTTP 服务
func StartHttp() {
	log.Println("启动 HTTP 服务")
	go logics.ConsumeRequestMessageChan()
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	engine.MaxMultipartMemory = 8 << 20
	engine.NoRoute(logics.HttpHandler)
	engine.NoMethod(logics.HttpHandler)
	engine.Any("/", logics.HttpHandler)
	err := engine.Run(fmt.Sprintf(":%d", app.Config.Frontend.HttpPort))
	if err != nil {
		log.Println("启动 HTTP 服务失败", err)
	}
}
