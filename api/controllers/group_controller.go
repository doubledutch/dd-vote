package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

type GroupController struct {
	db gorm.DB
}

func NewGroupController(db gorm.DB) *GroupController {
	return &GroupController{db: db}
}

func (gc GroupController) GetOrCreateGroup(c *gin.Context) {
	var group table.Group
	c.Bind(&group)
	gc.db.FirstOrCreate(&group, table.Group{Name: group.Name})
	// TODO http response code should reflect get/create outcome
	c.JSON(http.StatusCreated, resp.ApiResponse{IsError: false, Value: group})
}
