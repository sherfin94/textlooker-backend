package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"textlooker-backend/deployment"
	"textlooker-backend/models"

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

func Get(url string, token string) (map[string]interface{}, int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response, w.Code
}

func TestMain(m *testing.M) {
	models.ConnectDatabase("gorm_test")
	CleanupDatabase()
	router = SetupRouter(deployment.Test)
	m.Run()
}

func CleanupDatabase() {
	models.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Source{})
	models.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
	models.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.UserRegistration{})
}
