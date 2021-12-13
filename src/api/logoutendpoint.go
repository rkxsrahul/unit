package api

import (
	"log"
	"strings"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/redisdb"
)

// Logout is a  api handler for logging out user means delete token or session detail from redis
func Logout(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "logout user")
	// fetch token from header
	span.LogKV("task", "fetch token from request header")
	token := c.Request.Header.Get("Authorization")
	// trim bearer from token
	token = strings.TrimPrefix(token, "Bearer ")
	// call delete token go function
	span.LogKV("task", "delete token from redis")
	err := redisdb.DeleteToken(token)
	if err != nil {
		span.LogKV("task", "send final output when error in deleting token")
		c.JSON(501, gin.H{
			"error":   err,
			"message": "Error in deleting token",
		})
	} else {
		// delete token from db means delete session from db
		go jwtToken.DeleteTokenFromDb(token)
		span.LogKV("task", "send final output")
		c.JSON(200, gin.H{
			"error":   false,
			"message": "Successfully logout",
		})
	}
}
