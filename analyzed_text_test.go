package main

import (
	"testing"
	"textlooker-backend/models"
	"textlooker-backend/util"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyzedTextTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Token            string
	Source           *models.Source
}

func (suite *AnalyzedTextTestSuite) SetupSuite() {
	email, password := "test4@test.com", "Abcd124!"
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

func (suite *AnalyzedTextTestSuite) CleanupSuite() {
	CleanupDatabase()
}

func TestAnalyzedTextTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzedTextTestSuite))
}

func (suite *AnalyzedTextTestSuite) TestGetAnalyzedText() {
	randomAuthor := util.RandStringRunes(10)
	text, _ := models.NewText("Abraham Lincoln is the first president of United Stated Of America.", randomAuthor, time.Now(), int(suite.Source.ID))
	models.NewAnalyzedText(text)
	analyzedTexts, _ := models.GetAnalyzedTexts("Abraham", randomAuthor, time.Now().Add(-3*time.Hour), time.Now(), int(suite.Source.ID))

	assert.Contains(suite.T(), analyzedTexts[0].Content, "first president")
	assert.Contains(suite.T(), analyzedTexts[0].People, "Abraham Lincoln")
	assert.Equal(suite.T(), analyzedTexts[0].Author, randomAuthor)
}
