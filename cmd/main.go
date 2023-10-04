package main

import (
	"log"
	"os"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/database"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/handlers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
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

	err = db.AutoMigrate(&models.Organization{})

	if err != nil {
		log.Fatal("could not migrate the db")
	}

	handler := &handlers.Handler{
		DB: db,
	}

	router := gin.Default()

	handler.SetupRoutes(router)

	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "3000"
	}

	err = router.Run(":" + port)

	if err != nil {
		log.Fatalf("could not start server. %s", err.Error())
	}
}
