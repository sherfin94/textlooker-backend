package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"textlooker-backend/elastic"
	"textlooker-backend/models"
	"textlooker-backend/util"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TextTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Cookies          []*http.Cookie
	Source           *models.Source
}

func (suite *TextTestSuite) SetupSuite() {
	email, password := "test2@test.com", "Abcd124!"
	suite.UserRegistration, _ = models.NewUserRegistration(email, password)
	suite.User, _ = models.NewUser(email, *suite.UserRegistration)
	suite.Source, _ = models.NewSource("My Source", suite.User, true, true)

	data := map[string]interface{}{
		"password": password,
		"email":    email,
	}

	_, _, cookies := Post("/login", data, nil)
	suite.Cookies = cookies
}

func (suite *TextTestSuite) CleanupSuite() {
	CleanupDatabase()
}

func TestTextTestSuite(t *testing.T) {
	suite.Run(t, new(TextTestSuite))
}

func (suite *TextTestSuite) TestPostText() {
	text := map[string]interface{}{
		"content": "Abraham Lincoln is an amazing President. The United States of America is a good country.",
		"author":  []string{"Some person", "some other person"},
		"date":    strconv.FormatInt(time.Now().Unix(), 10),
	}

	data := map[string]interface{}{
		"batch":    []interface{}{text},
		"sourceID": suite.Source.ID,
	}

	response, code, _ := Post("/auth/text", data, suite.Cookies)

	assert.Equal(suite.T(), 200, code)
	assert.Equal(suite.T(), float64(1), response["savedTextCount"])

	text["sourceID"] = 0
	data = map[string]interface{}{
		"batch": []interface{}{text},
	}
	response, code, _ = Post("/auth/text", data, suite.Cookies)
	assert.Equal(suite.T(), 400, code)
}

func (suite *TextTestSuite) TestGetTextsFunc() {
	randomText := util.RandStringRunes(20)
	randomAuthor := util.RandStringRunes(10)
	models.BulkSaveText([]models.Text{
		{
			Content:  randomText,
			Author:   []string{randomAuthor},
			Date:     time.Now(),
			SourceID: int(suite.Source.ID),
		},
	})
	texts, _ := models.GetTexts(randomText, []elastic.FilterItem{}, time.Now().Add(-3*time.Hour), time.Now(), int(suite.Source.ID))

	assert.Contains(suite.T(), texts[0].Content, randomText)
	assert.Equal(suite.T(), texts[0].Author, []string{randomAuthor})
}

func (suite *TextTestSuite) TestGetTexts() {
	randomText := util.RandStringRunes(20)
	randomAuthor := util.RandStringRunes(10)
	models.BulkSaveText([]models.Text{
		{
			Content:  randomText,
			Author:   []string{randomAuthor},
			Date:     time.Now(),
			SourceID: int(suite.Source.ID),
		},
	})

	data := map[string]string{
		"content":   randomText,
		"author[]":  randomAuthor,
		"startDate": time.Now().Add(-3 * time.Hour).Format("Jan 2 15:04:05 -0700 MST 2006"),
		"endDate":   time.Now().Add(5 * time.Second).Format("Jan 2 15:04:05 -0700 MST 2006"),
		"sourceID":  fmt.Sprint(suite.Source.ID),
	}

	response, code := Get("/auth/text", data, suite.Cookies)

	assert.Equal(suite.T(), 200, code)
	assert.Contains(suite.T(), (response["texts"].([]interface{})[0].(map[string]interface{})["content"]), randomText)
	assert.Contains(suite.T(), (response["texts"].([]interface{})[0].(map[string]interface{})["author"]), randomAuthor)
}
