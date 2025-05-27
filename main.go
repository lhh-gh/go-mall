package main

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/middleware"
	"github/lhh-gh/go-mall/config"
	"net/http"
)

func main() {
	g := gin.New()
	// TODO: 后面会把应用日志统一收集到文件， 这里根据运行环境判断, 只在dev环境下才使用gin.Logger()输出信息到控制台
	//g.Use(gin.Logger(), gin.Recovery())
	//调用  处理请求追踪的中间件
	g.Use(gin.Logger(), middleware.StartTrace())
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		// 测试Zap 初始化的临时代码,
		logger.ZapLoggerTest(c)

		c.JSON(http.StatusOK, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	//  测试log
	g.GET("/logger-test", func(c *gin.Context) {
		// 使用原始的 logger 方式
		logger.New(c).Info("logger test", "key", "keyName", "val", 2)

		// 使用新的 v1 门面方式
		logger.InfoV1(c, "logger test v1", "key", "keyName", "val", 2)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
