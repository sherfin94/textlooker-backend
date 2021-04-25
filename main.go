package main

import (
	"net/http"
	"os"
	"textlooker-backend/controllers"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	// gin.SetMode(gin.ReleaseMode)
	// r := gin.New()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/users", controllers.PostUser)

	return r
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
		r := SetupRouter()
		r.Run(":8080")
	}

}
