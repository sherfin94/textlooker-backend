package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"textlooker-backend/database"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/models"
	"textlooker-backend/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

func Post(url string, data map[string]interface{}, token string) (map[string]interface{}, int) {
	marshalledData, _ := json.Marshal(data)
	postBody := bytes.NewBuffer(marshalledData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, postBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response, w.Code
}

func Get(url string, data map[string]string, token string) (map[string]interface{}, int) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", url, nil)
	if data != nil {
		query := req.URL.Query()
		for key, value := range data {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response, w.Code
}

func Delete(url string, token string) (map[string]interface{}, int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response, w.Code
}

func TestMain(m *testing.M) {
	util.RandomStringGeneratorInit()
	database.ConnectDatabase("gorm_test", database.Silent)
	elastic.Initiate()
	CleanupDatabase()
	router = SetupRouter(deployment.Test)
	m.Run()
}

func CleanupDatabase() {
	database.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Source{})
	database.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
	database.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.UserRegistration{})
}
