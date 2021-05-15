package main

import (
	"fmt"
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

func (suite *AnalyzedTextTestSuite) TestGetAnalyzedTextFunc() {
	randomAuthor := util.RandStringRunes(10)
	text, _ := models.NewText("Abraham Lincoln is the first president of United Stated Of America.", randomAuthor, time.Now(), int(suite.Source.ID))
	models.NewAnalyzedText(text)
	analyzedTexts,_, _ := models.GetAnalyzedTexts("Abraham", randomAuthor, []string{"Abraham Lincoln"}, []string{"America"}, time.Now().Add(-3*time.Hour), time.Now().Add(time.Minute), int(suite.Source.ID))

	assert.Contains(suite.T(), analyzedTexts[0].Content, "first president")
	assert.Contains(suite.T(), analyzedTexts[0].People, "Abraham Lincoln")
	assert.Equal(suite.T(), analyzedTexts[0].Author, randomAuthor)
}

func (suite *TextTestSuite) TestGetAnalyzedTexts() {
	randomText := "Abraham Lincoln is a good president of the United States Of America."
	randomAuthor := util.RandStringRunes(10)
	text, _ := models.NewText(randomText, randomAuthor, time.Now(), int(suite.Source.ID))
	models.NewAnalyzedText(text)

	data := map[string]string{
		"startDate": time.Now().Add(-3 * time.Hour).Format("Jan 2 15:04:05 -0700 MST 2006"),
		"endDate":   time.Now().Add(5 * time.Second).Format("Jan 2 15:04:05 -0700 MST 2006"),
		"sourceID":  fmt.Sprint(suite.Source.ID),
		"people":    "Abraham Lincoln",
	}

	response, code := Get("/auth/analyzed_text", data, suite.Token)

	assert.Equal(suite.T(), 200, code)
	assert.Contains(suite.T(), (response["texts"].([]interface{})[0].(map[string]interface{})["content"]), "good president")
	assert.Contains(suite.T(), (response["texts"].([]interface{})[0].(map[string]interface{})["author"]), randomAuthor)
}
