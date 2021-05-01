package controllers

import (
	"fmt"
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type Source struct {
	Name string `json:"name" binding:"required"`
}

func PostSource(context *gin.Context) {
	var source models.Source

	if err := context.ShouldBindJSON(&source); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := context.Get("user")
	fmt.Println(user)

	if _, err := models.NewSource(source.Name, user.(*models.User)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Source created"})
}
