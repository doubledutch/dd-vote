package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/doubledutch/dd-vote/api/auth"
	"github.com/doubledutch/dd-vote/api/models/resp"
	"github.com/doubledutch/dd-vote/api/models/table"
	"github.com/gin-gonic/gin"
)

// ExportHandler manages api endpoints for exporting reports
type ExportHandler struct {
	db gorm.DB
}

// NewExportHandler creates a new instance
func NewExportHandler(db gorm.DB) *ExportHandler {
	return &ExportHandler{db: db}
}

// GetAllQuestionsCSV serves a csv files report of all the questions in a group
func (handler ExportHandler) GetAllQuestionsCSV(c *gin.Context) {
	gname := c.Param("gname")
	if !auth.HasAccessToGroup(auth.GetUserIDFromCookie(c), gname, handler.db) {
		c.JSON(http.StatusForbidden, resp.APIResponse{IsError: true, Message: "You don't have permission to access this group"})
		return
	}

	var group table.Group
	if err := handler.db.Where("name = ?", gname).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, resp.APIResponse{IsError: true, Message: "Group does not exist"})
		return
	}

	// get all posts for a group with comments and users for those comments
	var posts []table.Post
	handler.db.First(&group, table.Group{Name: gname})
	handler.db.Model(&group).Order("id").Preload("User").Association("Posts").Find(&posts)
	output := "Question,Upvotes,Downvotes,Created by" + "\n"
	for _, post := range posts {
		data := []string{post.Name, fmt.Sprintf("%v", post.Upvotes), fmt.Sprintf("%v", post.Downvotes), fmt.Sprintf("%v", post.User.FullName())}
		for i := range data {
			data[i] = escapeForCSV(data[i])
		}
		output += strings.Join(data[:], ",") + "\n"
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=questions.csv")
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	c.String(http.StatusOK, output)
}

func escapeForCSV(data string) string {
	if strings.ContainsAny(data, "\",") {
		data = strings.Replace(data, "\"", "\"\"", -1)
		return "\"" + data + "\""
	}
	return data
}
