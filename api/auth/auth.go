package auth

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jordanjoz/dd-vote/api/models/table"
)

// User roles
const (
	SuperAdmin = 0
	EventAdmin = 1
	GroupAdmin = 2
)

// HasAccessToGroup checks if a user is an admin for a group
func HasAccessToGroup(uid uint, groupUUID string, db gorm.DB) bool {
	var permissions []table.Permission
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

// IsLoggedIn checks if the current user is logged in
func IsLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	return session.Get("uid") != nil
}

// GetUserIDFromCookie gets the userId for the current user
func GetUserIDFromCookie(c *gin.Context) uint {
	session := sessions.Default(c)
	uid := session.Get("uid").(uint)
	return uid
}

// StoreUserIDInCookie stores the userId for teh current user
func StoreUserIDInCookie(c *gin.Context, userID uint) {
	session := sessions.Default(c)
	session.Set("uid", userID)
	session.Save()
}

// ClearUserIDFromCookie clears the userId for the current suer
func ClearUserIDFromCookie(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("uid", nil)
	session.Save()
}
