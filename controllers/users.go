package controllers

import (
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostUser(context *gin.Context) {
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := models.NewUser(user.Email, user.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "User created"})
}
