package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/database"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testOrgs = []string{"org1", "org2", "org3"}
var testUsers = []string{
	// username,password,user type,organization
	"user1,password1,ADMIN,org1",
	"user11,password1,USER,org1",
	"user111,password1,USER,org1",
	"user2,password2,ADMIN,org2",
	"user22,password2,USER,org2",
	"user222,password2,USER,org2",
	"user3,password3,ADMIN,org3",
	"user33,password3,USER,org3",
	"user333,password3,USER,org3",
}

func setupTestData() *Handler {
	err := godotenv.Load("../../../.env.testing")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbConfig := database.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := database.NewConnection(&dbConfig)

	if err != nil {
		log.Fatal("could not load the database")
	}

	handler := &Handler{
		DB: db,
	}

	err = handler.SetSearchPath(orgSearchPath)

	if err != nil {
		log.Fatal("error changing search path")
	}

	err = db.AutoMigrate(&models.Organization{})

	if err != nil {
		log.Fatal("could not migrate the db")
	}

	orgs := &[]models.Organization{}

	for _, name := range testOrgs {
		orgName := name
		org := models.Organization{
			Name: &orgName,
		}

		err = db.Save(&org).Error

		if err != nil {
			log.Fatal("error creating organization")
		}

		schemaName := handler.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))
		err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName)).Error

		if err != nil {
			log.Fatal("error creating schema")
		}

		*orgs = append(*orgs, org)
	}

	users := &[]models.User{}

	for _, user := range testUsers {
		userInfo := strings.Split(user, ",")

		password, err := helpers.HashPassword(userInfo[1])

		if err != nil {
			log.Fatal("error hashing password")
		}

		user := models.User{
			Username: &userInfo[0],
			Password: &password,
			UserType: &userInfo[2],
		}

		for _, org := range *orgs {

			if *org.Name != userInfo[3] {
				continue
			}

			schemaName := handler.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))
			err = handler.SetSearchPath(schemaName)

			if err != nil {
				log.Fatal("error changing search path")
			}

			err = db.AutoMigrate(&models.User{})

			if err != nil {
				log.Fatal("error migrating user table")
			}

			err = db.Save(&user).Error

			if err != nil {
				log.Fatal("error creating user")
			}

			break
		}

		*users = append(*users, user)
	}

	return handler
}

func teardownTestData(h *Handler) {
	err := h.SetSearchPath(orgSearchPath)

	if err != nil {
		log.Fatal("error finding organizations")
	}

	orgs := &[]models.Organization{}

	err = h.DB.Find(orgs).Error

	if err != nil {
		log.Fatal("error finding organizations")
	}

	for _, org := range *orgs {
		schema := h.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))

		err = h.DB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error

		if err != nil {
			log.Fatal("error deleting schema")
		}
	}

	err = h.DB.Exec("DROP TABLE IF EXISTS ORGANIZATIONS").Error

	if err != nil {
		log.Fatal("error deleting organization table")
	}
}

func sampleFromDB(h *Handler, userType string) (*models.User, *models.Organization, error) {
	org := &models.Organization{}
	user := &models.User{}

	err := h.SetSearchPath(orgSearchPath)

	if err != nil {
		return user, org, err
	}

	err = h.DB.First(org).Error

	if err != nil {
		return user, org, err
	}

	schema := h.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))
	err = h.SetSearchPath(schema)

	if err != nil {
		return user, org, err
	}

	err = h.DB.Where("user_type = ?", userType).First(&user).Error

	return user, org, err
}

func TestLogin(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// username,password,org,expected status
	tests := []string{"user1,password1,org10,404", "user2,password1,org1,401", "user1,password2,org1,401", "user1,password1,org1,200"}

	for _, test := range tests {
		testInfo := strings.Split(test, ",")

		userLoginData := models.UserLogin{
			Username:     &testInfo[0],
			Password:     &testInfo[1],
			Organization: &testInfo[2],
		}

		expectedStatus, _ := strconv.Atoi(testInfo[3])

		reqBody, err := json.Marshal(userLoginData)

		if err != nil {
			log.Panic("error marshalling test data to json")
		}

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		t.Log(resp.Result().StatusCode)
		t.Log(resp.Body)
		assert.Equal(t, expectedStatus, resp.Result().StatusCode)
	}
}

func TestLogout(t *testing.T) {
	handler := Handler{}

	router := gin.Default()
	handler.SetupRoutes(router)

	req, _ := http.NewRequest("GET", "/logout", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

func TestRefresh(t *testing.T) {
	handler := setupTestData()
	defer teardownTestData(handler)

	router := gin.Default()
	handler.SetupRoutes(router)

	// Testing without refresh cookie
	req, _ := http.NewRequest("GET", "/refresh", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 401, resp.Result().StatusCode)

	// Testing with refresh cookie for invalid user
	_, refreshToken, err := helpers.GenerateAllTokens(100, "test", 100, "USER")

	if err != nil {
		log.Panic("error generating refresh token")
	}

	cookie := &http.Cookie{Name: "refresh_token", Value: refreshToken}
	req, _ = http.NewRequest("GET", "/refresh", nil)
	req.AddCookie(cookie)

	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 404, resp.Result().StatusCode)

	// Testing with valid refresh cookie
	user, org, err := sampleFromDB(handler, "USER")

	if err != nil {
		log.Panic("error sampling from db")
	}

	_, refreshToken, err = helpers.GenerateAllTokens(user.ID, *user.Username, org.ID, *user.UserType)

	if err != nil {
		log.Panic("error generating refresh token")
	}

	cookie = &http.Cookie{Name: "refresh_token", Value: refreshToken}
	req, _ = http.NewRequest("GET", "/refresh", nil)
	req.AddCookie(cookie)

	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}
