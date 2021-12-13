package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/activities"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/login"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// TokenLogin is an api handler used to login with token(token used in fetching workspace list)
func TokenLogin(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "login user with token")

	type ToLogin struct {
		Token     string `json:"token" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Workspace string `json:"workspace" binding:"required"`
	}

	//fetch token, email and WorkspaceID from request body
	span.LogKV("task", "binding body data")
	var data ToLogin
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass token, email and workspace",
		})
		return
	}

	span.LogKV("task", "call function to login with token")
	code, mapd := login.TokenLogin(strings.ToLower(data.Email), data.Token, data.Workspace, span)
	if code == 200 {
		span.LogKV("task", "save activity using auth-core schema")
		// recording user activity of login with team
		activities.RecordActivity(database.Activities{Email: data.Email,
			ActivityName: "Login with token",
			ClientIP:     c.ClientIP(),
			ClientAgent:  c.Request.Header.Get("User-Agent"),
			Timestamp:    time.Now().Unix()})
	}

	span.LogKV("task", "send final output")
	c.JSON(code, mapd)
}

//===========================================================================//

// LoginEndpoint is an api handler used to login with workspace
func LoginEndpoint(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "login user with workspace")

	type Login struct {
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Workspace string `json:"workspace"`
	}

	//fetch password, email and WorkspaceID from request body
	span.LogKV("task", "binding body data")
	var data Login
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass password and email",
		})
		return
	}

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(data.Email) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	//=============================================
	// recording user activity
	activity := database.Activities{Email: data.Email,
		ClientIP:    c.ClientIP(),
		ClientAgent: c.Request.Header.Get("User-Agent"),
		Timestamp:   time.Now().Unix()}
	//=============================================

	span.LogKV("task", "call function to login with workspace")
	code, mapd := login.NormalLogin(strings.ToLower(data.Email), data.Password, data.Workspace, span)
	if code == 200 {
		span.LogKV("task", "save login activity using auth-core schema")
		// recording user activity of login
		activity.ActivityName = "login"
		activities.RecordActivity(activity)
	} else if code == 500 {
		span.LogKV("task", "save failedlogin activity using auth-core schema")
	} else {
		span.LogKV("task", "save failedlogin activity using auth-core schema")
		// recording user activity of failed login
		activity.ActivityName = "failedlogin"
		activities.RecordActivity(activity)
	}
	span.LogKV("task", "send final output")
	c.JSON(code, mapd)
	return
}
