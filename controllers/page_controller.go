package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/doubledutch/dd-vote/api/models/table"

	"github.com/gin-gonic/gin"
)

// PageController manages browser requests related to group pages
type PageController struct {
	db gorm.DB
}

// NewPageController creates a new instance
func NewPageController(db gorm.DB) *PageController {
	return &PageController{db: db}
}

// ShowGroupPage shows the group page, which will use its front-end to make
// api requests for loading the appropriate data
func (gc PageController) ShowGroupPage(c *gin.Context) {
	// create new group if it does not exist
	gname := c.Param("gname")
	var group table.Group
	gc.db.FirstOrCreate(&group, table.Group{Name: gname})

	// show view
	http.ServeFile(c.Writer, c.Request, "views/group_page.html")
}
