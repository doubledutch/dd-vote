package controllers

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	UserController struct {
		db gorm.DB
	}
)

func NewUserController(db gorm.DB) *UserController {
	return &UserController{db: db}
}

func (uc UserController) LoginWithClientID(c *gin.Context) {
	// deserialize post
	var userReq models.UserRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Unable to parse user: %s", err)
		c.JSON(200, models.ApiResponse{IsError: true, Message: "Error logging in"})
		return
	}

	var user models.User
	user.ClientID = userReq.ClientID
	user.Firstname = userReq.Firstname
	user.Lastname = userReq.Lastname

	// create or get user
	if err := uc.db.FirstOrCreate(&user, models.User{ClientID: user.ClientID}).Error; err != nil {
		c.JSON(200, models.ApiResponse{IsError: true, Message: "Error logging in"})
		return
	}

	// set user logged in
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()

	c.JSON(200, models.ApiResponse{IsError: false, Value: user})
}
