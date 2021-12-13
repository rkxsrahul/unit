package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/workspace"
)

// ForgotWorkspaceEp is an api handler for forgot workspace
// this will send a mail for recover workspace link
func ForgotWorkspaceEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "forgot workspace")

	// Forgot is a  structure for binding data in body during forget workspace request
	type Forgot struct {
		Email string `json:"email" binding:"required"`
	}

	//bind email from body
	span.LogKV("task", "binding body data")
	var email Forgot
	if err := c.BindJSON(&email); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please enter Email Address.",
		})
	}

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(email.Email) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	//send token in mail to recover workspaces
	span.LogKV("task", "call function for sending recover workspace link")
	err := workspace.Forgot(strings.ToLower(email.Email), span)
	if err != nil {
		span.LogKV("task", "any error in sending mail")
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "forgot workspace done")
	c.JSON(200, gin.H{
		"error":   false,
		"message": "We have emailed a special link to recover your workspaces. Please check your email.",
	})
}

//=====================================================================================//

// GetWorkspaceList is a api handler to fetch workspace list
// on basis of userid correspondance to that valid token
func GetWorkspaceList(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "login user with workspace")

	type TokenMeta struct {
		Token string `json:"token" binding:"required"`
	}
	//fetch token from request body
	span.LogKV("task", "binding body data")
	var tm TokenMeta
	if err := c.BindJSON(&tm); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Token is required."})
		return
	}
	span.LogKV("task", "call function to recover worksapce")
	code, mapd := workspace.RecoverWorkspace(tm.Token, span)
	span.LogKV("task", "send final output list")
	c.JSON(code, mapd)
}
