package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type apiGatewayParams struct {
	ApiToken string `json:"token"`
}

func APIGateway() gin.HandlerFunc {
	return func(context *gin.Context) {
		var params apiGatewayParams

		requestBody := context.Request.Body
		bodyBytes, err := ioutil.ReadAll(requestBody)
		context.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err = json.Unmarshal(bodyBytes, &params)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		source, err := models.GetSourceByToken(params.ApiToken)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Source not found"})
			return
		} else {
			context.Set("source", source)
			context.Next()
		}

	}
}
