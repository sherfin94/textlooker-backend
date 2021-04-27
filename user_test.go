package main

import (
	"testing"
	"textlooker-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {

	userRegistration, _ := models.NewUserRegistration("Tfff@example.com", "hellosjkfio")

	data := map[string]interface{}{
		"password":          "hellosjkfio",
		"email":             "Tfff@example.com",
		"verificationToken": userRegistration.VerificationToken,
	}

	response, code := Post("/users", data)

	assert.Equal(t, 200, code)
	assert.Equal(t, "User created", response["status"])
}

func TestPostUserRegistration(t *testing.T) {
	data := map[string]interface{}{
		"password": "hellosjkfio",
		"email":    "Tfffex@ample.com",
	}

	response, code := Post("/user_registrations", data)

	assert.Equal(t, 200, code)
	assert.Equal(t, "User registration created", response["status"])
}
