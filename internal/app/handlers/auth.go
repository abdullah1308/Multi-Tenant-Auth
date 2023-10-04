package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		userLogin := &models.UserLogin{}
		foundUser := &models.User{}

		err := c.BindJSON(&userLogin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.SetSearchPath(orgSearchPath)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not find organization. %s", err.Error())})
			return
		}

		org := &models.Organization{}
		orgName := userLogin.Organization
		err = h.DB.Where("name = ?", orgName).First(&org).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
			return
		}

		schema := h.GetOrgSchema(strconv.FormatUint(uint64(org.ID), 10))
		err = h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not login. %s", err.Error())})
			return
		}

		err = h.DB.Where("username = ?", userLogin.Username).First(foundUser).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password is incorrect"})
			return
		}

		passwordIsValid := helpers.VerifyPassword(*userLogin.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password is incorrect"})
			return
		}

		accessToken, refreshToken, err := helpers.GenerateAllTokens(foundUser.ID, *foundUser.Username, org.ID, *foundUser.UserType)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("refresh_token", refreshToken, 86400, "/", os.Getenv("CLIENT_DOMAIN"), false, true)

		response := gin.H{
			"id":        foundUser.ID,
			"username":  foundUser.Username,
			"org_id":    org.ID,
			"user_type": foundUser.UserType,
			"token":     accessToken,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
			return
		}

		claims, err := helpers.ValidateRefreshToken(refreshToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		foundUser := &models.User{}

		schema := h.GetOrgSchema(strconv.FormatUint(uint64(claims.OrgID), 10))
		err = h.SetSearchPath(schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not find user. %s", err.Error())})
			return
		}

		err = h.DB.Where("id = ?", claims.UserID).First(foundUser).Error

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		accessToken, _, err := helpers.GenerateAllTokens(foundUser.ID, *foundUser.Username, claims.OrgID, *foundUser.UserType)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to generate token. %s", err.Error())})
			return
		}

		response := gin.H{
			"id":        foundUser.ID,
			"username":  foundUser.Username,
			"org_id":    claims.OrgID,
			"user_type": foundUser.UserType,
			"token":     accessToken,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("refresh_token", "", -1, "/", os.Getenv("CLIENT_DOMAIN"), false, true)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}
