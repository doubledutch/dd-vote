package handlers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/auth"
	"github.com/jordanjoz/dd-vote/api/models/req"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// UserController manages api endpoints for users
type UserController struct {
	db gorm.DB
}

// NewUserController creates a new instance
func NewUserController(db gorm.DB) *UserController {
	return &UserController{db: db}
}

// LoginWithClientID logs a user in using the data from the client, which
// gives them permission to post a question, comment, and vote in a group
func (uc UserController) LoginWithClientID(c *gin.Context) {
	// deserialize post
	var userReq req.UserRequest
	if err := c.BindJSON(&userReq); err != nil {
		log.Printf("Unable to parse user: %s", err)
		c.JSON(http.StatusBadRequest, resp.APIResponse{IsError: true, Message: "Error logging in"})
		return
	}

	// create user object from request
	user := userReq.ToUser()

	// create or get user from db
	if err := uc.db.FirstOrCreate(&user, table.User{ClientID: user.ClientID}).Error; err != nil {
		c.JSON(http.StatusBadRequest, resp.APIResponse{IsError: true, Message: "Error logging in"})
		return
	}

	// set user logged in
	auth.StoreUserIDInCookie(c, user.ID)

	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Value: user})
}

// Logout logs a user out
func (uc UserController) Logout(c *gin.Context) {
	auth.ClearUserIDFromCookie(c)
	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Message: "User logged out"})
}
