package controllers

import (
	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/models"
)

type (
	GroupController struct {
		db gorm.DB
	}
)

func NewGroupController(db gorm.DB) *GroupController {
	return &GroupController{db: db}
}

func (gc GroupController) GetOrCreateGroup(c *gin.Context) {
	var group models.Group
	c.Bind(&group)
	gc.db.FirstOrCreate(&group, models.Group{Name: group.Name})
	c.JSON(201, group)
}
