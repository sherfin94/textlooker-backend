package main

import (
	"testing"
	"textlooker-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {

	userRegistration, _ := models.NewUserRegistration("Tfff@example.com", "Abcd1432!")

	data := map[string]interface{}{
		"email":             "Tfff@example.com",
		"verificationToken": userRegistration.VerificationToken,
	}

	response, code, _ := Post("/users", data, nil)

	assert.Equal(t, 200, code)
	assert.Equal(t, "User created", response["status"])
}

func TestPostUserRegistration(t *testing.T) {
	data := map[string]interface{}{
		"email":    "Tfffex@ample.com",
		"password": "some password",
	}

	response, code, _ := Post("/user_registrations", data, nil)

	assert.Equal(t, 200, code)
	assert.Equal(t, "User registration created", response["status"])
}

func TestPostLogin(t *testing.T) {
	email := "test@te2st.com"
	password := "myawesomepassword123"
	userRegistration, _ := models.NewUserRegistration(email, password)
	models.NewUser(email, *userRegistration)

	data := map[string]interface{}{
		"password": password,
		"email":    email,
	}

	_, code, cookies := Post("/login", data, nil)

	assert.Equal(t, 200, code)
	assert.Equal(t, true, len(cookies[0].Value) > 0)
}
