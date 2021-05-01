package main

import (
	"testing"
	"textlooker-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SourceTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Token            string
}

func (suite *SourceTestSuite) SetupTest() {
	email, password := "test@test.com", "Abcd124!"
	suite.UserRegistration, _ = models.NewUserRegistration(email)
	suite.User, _ = models.NewUser(email, password, *suite.UserRegistration)

	data := map[string]interface{}{
		"password": password,
		"email":    email,
	}

	response, _ := Post("/login", data, suite.Token)
	suite.Token = response["token"].(string)
}

func (suite *SourceTestSuite) CleanupTest() {
	CleanupDatabase()
}

func TestSourceTestSuite(t *testing.T) {
	suite.Run(t, new(SourceTestSuite))
}

func (suite *SourceTestSuite) TestPostSource() {
	data := map[string]interface{}{
		"name": "My source name",
	}
	response, code := Post("/auth/sources", data, suite.Token)

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), "Source created", response["status"])
	assert.NotNil(suite.T(), response["sourceID"])

	_, code = Post("/auth/sources", data, suite.Token)
	assert.NotEqual(suite.T(), 200, code)
}
