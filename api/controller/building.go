package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/app"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/config"
	"github/lhh-gh/go-mall/logic/appservice"
	"net/http"
)

func TestConfigRead(c *gin.Context) {
	database := config.Database
	c.JSON(http.StatusOK, gin.H{
		"type":     database.Type,
		"max_life": database.Master.MaxLifeTime,
	})
	return
}

func TestLogger(c *gin.Context) {
	logger.New(c).Info("logger test", "key", "keyName", "val", 2)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

func TestAccessLog(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
	return
}

func TestPanicLog(c *gin.Context) {
	var a map[string]string
	a["k"] = "v"
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   a,
	})
	return
}

func TestAppError(c *gin.Context) {

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
	return
}

func TestResponseObj(c *gin.Context) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	app.NewResponse(c).Success(data)
	return
}

func TestResponseList(c *gin.Context) {

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
}

func TestResponseError(c *gin.Context) {
	baseErr := errors.New("a dao error")
	// 这一步正式开发时写在service层
	err := errcode.Wrap("encountered an error when xxx service did xxx", baseErr)
	app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	return
}

func TestPing(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	return
}
func TestGormLogger(c *gin.Context) {
	svc := appservice.NewDemoAppSvc(c)
	list, err := svc.GetDemoIdentities()
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(c).Success(list)
	return
}
