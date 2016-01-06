package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/api/auth"

	"github.com/gin-gonic/gin"
)

// AdminController manages browser requests for admin pages
type AdminController struct {
	db gorm.DB
}

// NewAdminController creates a new instance
func NewAdminController(db gorm.DB) *AdminController {
	return &AdminController{db: db}
}

// ShowAdminPage either shows the login page or the admin panel, depending
// on whether the user is logged in and has access to the group
func (gc AdminController) ShowAdminPage(c *gin.Context) {
	gname := c.Param("gname")
	if auth.IsLoggedIn(c) && auth.HasAccessToGroup(auth.GetUserIDFromCookie(c), gname, gc.db) {
		// user has admin access
		http.ServeFile(c.Writer, c.Request, "views/admin_panel.html")
	} else {
		// show login page
		http.ServeFile(c.Writer, c.Request, "views/admin_login.html")
	}
}
