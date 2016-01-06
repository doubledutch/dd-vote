package controllers

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/auth"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

type ExportController struct {
	db gorm.DB
}

func NewExportController(db gorm.DB) *ExportController {
	return &ExportController{db: db}
}

func (ec ExportController) GetAllQuestionsCSV(c *gin.Context) {
	groupUUID := c.Param("gid")
	if !auth.HasAccessToGroup(auth.GetUserIDFromCookie(c), groupUUID, ec.db) {
		c.JSON(http.StatusForbidden, resp.ApiResponse{IsError: true, Message: "You don't have permission to access this group"})
		return
	}

	var group table.Group
	if err := ec.db.Where("name = ?", groupUUID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.ApiResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	// get all posts for a group with comments and users for those comments
	var posts []table.Post
	ec.db.First(&group, table.Group{Name: groupUUID})
	ec.db.Model(&group).Order("id").Preload("Comments").Preload("Comments.User").Association("Posts").Find(&posts)

	output := "Question,Upvotes,Downvotes,Created by" + "\n"
	for _, post := range posts {
		data := []string{post.Name, "strconv.FormatUint(post.Upvotes, 10)", "strconv.FormatUint(post.Downvotes, 10)", "strconv.FormatUint(post.CreatedBy, 16)"}
		//TODO function for escaping commas, quotes, and new lines
		output += strings.Join(data[:], ",") + "\n"
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=top.csv")
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	c.String(http.StatusOK, output)
}
