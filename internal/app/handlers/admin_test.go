package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Trying to create user without an access token
	req, _ := http.NewRequest("POST", "/users", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for a user from db
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err := helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("POST", "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for an admin from db
	user, org, err = sampleFromDB(handler, "ADMIN")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err = helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	testUsername := "test"
	testPassword := "testing123"
	testUserType := "USER"

	testUser := models.UserCreate{
		Username: &testUsername,
		Password: &testPassword,
		UserType: &testUserType,
	}

	reqBody, err := json.Marshal(testUser)

	if err != nil {
		log.Panic("error marshalling test data to json")
	}

	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 201, resp.Result().StatusCode)
}

func TestUpdateUser(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Trying to update user without an access token
	req, _ := http.NewRequest("PATCH", "/users/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for a user from db
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err := helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("PATCH", "/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for an admin from db
	user, org, err = sampleFromDB(handler, "ADMIN")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err = helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	testPasswordUpdate := "testing123"

	testUser := models.UserCreate{
		Password: &testPasswordUpdate,
	}

	reqBody, err := json.Marshal(testUser)

	if err != nil {
		log.Panic("error marshalling test data to json")
	}

	req, _ = http.NewRequest("PATCH", "/users/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

func TestDeleteUser(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Trying to delete user without an access token
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for a user from db
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err := helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Generating access token for an admin from db
	user, org, err = sampleFromDB(handler, "ADMIN")

	if err != nil {
		log.Panic("error sampling from DB")
	}

	accessToken, _, err = helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	req, _ = http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}
