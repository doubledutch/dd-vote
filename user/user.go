package user

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/models"
)

// User roles
const (
	SuperAdmin = 0
	EventAdmin = 1
	GroupAdmin = 2
)

func HasAccessToGroup(uid uint, groupUUID string, db gorm.DB) bool {
	var permissions []models.Permission
	db.Where("user_id = ?", uid).Find(&permissions)
	for _, p := range permissions {
		if p.Role == SuperAdmin {
			// SuperAdmin
			return true
		} else if p.Role == GroupAdmin && p.Metadata == groupUUID {
			// GroupAdmin
			return true
		}
	}
	return false
}

func IsLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	return session.Get("uid") != nil
}

func GetUserIDFromCookie(c *gin.Context) uint {
	session := sessions.Default(c)
	uid := session.Get("uid").(uint)
	return uid
}

func StoreUserIDInCookie(c *gin.Context, userID uint) {
	session := sessions.Default(c)
	session.Set("uid", userID)
	session.Save()
}

func ClearUserIDFromCookie(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("uid", nil)
	session.Save()
}
