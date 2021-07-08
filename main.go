package main

import (
	"log"
	"net/http"
	"os"
	"textlooker-backend/api"
	"textlooker-backend/controllers"
	"textlooker-backend/database"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/kafka"
	"textlooker-backend/middleware"
	"textlooker-backend/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("auth"))
	router.Use(sessions.Sessions("authsession", store))

	authMiddleware := middleware.GenerateJWTAuthMiddleware()
	apiGatewayMiddleware := middleware.APIGateway()

	corsMiddleware := middleware.CORSMiddleware()
	router.Use(corsMiddleware)

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST("/users", controllers.PostUser)
	router.POST("/user_registrations", controllers.PostUserRegistration)
	router.POST("/login", authMiddleware.LoginHandler)
	router.GET("/logout", authMiddleware.LogoutHandler)

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())

	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	auth.POST("/sources", controllers.PostSource)
	auth.GET("/sources", controllers.GetSources)
	auth.DELETE("/sources/:sourceID", controllers.DeleteSource)

	auth.POST("/text", controllers.PostText)
	auth.GET("/text", controllers.GetTexts)

	auth.GET("/analyzed_text", controllers.GetAnalyzedTexts)

	auth.GET("/general_aggregation", controllers.GetGeneralAggregation)
	auth.GET("/per_date_aggregation", controllers.GetPerDateAggregation)
	auth.GET("/dateless_aggregation", controllers.GetDatelessGeneralAggregation)

	auth.POST("/insights", controllers.PostInsight)
	auth.GET("/insights", controllers.GetInsights)
	auth.DELETE("/insights", controllers.DeleteInsight)

	auth.POST("/dashboards", controllers.PostDashboard)
	auth.GET("/dashboards", controllers.GetDashboards)
	auth.DELETE("/dashboards", controllers.DeleteDashboard)

	auth.POST("/dashboard_insights", controllers.AddInsightToDashboard)

	router.GET("/dashboards", controllers.GetDashboardViaToken)
	router.GET("/dashboard_insights", controllers.GetInsightsAggregationViaToken)

	apiGateway := router.Group("/api")
	apiGateway.Use(apiGatewayMiddleware)

	apiGateway.POST("/text", api.PostText)

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
		database.ConnectDatabase(os.Getenv("DATABASE_NAME"), database.Silent)

		channel := make(chan kafka.TextSet)
		go kafka.InitializeProducer(&channel)

		var runMode deployment.RunMode

		switch os.Getenv("DEPLOYMENT_MODE") {
		case "DEVELOPMENT":
			{
				runMode = deployment.Development
			}
		case "PRODUCTION":
			{
				runMode = deployment.Development
			}
		default:
			{
				log.Fatalf("DEPLOYMENT_MODE is not set")
			}
		}

		r := SetupRouter(runMode)
		r.Run(":8080")
	}
}
