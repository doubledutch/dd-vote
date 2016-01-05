package viewcontrollers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/jordanjoz/dd-vote/api/models/table"

	"github.com/gin-gonic/gin"
)

type (
	PageViewController struct {
		db gorm.DB
	}
)

func NewPageViewController(db gorm.DB) *PageViewController {
	return &PageViewController{db: db}
}

func (gc PageViewController) ShowGroupPage(c *gin.Context) {
	// create new group if it does not exist
	gid := c.Param("gid")
	var group table.Group
	gc.db.FirstOrCreate(&group, table.Group{Name: gid})

	// show view
	http.ServeFile(c.Writer, c.Request, "views/group_page.html")
}
