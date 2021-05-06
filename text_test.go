package main

import (
	"testing"
	"textlooker-backend/models"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TextTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Token            string
	Source           *models.Source
}

func (suite *TextTestSuite) SetupSuite() {
	email, password := "test2@test.com", "Abcd124!"
	suite.UserRegistration, _ = models.NewUserRegistration(email)
	suite.User, _ = models.NewUser(email, password, *suite.UserRegistration)
	suite.Source, _ = models.NewSource("My Source", suite.User)

	data := map[string]interface{}{
		"password": password,
		"email":    email,
	}

	response, _ := Post("/login", data, suite.Token)
	suite.Token = response["token"].(string)
}

func (suite *TextTestSuite) CleanupSuite() {
	CleanupDatabase()
}

func TestTextTestSuite(t *testing.T) {
	suite.Run(t, new(TextTestSuite))
}

func (suite *TextTestSuite) TestPostText() {
	data := map[string]interface{}{
		"content":  "My awesome new text",
		"author":   "Some person",
		"time":     time.Now().Format("Jan 2 15:04:05 -0700 MST 2006"),
		"sourceID": suite.Source.ID,
	}

	response, code := Post("/auth/text", data, suite.Token)

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), "Text saved", response["status"])
}
