package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrg(t *testing.T) {
	// organization name,username,password,expected status
	tests := []string{"org1,user1,password,403", "newOrg,user1,password,201"}

	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	for _, test := range tests {
		testInfo := strings.Split(test, ",")

		orgCreateData := models.OrganizationCreate{
			Name:     &testInfo[0],
			Username: &testInfo[1],
			Password: &testInfo[2],
		}

		expectedStatus, _ := strconv.Atoi(testInfo[3])

		reqBody, err := json.Marshal(orgCreateData)

		if err != nil {
			log.Panic("error marshalling test data to json")
		}

		req, _ := http.NewRequest("POST", "/orgs", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		t.Log(resp.Result().StatusCode)
		t.Log(resp.Body)
		assert.Equal(t, expectedStatus, resp.Result().StatusCode)
	}
}

func TestGetOrgs(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	req, _ := http.NewRequest("GET", "/orgs", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

func TestDeleteOrg(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Deleting an organization that does not exist
	req, _ := http.NewRequest("DELETE", "/orgs/1000", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 404, resp.Result().StatusCode)

	// Deleting a valid organization
	req, _ = http.NewRequest("DELETE", "/orgs/1", nil)
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}
