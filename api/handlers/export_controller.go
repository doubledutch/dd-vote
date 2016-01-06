package handlers

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jordanjoz/dd-vote/api/auth"
	"github.com/jordanjoz/dd-vote/api/models/resp"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// ExportController manages api endpoints for exporting reports
type ExportController struct {
	db gorm.DB
}

// NewExportController creates a new instance
func NewExportController(db gorm.DB) *ExportController {
	return &ExportController{db: db}
}

// GetAllQuestionsCSV serves a csv files report of all the questions in a group
func (ec ExportController) GetAllQuestionsCSV(c *gin.Context) {
	gname := c.Param("gname")
	if !auth.HasAccessToGroup(auth.GetUserIDFromCookie(c), gname, ec.db) {
		c.JSON(http.StatusForbidden, resp.APIResponse{IsError: true, Message: "You don't have permission to access this group"})
		return
	}

	var group table.Group
	if err := ec.db.Where("name = ?", gname).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	// get all posts for a group with comments and users for those comments
	var posts []table.Post
	ec.db.First(&group, table.Group{Name: gname})
	ec.db.Model(&group).Order("id").Preload("Comments").Preload("Comments.User").Association("Posts").Find(&posts)

	output := "Question,Upvotes,Downvotes,Created by" + "\n"
	for _, post := range posts {
		data := []string{post.Name, "strconv.FormatUint(post.Upvotes, 10)", "strconv.FormatUint(post.Downvotes, 10)", "strconv.FormatUint(post.CreatedBy, 16)"}
		//TODO function for escaping commas, quotes, and new lines
		output += strings.Join(data[:], ",") + "\n"
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=questions.csv")
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	c.String(http.StatusOK, output)
}
