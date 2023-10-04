package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateOrg() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgCreate := &models.OrganizationCreate{}

		err := c.BindJSON(&orgCreate)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		err = validate.Struct(orgCreate)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		err = h.SetSearchPath(orgSearchPath)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create organization. %s", err.Error())})
			return
		}

		existingOrg := &models.Organization{}
		err = h.DB.Where("name = ?", orgCreate.Name).First(&existingOrg).Error
		if err == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "organization already exists"})
			return
		}

		org := &models.Organization{}
		org.Name = orgCreate.Name

		err = h.DB.Create(&org).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create organization. %s", err.Error())})
			return
		}

		schemaName := h.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))
		err = h.DB.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName)).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create organization. %s", err.Error())})
			return
		}

		orgUser := &models.User{}
		orgUser.Username = orgCreate.Username
		userType := "ADMIN"
		orgUser.UserType = &userType
		hashedPassword, err := helpers.HashPassword(*orgCreate.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		orgUser.Password = &hashedPassword

		err = h.SetSearchPath(schemaName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		err = h.DB.AutoMigrate(&models.User{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not migrate user table. %s", err.Error())})
			return
		}

		err = h.DB.Create(&orgUser).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not create user. %s", err.Error())})
			return
		}

		response := gin.H{
			"org_id":       org.ID,
			"organization": org.Name,
			"user_id":      orgUser.ID,
			"username":     orgUser.Username,
			"user_type":    "ADMIN",
		}

		c.JSON(http.StatusCreated, response)
	}
}

func (h *Handler) GetOrgs() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs := &[]models.Organization{}

		err := h.SetSearchPath(orgSearchPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not get organizations. %s", err.Error())})
			return
		}

		err = h.DB.Find(orgs).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not get organizations. %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, orgs)
	}
}

func (h *Handler) DeleteOrg() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := c.Param("id")

		if orgID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not specified"})
			return
		}

		err := h.SetSearchPath(orgSearchPath)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete organization. %s", err.Error())})
			return
		}

		org := &models.Organization{}
		err = h.DB.Where("id = ?", orgID).First(&org).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}

		schema := h.GetOrgSchema(orgID)
		err = h.DB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete organization. %s", err.Error())})
			return
		}

		err = h.DB.Delete(&org).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete organization. %s", err.Error())})
			return
		}

		response := gin.H{
			"message": "organization deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
