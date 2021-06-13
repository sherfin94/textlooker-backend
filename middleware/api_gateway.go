package middleware

import (
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type apiGatewayParams struct {
	ApiToken string `json:"apiToken"`
}

func APIGateway() gin.HandlerFunc {
	return func(context *gin.Context) {
		var params apiGatewayParams

		if err := context.ShouldBindJSON(&params); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		source, err := models.GetSourceByToken(params.ApiToken)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Source not found"})
			return
		}

		context.Set("source-id", source.ID)

		context.Next()
	}
}
