package main

import (
	"fmt"
	"net/http"
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
	Cookies          []*http.Cookie
}

func (suite *SourceTestSuite) SetupSuite() {

	email, password := "test@test.com", "Abcd124!"
	suite.UserRegistration, _ = models.NewUserRegistration(email, password)
	suite.User, _ = models.NewUser(email, *suite.UserRegistration)

	data := map[string]interface{}{
		"password": password,
		"email":    email,
	}

	_, _, cookies := Post("/login", data, nil)
	suite.Cookies = cookies
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
	response, code, _ := Post("/auth/sources", data, suite.Cookies)

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), "Source created", response["status"])
	assert.NotNil(suite.T(), response["sourceID"])

	_, code, _ = Post("/auth/sources", data, suite.Cookies)
	assert.NotEqual(suite.T(), 200, code)
}

func (suite *SourceTestSuite) TestGetSource() {

	source, _ := models.NewSource("My new source", suite.User)

	response, code := Get("/auth/sources", nil, suite.Cookies)

	respondedSources := response["sources"].([]interface{})

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), float64(source.ID), respondedSources[0].(map[string]interface{})["id"])
}

func (suite *SourceTestSuite) TestDeleteSource() {
	source, _ := models.NewSource("Another source", suite.User)
	_, code := Delete("/auth/sources/"+fmt.Sprint(source.ID), suite.Cookies)

	assert.Equal(suite.T(), 200, code)
}
