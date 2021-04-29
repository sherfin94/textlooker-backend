package controllers

import (
	"net/http"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	VerificationToken string `json:"verificationToken" binding:"required"`
}

func PostUser(context *gin.Context) {
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userRegistration models.UserRegistration
	models.Database.Where("email = ?", user.Email).First(&userRegistration)

	if userRegistration.VerificationToken != user.VerificationToken {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Verifcation token is wrong"})
		return
	}

	if _, err := models.NewUser(user.Email, user.Password, userRegistration); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "User created"})
}
