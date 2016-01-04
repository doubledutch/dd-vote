package controllers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	VoteController struct {
		db gorm.DB
	}
)

func NewVoteController(db gorm.DB) *VoteController {
	return &VoteController{db: db}
}

func (cc VoteController) CreateOrUpdateVote(c *gin.Context) {
	var voteReq models.VoteCreateRequest
	if err := c.BindJSON(&voteReq); err != nil {
		log.Printf("Unable request: %s", err)
		c.JSON(200, models.ApiResponse{IsError: true, Message: "Error parsing request"})
		return
	}

	// lookup post by uuid
	var post models.Post
	if err := cc.db.Where("uuid = ?", voteReq.PostUUID).First(&post).Error; err != nil {
		c.JSON(http.StatusOK, models.ApiResponse{IsError: true, Message: "Question does not exist"})
		return
	}

	// get current user id
	userID := sessions.Default(c).Get("uid").(uint)

	// start transaction
	tx := cc.db.Begin()

	// attempt to get existing vote
	var vote models.Vote
	isChangingVote := false
	if err := tx.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&vote).Error; err != nil {
		// vote does not exist
		vote.PostID = post.ID
		vote.UserID = userID
		vote.Value = voteReq.Value
		if err := tx.Create(&vote).Error; err != nil {
			log.Printf("Unable to create vote: %s", err)
			c.JSON(http.StatusOK, models.ApiResponse{IsError: true, Message: "Unable to create vote"})
			return
		}
	} else {
		// don't allow user to vote same way multiple times
		if vote.Value == voteReq.Value {
			c.JSON(http.StatusOK, models.ApiResponse{IsError: true, Message: "Already voted that way"})
			return
		} else {
			// changing vote
			vote.Value = voteReq.Value
			isChangingVote = true
			tx.Save(&vote)
		}
	}

	// update question's vote counts
	if vote.Value == 1 {
		post.Upvotes++
		if isChangingVote {
			post.Downvotes--
		}
	} else {
		post.Downvotes++
		if isChangingVote {
			post.Upvotes--
		}
	}

	tx.Save(&post)

	// commit transaction
	tx.Commit()

	c.JSON(201, models.ApiResponse{IsError: false, Value: post})
}
