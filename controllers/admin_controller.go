package controllers

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	AdminController struct {
		db gorm.DB
	}
)

func NewAdminController(db gorm.DB) *AdminController {
	return &AdminController{db: db}
}

func (ac AdminController) Login(c *gin.Context) {
	var userReq models.AdminLoginRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Unable to parse user: %s", err)
		c.JSON(200, models.ApiResponse{IsError: true, Message: "Error logging in"})
		return
	}

	var user models.User
	user.Email = userReq.Email
	user.Password = userReq.Password

	// lookup user
	if err := ac.db.First(&user, models.User{Email: user.Email, Password: user.Password}).Error; err != nil {
		c.JSON(200, models.ApiResponse{IsError: true, Message: "Error logging in"})
		return
	}

	// set user logged in
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.JSON(200, models.ApiResponse{IsError: false, Value: user})
}
