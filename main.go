package main

import (
	"net/http"
	"os"
	"textlooker-backend/controllers"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
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
		models.ApplyMigrations()
	case "run":
		r := setupRouter()
		r.Run(":8080")
	case "test":
		models.NewUser("sherfin@04.com", "hellohell")

	}

}
