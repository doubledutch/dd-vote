package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	CommentController struct {
		db gorm.DB
	}
)

func NewCommentController(db gorm.DB) *CommentController {
	return &CommentController{db: db}
}

func (cc CommentController) CreateComment(c *gin.Context) {
	// lookup post by uuid
	postUUID := c.Query("post")
	var post models.Post
	if err := cc.db.Where("uuid = ?", postUUID).First(&post).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Question does not exist"})
		return
	}

	// deserialize comment
	var comment models.Comment
	c.Bind(&comment)
	comment.PostID = post.ID

	// create new comment
	if err := cc.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusOK, models.Error{IsError: true, Message: "Unknown error"})
		return
	}

	c.JSON(201, comment)
}
