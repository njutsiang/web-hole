package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/njutsiang/web-hole/app"
	"net/http"
)

// 启动 Backend 服务
func StartBackend() {
	app.Log.Info("启动 Backend 服务")
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	engine.MaxMultipartMemory = 8 << 20
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": "demo",
		})
	})
	engine.GET("/demo1", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": "demo1",
		})
	})
	engine.GET("/demo2", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": "demo2",
		})
	})
	err := engine.Run(fmt.Sprintf(":%d", app.Config.Backend.HttpPort))
	if err != nil {
		app.Log.Error("启动 Backend 服务失败 " + err.Error())
	}
}