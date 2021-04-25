package main

import (
	"fmt"
	"net/http"
	"os"
	"textlooker-backend/controllers"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type RunMode uint8

const Production, Development, Test = 1, 2, 3

func SetupRouter(runMode RunMode) *gin.Engine {
	var router *gin.Engine
	switch runMode {
	case Development:
		router = gin.Default()
		fmt.Println("shashi")
	case Test:
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	case Production:
		router = gin.New()
	}

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/users", controllers.PostUser)

	return router
}

func main() {
	argument := os.Args[1]

	switch argument {
	case "migrate":
		models.ApplyMigrations("gorm")
	case "testdbsetup":
		models.ApplyMigrations("gorm_test")
	case "run":
		models.ConnectDatabase("gorm")
		r := SetupRouter(Development)
		r.Run(":8080")
	}

}
