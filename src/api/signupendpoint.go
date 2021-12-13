package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/signup"
)

// SignupData defining structure for binding signup data
type SignupData struct {
	Name     string `json:"name"`
	Contact  string `json:"contact"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignupEndpoint is a api handler for creating accounts
func SignupEndpoint(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "signup user")
	// binding body data
	span.LogKV("task", "binding body data")
	var signupdt SignupData
	if err := c.BindJSON(&signupdt); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "email, password and name are required fields."})
		return
	}

	//saving data in account structure
	acs := database.Accounts{}
	acs.Name = signupdt.Name
	//validation check on email
	span.LogKV("task", "validating email")
	if !methods.ValidateEmail(strings.ToLower(signupdt.Email)) {
		span.LogKV("task", "send final output when email is invalid")
		c.JSON(400, gin.H{"error": true, "message": "Please pass valid email address"})
		return
	}
	acs.Email = strings.ToLower(signupdt.Email)
	acs.ContactNo = signupdt.Contact
	acs.VerifyStatus = "not_verified"
	acs.CreationDate = time.Now().Unix()

	//validation check on password
	span.LogKV("task", "validating password")
	if !methods.CheckPassword(signupdt.Password) {
		span.LogKV("task", "send final output when password is invalid")
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}

	// save hash password insted of normal password
	span.LogKV("task", "save hash password")
	acs.Password = methods.HashForNewPassword(signupdt.Password)

	// passing account details to save in db and send mail for verification
	span.LogKV("task", "signup user")
	msg, ok := signup.Signup(acs, signupdt.Password, span)
	span.LogKV("task", "send final output")
	c.JSON(200, gin.H{"error": !(ok), "message": msg})
}

//==============================================================================

// Email defining structure for binding send code again data
type Email struct {
	Email string `json:"email" binding:"required"`
}

// SendCodeAgain is api handler for sending verification code again
func SendCodeAgain(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "send code again")
	// binding request body data
	span.LogKV("task", "binding body data")
	var email Email
	if err := c.BindJSON(&email); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Email is required."})
		return
	}

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(email.Email) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	// passing passed email to sendcodeagain function and in response boolean or message
	span.LogKV("task", "call SendCodeAgain function")
	msg, ok := signup.SendCodeAgain(strings.ToLower(email.Email), span)

	span.LogKV("task", "send final output")
	// checking boolean is true or false
	if !ok {
		// if false sending unable to send code again
		c.JSON(400, gin.H{"error": !(ok), "message": msg})
		return
	}

	c.JSON(200, gin.H{"error": !(ok), "message": msg})
}

//==============================================================================
