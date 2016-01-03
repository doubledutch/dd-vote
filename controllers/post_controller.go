package controllers

import (
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

func (pc PostController) GetAllPostsForGroup(c *gin.Context) {
	groupname := c.Query("group")
	var group models.Group
	if err := pc.db.Where("name = ?", groupname).First(&group).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Group does not exist"})
		return
	}

	var posts []models.Post
	pc.db.First(&group, models.Group{Name: group.Name})
	pc.db.Model(&group).Association("Posts").Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// CreatePost creates a new user resource
func (pc PostController) CreatePost(c *gin.Context) {
	// lookup group by name
	groupname := c.Query("group")
	var group models.Group
	if err := pc.db.Where("name = ?", groupname).First(&group).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Group does not exist"})
		return
	}

	// deserialize post
	var post models.Post
	c.Bind(&post)
	post.GroupID = group.ID

	// create new question
	if err := pc.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Question has already been asked"})
		return
	}

	c.JSON(201, post)
}
