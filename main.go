package main

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/config"
	"net/http"
)

func main() {
	g := gin.New()
	// TODO: 后面会把应用日志统一收集到文件， 这里根据运行环境判断, 只在dev环境下才使用gin.Logger()输出信息到控制台
	g.Use(gin.Logger(), gin.Recovery())
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		// 测试Zap 初始化的临时代码, 下节课会删掉
		logger.ZapLoggerTest(c)

		c.JSON(http.StatusOK, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
