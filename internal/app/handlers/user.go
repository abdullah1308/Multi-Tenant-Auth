package handlers

import (
	"fmt"
	"net/http"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users := &[]models.User{}

		orgID, _ := c.Get("orgID")
		schema := h.GetOrgSchema(fmt.Sprintf("%v", orgID))
		err := h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get users"})
			return
		}

		err = h.DB.Find(users).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get users"})
			return
		}

		usersRes := &[]models.UserResponse{}

		for _, user := range *users {
			*usersRes = append(*usersRes, *models.EncodeToUserResponse(&user))
		}

		c.JSON(http.StatusOK, usersRes)
	}
}

func (h *Handler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.User{}

		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not specified"})
			return
		}

		orgID, _ := c.Get("orgID")
		schema := h.GetOrgSchema(fmt.Sprintf("%v", orgID))
		err := h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not get the user. %s", err.Error())})
			return
		}

		err = h.DB.Where("id = ?", id).First(user).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not get the user. %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, models.EncodeToUserResponse(user))
	}
}
