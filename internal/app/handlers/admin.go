package handlers

import (
	"fmt"
	"net/http"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCreate := &models.UserCreate{}

		err := c.BindJSON(&userCreate)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		user := &models.User{
			Username: userCreate.Username,
			Password: userCreate.Password,
			UserType: userCreate.UserType,
		}

		err = validate.Struct(user)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		orgID, _ := c.Get("orgID")
		schema := h.GetOrgSchema(fmt.Sprintf("%v", orgID))
		err = h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		existingUser := &models.User{}
		err = h.DB.Where("username = ?", user.Username).First(&existingUser).Error
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		hashedPassword, err := helpers.HashPassword(*user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		user.Password = &hashedPassword

		err = h.DB.Create(&user).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		c.JSON(http.StatusCreated, models.EncodeToUserResponse(user))
	}
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not specified"})
			return
		}

		orgID, _ := c.Get("orgID")
		schema := h.GetOrgSchema(fmt.Sprintf("%v", orgID))
		err := h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not find user. %s", err.Error())})
			return
		}

		user := &models.User{}

		err = h.DB.Where("id = ?", id).First(&user).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find user. %s", err.Error())})
			return
		}

		userCreate := &models.UserCreate{}
		err = c.BindJSON(&userCreate)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if userCreate.Username != nil {
			user.Username = userCreate.Username
		}

		if userCreate.Password != nil {
			user.Password = userCreate.Password
		}

		if userCreate.UserType != nil {
			user.UserType = userCreate.UserType
		}

		err = validate.Struct(user)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if user.Password != nil {
			hashedPassword, err := helpers.HashPassword(*user.Password)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update user. %s", err.Error())})
				return
			}

			user.Password = &hashedPassword
		}

		err = h.DB.Save(&user).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update user. %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, models.EncodeToUserResponse(user))
	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not specified"})
			return
		}

		orgID, _ := c.Get("orgID")
		schema := h.GetOrgSchema(fmt.Sprintf("%v", orgID))
		err := h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not find user. %s", err.Error())})
			return
		}

		user := &models.User{}

		err = h.DB.Where("id = ?", id).First(&user).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find user. %s", err.Error())})
			return
		}

		err = h.DB.Delete(&user).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}
