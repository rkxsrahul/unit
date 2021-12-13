package api

import (
	"log"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
)

// ChangePasswordEp is a api handler to change password of a account
func ChangePasswordEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "change password")
	// defining type for fetching password from body of post request
	type NewPassword struct {
		CurrentPassword string `json:"current_password" form:"current_password" binding:"required"`
		Password        string `json:"password" form:"password" binding:"required"`
	}
	var newPass NewPassword
	// binding body json with above variable and checking error
	span.LogKV("task", "binding body data")
	err := c.BindJSON(&newPass)
	if err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Password is required field."})
		return
	}

	//extracting jwt claims for getting user id
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)
	// passing new password and userid and in return getting status code, msg and error
	span.LogKV("task", "call change password function")
	code, ok, msg := accounts.ChangePassword(claims["id"].(string), newPass.CurrentPassword, newPass.Password, span)
	span.LogKV("task", "send final output")
	c.JSON(code, gin.H{"error": !ok, "message": msg})
}

//==============================================================================
