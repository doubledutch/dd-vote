package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

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
	http.ServeFile(c.Writer, c.Request, "views/group_page.html")
}
