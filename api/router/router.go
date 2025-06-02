package router

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/comon/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	// use global middlewares
	engine.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engine.Group("")
	registerBuildingRoutes(routeGroup)
}
