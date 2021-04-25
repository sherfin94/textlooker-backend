package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"textlooker-backend/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	models.ConnectDatabase("gorm_test")
	models.Database.Unscoped().Session(
		&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
	m.Run()
}
func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	data, _ := json.Marshal(map[string]string{
		"password": "hellosjkfio",
		"email":    "Tfff@example.com",
	})

	postBody := bytes.NewBuffer(data)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", postBody)
	router.ServeHTTP(w, req)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, response["status"], "User created")
}
