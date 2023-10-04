package middlewares

import (
	"net/http"
	"strings"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/helpers"
	"github.com/gin-gonic/gin"
)

func AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")

		var tokenString string
		if strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			tokenString = authHeader[len(BEARER_SCHEMA):]
		} else {
			tokenString = ""
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header provided"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateAccessToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("orgID", claims.OrgID)
		c.Set("userType", claims.UserType)
		c.Next()
	}
}

func AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Keys["userType"] != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "only admins can access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
