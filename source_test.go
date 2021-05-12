package main

import (
	"fmt"
	"testing"
	"textlooker-backend/database"
	"textlooker-backend/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SourceTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Token            string
}

func (suite *SourceTestSuite) SetupSuite() {

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

func (suite *SourceTestSuite) CleanupSuite() {
	CleanupDatabase()
}

func (suite *SourceTestSuite) BeforeTest(suiteName, testName string) {
	database.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Source{})
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

func (suite *SourceTestSuite) TestGetSource() {

	source, _ := models.NewSource("My new source", suite.User)

	response, code := Get("/auth/sources", nil, suite.Token)

	respondedSources := response["sources"].([]interface{})

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), float64(source.ID), respondedSources[0].(map[string]interface{})["id"])
}

func (suite *SourceTestSuite) TestDeleteSource() {
	source, _ := models.NewSource("Another source", suite.User)

	_, code := Delete("/auth/sources/"+fmt.Sprint(source.ID), suite.Token)

	assert.Equal(suite.T(), 200, code)
}
