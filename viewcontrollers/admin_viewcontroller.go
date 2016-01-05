package viewcontrollers

import (
	"net/http"

	"github.com/jinzhu/gorm"

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
	http.ServeFile(c.Writer, c.Request, "views/admin_login.html")
}
