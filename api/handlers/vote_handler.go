package handlers

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/doubledutch/dd-vote/api/auth"
	"github.com/doubledutch/dd-vote/api/models/req"
	"github.com/doubledutch/dd-vote/api/models/resp"
	"github.com/doubledutch/dd-vote/api/models/table"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// VoteHandler manages api endpoints for voting
type VoteHandler struct {
	db gorm.DB
}

// NewVoteHandler creates a new instance
func NewVoteHandler(db gorm.DB) *VoteHandler {
	return &VoteHandler{db: db}
}

// GetUserVotes gets the user's votes for a group
func (handler VoteHandler) GetUserVotes(c *gin.Context) {
	gname := c.Param("gname")
	var group table.Group
	if err := handler.db.Where("name = ?", gname).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	userID := auth.GetUserIDFromCookie(c)
	var votes []table.Vote
	if err := handler.db.Joins("left join posts on posts.id = votes.post_id").Where("posts.group_id = ? and votes.user_id = ?", group.ID, userID).Find(&votes).Error; err != nil {
		// return empty slice instead of nil for no data
		votes = make([]table.Vote, 0)
	}

	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Value: votes})
}

// CreateOrUpdateVote create a new vote or update the user's existing one
func (handler VoteHandler) CreateOrUpdateVote(c *gin.Context) {
	puuid := c.Param("puuid")
	var voteReq req.VoteCreateRequest
	if err := c.BindJSON(&voteReq); err != nil {
		log.Printf("Unable request: %s", err)
		c.JSON(http.StatusBadRequest, resp.APIResponse{IsError: true, Message: "Error parsing request"})
		return
	}

	// lookup post by uuid
	var post table.Post
	if err := handler.db.Where("uuid = ?", puuid).First(&post).Error; err != nil {
		c.JSON(http.StatusOK, resp.APIResponse{IsError: true, Message: "Question does not exist"})
		return
	}

	// get current user id
	userID := sessions.Default(c).Get("uid").(uint)

	// start transaction
	tx := handler.db.Begin()

	// attempt to get existing vote
	var vote table.Vote
	isChangingVote := false
	if err := tx.Where("user_id = ? AND post_id = ?", userID, post.ID).First(&vote).Error; err != nil {
		// vote does not exist
		vote.PostID = post.ID
		vote.PostUUID = post.UUID
		vote.UserID = userID
		vote.Value = voteReq.Value
		if err := tx.Create(&vote).Error; err != nil {
			log.Printf("Unable to create vote: %s", err)
			c.JSON(http.StatusInternalServerError, resp.APIResponse{IsError: true, Message: "Unable to create vote"})
			return
		}
	} else {
		// don't allow user to vote same way multiple times
		if vote.Value == voteReq.Value {
			c.JSON(http.StatusConflict, resp.APIResponse{IsError: true, Message: "Already voted that way"})
			return
		}

		// changing vote
		vote.Value = voteReq.Value
		isChangingVote = true
		tx.Save(&vote)
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

	c.JSON(http.StatusOK, resp.APIResponse{IsError: false, Value: post})
}
