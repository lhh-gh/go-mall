package main

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/api/router"
	"github/lhh-gh/go-mall/comon/enum"
	"github/lhh-gh/go-mall/config"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()

	router.RegisterRoutes(g)

	g.Run(":8080")

}
