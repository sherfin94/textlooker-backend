package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {

	data := map[string]interface{}{
		"password": "hellosjkfio",
		"email":    "Tfff@example.com",
	}

	response, code := Post("/users", data)

	assert.Equal(t, 200, code)
	assert.Equal(t, "User created", response["status"])
}
