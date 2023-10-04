package handlers

import (
	"fmt"

	"github.com/HousewareHQ/houseware---backend-engineering-octernship-abdullah1308/internal/app/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

const orgSearchPath = "public"

var validate = validator.New()

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) GetOrgSchema(orgID string) string {
	return "org_schema_" + orgID
}

func (h *Handler) SetSearchPath(schema string) error {
	err := h.DB.Exec(fmt.Sprintf("SET search_path = %s", schema)).Error
	return err
}

func (h *Handler) SetupRoutes(incomingRoutes *gin.Engine) {
	// Auth Routes
	incomingRoutes.POST("/login", h.Login())
	incomingRoutes.GET("/logout", h.Logout())
	incomingRoutes.GET("/refresh", h.Refresh())

	// Organization Routes
	incomingRoutes.POST("/orgs", h.CreateOrg())
	incomingRoutes.GET("/orgs", h.GetOrgs())
	incomingRoutes.DELETE("/orgs/:id", h.DeleteOrg())

	// User Routes
	incomingRoutes.GET("/users", middlewares.AuthorizeUser(), h.GetUsers())
	incomingRoutes.GET("/users/:id", middlewares.AuthorizeUser(), h.GetUserByID())

	// Admin Routes
	incomingRoutes.POST("/users", middlewares.AuthorizeUser(), middlewares.AuthorizeAdmin(), h.CreateUser())
	incomingRoutes.PATCH("/users/:id", middlewares.AuthorizeUser(), middlewares.AuthorizeAdmin(), h.UpdateUser())
	incomingRoutes.DELETE("/users/:id", middlewares.AuthorizeUser(), middlewares.AuthorizeAdmin(), h.DeleteUser())
}
