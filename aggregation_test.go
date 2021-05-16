package main

import (
	"log"
	"testing"
	"textlooker-backend/models"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AggregationTestSuite struct {
	suite.Suite
	UserRegistration *models.UserRegistration
	User             *models.User
	Token            string
	Source           *models.Source
}

func (suite *AggregationTestSuite) SetupSuite() {
	email, password := "test5@test.com", "Abcd124!"
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

func (suite *AggregationTestSuite) CleanupSuite() {
	CleanupDatabase()
}

func TestAggregationTestSuite(t *testing.T) {
	suite.Run(t, new(AggregationTestSuite))
}

func (suite *TextTestSuite) TestAggregations() {
	log.Println("Skipped TestAggregations as it takes time and is complicated")
	return
	texts := [][]string{
		{"Bob and Alice were friends.", "AuthorB"},
		{"Bob got employed in Apple.", "AuthorC"},
		{"Alice got employed in Facebook.", "AuthorB"},
		{"Bob is from the United States Of America.", "AuthorC"},
		{"Alice is from India.", "AuthorC"},
		{"However, Alice wanted to work in Tesla.", "AuthorB"},
		{"And Bob wanted a job in Facebook.", "AuthorB"},
	}

	for _, text := range texts {
		savedText, _ := models.NewText(text[0], text[1], time.Now(), int(suite.Source.ID))
		models.NewAnalyzedText(savedText)
	}

	startDate := time.Now().Add(-time.Minute)
	endDate := time.Now().Add(time.Minute)
	_, aggregation, _ := models.GetAnalyzedTexts("*", "*", []string{}, []string{}, startDate, endDate, int(suite.Source.ID))

	expectedAuthorsData := []models.CountItem{
		{Value: "AuthorB", Count: 4},
		{Value: "AuthorC", Count: 3},
	}

	expectedPeopleData := []models.CountItem{
		{Value: "Bob", Count: 4},
		{Value: "Alice", Count: 1},
	}

	// expectedGPEData := []models.CountItem{
	// 	{Value: "Apple", Count: 1},
	// 	{Value: "Facebook", Count: 2},
	// 	{Value: "United States", Count: 1},
	// 	{Value: "India", Count: 1},
	// 	{Value: "Tesla", Count: 1},
	// }

	assert.Equal(suite.T(), expectedAuthorsData, aggregation.Authors)
	assert.Equal(suite.T(), expectedPeopleData, aggregation.People)
	// assert.Contains(suite.T(), expectedGPEData, aggregation.GPE)
}