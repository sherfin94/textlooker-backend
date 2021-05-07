package main

import (
	"net/http"
	"os"
	"textlooker-backend/controllers"
	"textlooker-backend/database"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/middleware"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter(runMode deployment.RunMode) *gin.Engine {
	var router *gin.Engine

	deployment.CurrentRunMode = runMode
	deployment.InitiateEnv()

	switch runMode {
	case deployment.Development:
		router = gin.Default()
	case deployment.Test:
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	case deployment.Production:
		router = gin.New()
	}

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

	auth.POST("/sources", controllers.PostSource)
	auth.GET("/sources", controllers.GetSources)
	auth.DELETE("/sources/:sourceID", controllers.DeleteSource)

	auth.POST("/text", controllers.PostText)

	return router
}

func main() {
	argument := os.Args[1]

	switch argument {
	case "migrate":
		models.ApplyMigrations("gorm", database.Loud)
	case "migrate-test":
		models.ApplyMigrations("gorm_test", database.Loud)
	case "run":
		elastic.Initiate()
		database.ConnectDatabase("gorm", database.OnlyErrors)
		r := SetupRouter(deployment.Development)
		r.Run(":8080")
	}

}
