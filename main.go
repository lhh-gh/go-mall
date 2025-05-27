package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/app"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/middleware"
	"github/lhh-gh/go-mall/config"
	"net/http"
)

func main() {
	g := gin.New()
	// 有了AccessLog 后, 就不需要gin.Logger这个中间件啦
	// g.Use(gin.Logger(), middleware.StartTrace())
	g.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
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
	g.POST("/access-log-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	g.GET("/customized-error-test", func(c *gin.Context) {

		// 使用 Wrap 包装原因error 生成 项目error
		err := errors.New("a dao error")
		appErr := errcode.Wrap("包装错误", err)
		bAppErr := errcode.Wrap("再包装错误", appErr)
		logger.New(c).Error("记录错误", "err", bAppErr)

		// 预定义的ErrServer, 给其追加错误原因的error
		err = errors.New("a domain error")
		apiErr := errcode.ErrServer.WithCause(err)
		logger.New(c).Error("API执行中出现错误", "err", apiErr)

		c.JSON(apiErr.HttpStatusCode(), gin.H{
			"code": apiErr.Code(),
			"msg":  apiErr.Msg(),
		})

	})
	g.GET("/response-list", func(c *gin.Context) {

		pagination := app.NewPagination(c)
		// Mock fetch list data from db
		data := []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			{
				Name: "Lily",
				Age:  26,
			},
			{
				Name: "Violet",
				Age:  25,
			},
		}
		pagination.SetTotalRows(2)
		app.NewResponse(c).SetPagination(pagination).Success(data)
		return
	})
	g.GET("/response-error", func(c *gin.Context) {

		baseErr := errors.New("a dao error")
		// 这一步正式开发时写在service层
		err := errcode.Wrap("encountered an error when xxx service did xxx", baseErr)
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	})
	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
