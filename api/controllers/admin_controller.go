package controllers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/auth"
	"github.com/jordanjoz/dd-vote/api/models/req"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// AdminController manages api endpoints for admins
type AdminController struct {
	db gorm.DB
}

// NewAdminController creates a new instance
func NewAdminController(db gorm.DB) *AdminController {
	return &AdminController{db: db}
}

// Login attempts to log an admin in
func (ac AdminController) Login(c *gin.Context) {
	var userReq req.AdminLoginRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Unable to parse user: %s", err)
		c.JSON(http.StatusBadRequest, resp.APIResponse{IsError: true, Message: "Error logging in"})
		return
	}

	// create user object from request
	user := userReq.ToUser()

	// lookup user in db
	// TODO user passwords should be hashed
	if err := ac.db.First(&user, table.User{Email: user.Email, Password: user.Password}).Error; err != nil {
		c.JSON(http.StatusBadRequest, resp.APIResponse{IsError: true, Message: "Email or password is incorrect"})
		return
	}

	if !auth.HasAccessToGroup(user.ID, userReq.GroupUUID, ac.db) {
		c.JSON(http.StatusForbidden, resp.APIResponse{IsError: true, Message: "You don't have permission to access this group"})
		return
	}

	// set user logged in
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Value: user})
}
