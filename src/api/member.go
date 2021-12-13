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
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/member"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// MemberSignup is a structure for binding data from body during member signup request
type MemberSignup struct {
	Workspace string `json:"workspace" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

// MemberSignupEp is an api handler
// It is used to send invite link on mail if user is invited by the owner
func MemberSignupEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "member signup")
	//fetch email and workspace from database
	span.LogKV("task", "binding body data")
	var mem MemberSignup
	if err := c.BindJSON(&mem); err != nil {
		// if there is some error passing bad status code
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Email and Workspace are required fields."})
		return
	}

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(mem.Email) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	// send again invite mail to user
	span.LogKV("task", "call function for sending invite link to invited user")
	code, msg := member.MemberSignup(strings.ToLower(mem.Email), mem.Workspace, span)
	span.LogKV("task", "send final output")
	c.JSON(code, gin.H{
		"error":   code != 200,
		"message": msg,
	})
}

//========================================================================//

// TokenPassword is a structure for binding data from body during set new password request
type TokenPassword struct {
	Name      string `json:"name" binding:"required"`
	Contact   string `json:"contact" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Token     string `json:"token" binding:"required"`
	Workspace string `json:"workspace" binding:"required"`
}

// MemberRegistration is an api handler
// It is used for saving member information in database
// but before saving using invite token member is verified
func MemberRegistration(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "save member information")
	// fetching data from request body
	span.LogKV("task", "binding body data")
	var tp TokenPassword
	if err := c.BindJSON(&tp); err != nil {
		// if there is some error passing bad status code
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "name, contact, password, token and workspace are required field."})
		return
	}

	// check password is valid
	span.LogKV("task", "check password validation")
	if !methods.CheckPassword(tp.Password) {
		span.LogKV("task", "send final output when password is wrong")
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}

	// update password in database
	span.LogKV("task", "call function for updating member details in database")
	code, mapd := member.SaveMemberInfo(tp.Token, tp.Password, methods.HashForNewPassword(tp.Password), tp.Workspace, tp.Name, tp.Contact, span)
	if code == 200 {
		span.LogKV("task", "save activity using auth-core schema")
		// recording user activity of reseting password
		activities.RecordActivity(database.Activities{Email: mapd["email"].(string),
			ActivityName: "registration",
			ClientIP:     c.ClientIP(),
			ClientAgent:  c.Request.Header.Get("User-Agent"),
			Timestamp:    time.Now().Unix()})
	}
	span.LogKV("task", "reset password done")
	c.JSON(code, mapd)
}
