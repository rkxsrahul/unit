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
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/forgotpass"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// ForgotPassData is a  structure for binding data in body during forget or reset password request
type ForgotPassData struct {
	// state defines the state of request is it forgot or reset
	State string `json:"state" binding:"required"`
	// email of user
	Email string `json:"email"`
	// token recieved in email for resetting password
	Token string `json:"token"`
	// new password
	Password string `json:"password"`
	// WorkSpace
	Workspace string `json:"workspace"`
}

// ForgotPassEp is an api handler for forgot and reset password with workspace
func ForgotPassEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "forgot and reset password")
	// fetching data from request body
	span.LogKV("task", "binding body data")
	var fpdt ForgotPassData
	if err := c.BindJSON(&fpdt); err != nil {
		// if there is some error passing bad status code
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "State field is missing."})
		return
	}

	// when state is forget email is passed
	if fpdt.State == "forgot" {

		span.LogKV("task", "check email validation")
		if !methods.ValidateEmail(fpdt.Email) {
			span.LogKV("task", "send final output when email is wrong")
			// if there is some error passing bad status code
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
			return
		}

		span.LogKV("task", "call function for sending reset password link")
		msg, ok := forgotpass.ForgotPassChallenge(strings.ToLower(fpdt.Email), fpdt.Workspace, span)
		// return status code and msg and error if any
		span.LogKV("task", "forgot password done")
		c.JSON(http.StatusOK, gin.H{"error": !(ok), "message": msg})
		return
	}
	// when state is reset token and new password is passed
	if fpdt.State == "reset" {
		span.LogKV("task", "call function for updating password in database")
		email, msg, ok := forgotpass.ResetForgottenPass(fpdt.Token, fpdt.Password, fpdt.Workspace, span)
		if ok {
			span.LogKV("task", "save activity using auth-core schema")
			// recording user activity of reseting password
			activities.RecordActivity(database.Activities{Email: email,
				ActivityName: "reset_password",
				ClientIP:     c.ClientIP(),
				ClientAgent:  c.Request.Header.Get("User-Agent"),
				Timestamp:    time.Now().Unix()})
		}
		// return status code and msg and error if any
		span.LogKV("task", "reset password done")
		c.JSON(http.StatusOK, gin.H{"error": !(ok), "message": msg})
		return
	}
	span.LogKV("task", "send final output when state field is wrong")
	c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "State field value should be forgot or reset only."})
}
