package controllers

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/models/req"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
	userHelper "github.com/jordanjoz/dd-vote/api/user"
)

type AdminController struct {
	db gorm.DB
}

func NewAdminController(db gorm.DB) *AdminController {
	return &AdminController{db: db}
}

func (ac AdminController) Login(c *gin.Context) {
	var userReq req.AdminLoginRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Unable to parse user: %s", err)
		c.JSON(200, resp.ApiResponse{IsError: true, Message: "Error logging in"})
		return
	}

	var user table.User
	user.Email = userReq.Email
	user.Password = userReq.Password

	// lookup user
	if err := ac.db.First(&user, table.User{Email: user.Email, Password: user.Password}).Error; err != nil {
		c.JSON(200, resp.ApiResponse{IsError: true, Message: "Email or password is incorrect"})
		return
	}

	if !userHelper.HasAccessToGroup(user.ID, userReq.GroupUUID, ac.db) {
		c.JSON(200, resp.ApiResponse{IsError: true, Message: "You don't have permission to access this group"})
		return
	}

	// set user logged in
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.JSON(200, resp.ApiResponse{IsError: false, Value: user})
}
