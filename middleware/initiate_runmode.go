package middleware

import (
	"textlooker-backend/deployment"

	"github.com/gin-gonic/gin"
)

func InitiateRunMode(runMode deployment.RunMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("run_mode", runMode)
		c.Next()
	}
}
