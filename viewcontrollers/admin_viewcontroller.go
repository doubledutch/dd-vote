package viewcontrollers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/user"

	"github.com/gin-gonic/gin"
)

type (
	AdminViewController struct {
		db gorm.DB
	}
)

func NewAdminViewController(db gorm.DB) *AdminViewController {
	return &AdminViewController{db: db}
}

func (gc AdminViewController) ShowAdminPage(c *gin.Context) {
	groupUUID := c.Param("gid")
	if user.IsLoggedIn(c) && user.HasAccessToGroup(user.GetUserIDFromCookie(c), groupUUID, gc.db) {
		// user has admin access
		http.ServeFile(c.Writer, c.Request, "views/admin_panel.html")
	} else {
		// show login page
		http.ServeFile(c.Writer, c.Request, "views/admin_login.html")
	}
}