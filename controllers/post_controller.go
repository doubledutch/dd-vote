package controllers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	// PostController represents the controller for operating on the User resource
	PostController struct {
		db gorm.DB
	}
)

func NewPostController(db gorm.DB) *PostController {
	return &PostController{db: db}
}

// GetPost retrieves an individual user resource
func (pc PostController) GetPost(c *gin.Context) {
	// Stub an example user
	u := models.Post{
		Name: "Why do we ask questions 2?",
	}

	// Write content-type, statuscode, payload
	c.JSON(http.StatusOK, u)
}

// CreatePost creates a new user resource
func (pc PostController) CreatePost(c *gin.Context) {
	//TODO lookup group id by name
	groupname := c.Query("groupname")
	var group models.Group
	pc.db.Where("name = ?", groupname).First(&group)

	log.Println(group.ID)

	var post models.Post
	c.Bind(&post)
	post.GroupID = group.ID

	//TODO lookup group id by name

	if err := pc.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Question has already been asked"})
		return
	}

	c.JSON(201, post)
}
