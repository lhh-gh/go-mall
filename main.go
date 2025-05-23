package main

import (
	"github.com/gin-gonic/gin"
	"github/lhh-gh/go-mall/config"
)

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		database := config.Database
		c.JSON(200, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	r.Run(":8080")
}
