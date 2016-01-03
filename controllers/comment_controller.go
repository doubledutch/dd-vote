package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
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
		c.JSON(http.StatusOK, models.ApiResponse{IsError: true, Message: "Question does not exist"})
		return
	}

	// deserialize comment
	var comment models.Comment
	c.Bind(&comment)
	comment.PostID = post.ID
	comment.UserID = sessions.Default(c).Get("uid").(uint)

	// create new comment
	if err := cc.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusOK, models.ApiResponse{IsError: true, Message: "Unknown error"})
		return
	}

	// get comment that we just inserted with user info
	cc.db.Preload("User").Find(&comment)

	c.JSON(201, models.ApiResponse{IsError: false, Value: comment})
}
