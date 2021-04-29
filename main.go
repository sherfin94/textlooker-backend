package main

import (
	"net/http"
	"os"
	"textlooker-backend/controllers"
	"textlooker-backend/deployment"
	"textlooker-backend/middleware"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter(runMode deployment.RunMode) *gin.Engine {
	var router *gin.Engine
	switch runMode {
	case deployment.Development:
		router = gin.Default()
	case deployment.Test:
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	case deployment.Production:
		router = gin.New()
	}

	router.Use(middleware.InitiateRunMode(runMode))

	authMiddleware := middleware.GenerateJWTAuthMiddleware()

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/users", controllers.PostUser)
	router.POST("/user_registrations", controllers.PostUserRegistration)
	router.POST("/login", authMiddleware.LoginHandler)

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())

	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hey!!")
	})

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
		r := SetupRouter(deployment.Development)
		r.Run(":8080")
	}

}
