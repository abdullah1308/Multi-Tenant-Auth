package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Trying to access users without an access token
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for a user from db and accessing users
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err := helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

func TestGetUserByID(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Trying to access users without an access token
	req, _ := http.NewRequest("GET", "/users/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for a user from db and accessing users
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err := helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}
