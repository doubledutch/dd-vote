package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// PostController manages api endpoints for posts (questions)
type PostController struct {
	db gorm.DB
}

// NewPostController creates a new instance
func NewPostController(db gorm.DB) *PostController {
	return &PostController{db: db}
}

// GetAllPostsForGroup returns all the questions in a group with nested
// data for their comments on the users on those comments
func (pc PostController) GetAllPostsForGroup(c *gin.Context) {
	groupname := c.Query("group")
	var group table.Group
	if err := pc.db.Where("name = ?", groupname).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	// get all posts for a group with comments and users for those comments
	var posts []table.Post
	pc.db.First(&group, table.Group{Name: group.Name})
	pc.db.Model(&group).Order("id").Preload("Comments").Preload("Comments.User").Association("Posts").Find(&posts)

	// make sure comments are empty slice and not nil
	for i := range posts {
		if posts[i].Comments == nil {
			posts[i].Comments = make([]table.Comment, 0)
		}
	}

	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Value: posts})
}

// CreatePost creates a new question
func (pc PostController) CreatePost(c *gin.Context) {
	// lookup group by name
	groupname := c.Query("group")
	var group table.Group
	if err := pc.db.Where("name = ?", groupname).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	// deserialize post
	var post table.Post
	c.Bind(&post)
	post.GroupID = group.ID
	post.UUID = uuid.NewV4().String() //TODO make sure this doesn't break everything
	post.CreatedBy = sessions.Default(c).Get("uid").(uint)

	// create new question
	if err := pc.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusConflict, resp.APIResponse{IsError: true, Message: "Question has already been asked"})
		return
	}

	// make sure comments are empty slice and not nil
	post.Comments = make([]table.Comment, 0)

	c.JSON(http.StatusCreated, resp.APIResponse{IsError: false, Value: post})
}
