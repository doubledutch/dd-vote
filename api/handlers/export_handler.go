package handlers

import (
	"fmt"
	"log"
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

// GetAllQuestionsCSV serves a csv file report of all the questions in a group
// Columns: Question, Score (upvotes - downvotes), Total votes (upvotes + downvotes), Upvotes, Downvotes, Created by
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
	output := "Question,Score (upvotes - downvotes),Total votes (upvotes + downvotes),Upvotes,Downvotes,Created by" + "\n"
	for _, post := range posts {
		data := []string{post.Name, fmt.Sprintf("%v", int(post.Upvotes)-int(post.Downvotes)), fmt.Sprintf("%v", post.Upvotes+post.Downvotes), fmt.Sprintf("%v", post.Upvotes), fmt.Sprintf("%v", post.Downvotes), fmt.Sprintf("%v", post.User.FullName())}
		for i := range data {
			data[i] = escapeForCSV(data[i])
		}
		output += strings.Join(data[:], ",") + "\n"
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=questions.csv")
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	c.String(http.StatusOK, output)
}

// GetTopUsersCSV serves a csv file report of users with the most votes
// Columns: Name, Total votes
func (handler ExportHandler) GetTopUsersCSV(c *gin.Context) {
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

	handler.db.First(&group, table.Group{Name: gname})
	rows, err := handler.db.Table("users").Select("users.firstname, users.lastname, count(*) as total_votes").Joins("join votes on votes.user_id = users.id join posts on posts.id = votes.post_id join groups on groups.id = posts.group_id").Where("groups.id = ?", group.ID).Group("users.id, users.firstname, users.lastname").Rows()
	if err != nil {
		log.Printf("Unable to get top user rows: %s", err)
		c.JSON(http.StatusInternalServerError, resp.APIResponse{IsError: true, Message: "Unable to export top users"})
		return
	}
	output := "Name,Total votes" + "\n"
	for rows.Next() {
		var firstname string
		var lastname string
		var totalVotes int
		rows.Scan(&firstname, &lastname, &totalVotes)
		data := []string{firstname + " " + lastname, fmt.Sprintf("%v", totalVotes)}
		for i := range data {
			data[i] = escapeForCSV(data[i])
		}
		output += strings.Join(data[:], ",") + "\n"
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=top_users.csv")
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
