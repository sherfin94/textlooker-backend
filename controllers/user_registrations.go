package controllers

import (
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type UserRegistration struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostUserRegistration(context *gin.Context) {
	var userRegistration UserRegistration
	if err := context.ShouldBindJSON(&userRegistration); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := models.NewUser(userRegistration.Email, userRegistration.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "User registration created"})
}
