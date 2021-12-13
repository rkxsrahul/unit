package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/activities"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// EmailVerifyToken defining structure for binding verification mail data
type EmailVerifyToken struct {
	VerificationCode string `json:"verification_code" binding:"required"`
	Email            string `json:"email" binding:"required"`
}

// VerifyMailEp is a api handler for verify email id by token
func VerifyMailEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	} 
	span.SetTag("event", "verify mail")
	// binding request body data
	span.LogKV("task", "binding body data")
	var tokendata EmailVerifyToken
	if c.BindJSON(&tokendata) != nil {
		span.LogKV("task", "send final output when request data is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Email and Verification Code are required field."})
		return
	}

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(tokendata.Email) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	// passing email and token for getting verified
	span.LogKV("task", "call verify mail function")
	account, ok := accounts.VerifyMail(strings.ToLower(tokendata.Email), tokendata.VerificationCode, span)
	if !ok {
		span.LogKV("task", "send final output invalid or expired code")
		// if there is some error then passing StatusUnauthorized and msg invalid token
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid or expired Verification Code."})
		return
	}

	//=============================================

	span.LogKV("task", "save activity")
	// saving user-activity
	activity := database.Activities{Email: tokendata.Email,
		ClientIP:     c.ClientIP(),
		ClientAgent:  c.Request.Header.Get("User-Agent"),
		Timestamp:    time.Now().Unix(),
		ActivityName: "email_verified"}
	activities.RecordActivity(activity)
	//=============================================

	// setting jwt token and claims to be used in other protected apis
	span.LogKV("task", "generate jwt token")
	mapd := jwtToken.JwtToken(account, span)
	mapd["name"] = account.Name
	mapd["email"] = account.Email
	mapd["role_id"] = account.RoleID
	mapd["error"] = false
	mapd["message"] = "Email verification done"
	span.LogKV("task", "Email verification done")
	c.JSON(200, mapd)
}

// ChangeMail is an api handler to toggle mail service
func ChangeMail(c *gin.Context) {
	mail.ToggleMail(c.Param("value"))
}

// ChangeOTP is an api handler to toggle OTP
func ChangeOTP(c *gin.Context) {
	verifyToken.ToggleOTP(c.Param("value"))
}
