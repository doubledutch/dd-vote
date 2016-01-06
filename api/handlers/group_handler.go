package handlers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// GroupHandler manges api endpoints for groups
type GroupHandler struct {
	db gorm.DB
}

// NewGroupHandler creates a new instance
func NewGroupHandler(db gorm.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

// GetOrCreateGroup returns a group with the provided name, and creates one
// if it doesn't already exist
func (handler GroupHandler) GetOrCreateGroup(c *gin.Context) {
	var group table.Group
	c.Bind(&group)
	handler.db.FirstOrCreate(&group, table.Group{Name: group.Name})
	// TODO http response code should reflect get/create outcome
	c.JSON(http.StatusCreated, resp.APIResponse{IsError: false, Value: group})
}
