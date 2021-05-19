package controllers

import (
	"net/http"
	"textlooker-backend/deployment"
	"textlooker-backend/mailer"
	"textlooker-backend/models"

	"github.com/gin-gonic/gin"
)

type UserRegistration struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

func PostUserRegistration(context *gin.Context) {
	var userRegistration UserRegistration
	if err := context.ShouldBindJSON(&userRegistration); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if createdUserRegistration, err := models.NewUserRegistration(userRegistration.Email, userRegistration.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if !deployment.IsTest() {
			go mailer.SendMail(
				"TextLooker",
				"hi@textlooker.com",
				userRegistration.Email,
				userRegistration.Email,
				"Verification token for Textlooker",
				"Your verification token is "+createdUserRegistration.VerificationToken,
			)
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": "User registration created"})
}
