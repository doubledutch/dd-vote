package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/jordanjoz/dd-vote/api/models/table"

	"github.com/gin-gonic/gin"
)

// PageViewController manages browser requests related to group pages
type PageViewController struct {
	db gorm.DB
}

// NewPageViewController creates a new instance
func NewPageViewController(db gorm.DB) *PageViewController {
	return &PageViewController{db: db}
}

// ShowGroupPage shows the group page, which will use its front-end to make
// api requests for loading the appropriate data
func (gc PageViewController) ShowGroupPage(c *gin.Context) {
	// create new group if it does not exist
	gname := c.Param("gname")
	var group table.Group
	gc.db.FirstOrCreate(&group, table.Group{Name: gname})

	// show view
	http.ServeFile(c.Writer, c.Request, "views/group_page.html")
}
