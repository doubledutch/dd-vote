package middleware

import (
	"net/http"

	"github.com/doubledutch/dd-vote/api/auth"
	"github.com/doubledutch/dd-vote/api/models/resp"
	"github.com/gin-gonic/gin"
)

// UseAuth rejects unauthorized api requests
func UseAuth(c *gin.Context) {
	if !auth.IsLoggedIn(c) {
		c.JSON(http.StatusUnauthorized, resp.APIResponse{IsError: false, Message: "User is not logged in"})
		c.Abort()
	}
}
