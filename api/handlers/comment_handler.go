package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// CommentHandler manages api endpoints for comments
type CommentHandler struct {
	db gorm.DB
}

// NewCommentHandler creates a new instance
func NewCommentHandler(db gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// CreateComment creates a new comment on a post
func (handler CommentHandler) CreateComment(c *gin.Context) {
	// lookup post by uuid
	postUUID := c.Param("puuid")
	var post table.Post
	if err := handler.db.Where("uuid = ?", postUUID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Question does not exist"})
		return
	}

	// deserialize comment
	var comment table.Comment
	c.Bind(&comment)
	comment.PostID = post.ID
	comment.UserID = sessions.Default(c).Get("uid").(uint)

	// create new comment
	if err := handler.db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, resp.APIResponse{IsError: true, Message: "Unknown error"})
		return
	}

	// get comment that we just inserted with user info
	handler.db.Preload("User").Find(&comment)

	c.JSON(http.StatusCreated, resp.APIResponse{IsError: false, Value: comment})
}