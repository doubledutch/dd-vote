package viewcontrollers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/api/auth"

	"github.com/gin-gonic/gin"
)

// AdminViewController manages browser requests for admin pages
type AdminViewController struct {
	db gorm.DB
}

// NewAdminViewController creates a new instance
func NewAdminViewController(db gorm.DB) *AdminViewController {
	return &AdminViewController{db: db}
}

// ShowAdminPage either shows the login page or the admin panel, depending
// on whether the user is logged in and has access to the group
func (gc AdminViewController) ShowAdminPage(c *gin.Context) {
	groupUUID := c.Param("gid")
	if auth.IsLoggedIn(c) && auth.HasAccessToGroup(auth.GetUserIDFromCookie(c), groupUUID, gc.db) {
		// user has admin access
		http.ServeFile(c.Writer, c.Request, "views/admin_panel.html")
	} else {
		// show login page
		http.ServeFile(c.Writer, c.Request, "views/admin_login.html")
	}
}
