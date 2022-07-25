package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/njutsiang/web-hole/app"
	"github.com/njutsiang/web-hole/logics"
)

// 启动 HTTP 服务
func StartHttp() {
	app.Log.Info("启动 HTTP 服务")
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.MaxMultipartMemory = 8 << 20
	engine.Use(gin.Recovery())
	engine.NoRoute(logics.HttpHandler)
	engine.NoMethod(logics.HttpHandler)
	engine.Any("/", logics.HttpHandler)
	var err error
	if app.Config.Frontend.HttpsCertFile != "" && app.Config.Frontend.HttpsKeyFile != "" {
		err = engine.RunTLS(fmt.Sprintf(":%d", app.Config.Frontend.HttpPort), app.Config.Frontend.HttpsCertFile, app.Config.Frontend.HttpsKeyFile)
	} else {
		err = engine.Run(fmt.Sprintf(":%d", app.Config.Frontend.HttpPort))
	}
	if err != nil {
		app.Log.Error("启动 HTTP 服务失败 " + err.Error())
	}
}
