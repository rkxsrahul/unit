package api

import (
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/aws"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// GetUserProfile is a api handler to fetch users details
func GetUserProfile(c *gin.Context) {
	acc, err := accounts.GetAccountForEmail(strings.ToLower(c.Param("email")))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"account": acc,
	})
}

// ViewProfile is a api handler for viewing user profile
func ViewProfile(c *gin.Context) {
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "view profile data")
	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)

	// fetch profile on basis of email
	span.LogKV("task", "call function for fetching profile from database")
	acc, err := accounts.GetAccountForEmail(claims["email"].(string))
	if err != nil {
		span.LogKV("task", "send final output when there is some error")
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "Profile fetched Succesfully")
	c.JSON(200, gin.H{
		"error":   false,
		"account": acc,
	})
}

// UpdateData is a structure for binding update profile data
type UpdateData struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

// UpdateProfile is a api handler for updating user profile
func UpdateProfile(c *gin.Context) { // fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "update profile data")
	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)

	email, ok := claims["email"]

	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	// fetching data from request body
	span.LogKV("task", "binding body data")
	var data UpdateData
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass valid name and contact number",
		})
		return
	}

	//update name and contact of user
	span.LogKV("task", "call function for updating name and contact in database")
	err := accounts.UpdateProfile(email.(string), data.Name, data.Contact, span)
	if err != nil {
		span.LogKV("task", "send final output when there is some error")
		log.Println(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "Profile Updated Successfully")
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Profile Updated Successfully",
	})
}

func MailAccessKeys(c *gin.Context) {
	// extracting jwt claims    
	claims := jwt.ExtractClaims(c)

	if config.Conf.Service.ISAWS == "true" {
		// connecting to db
		db := config.DB

		// intialize variable with type accounts
		var acs []database.Accounts
		// fetching data on basis of userid
		db.Where("userid= ?", claims["id"].(string)).Find(&acs)
		if len(acs) == 0 {
			// if there is no account
			c.JSON(400, gin.H{"error": true, "message": "Account not found."})
			return
		}

		value := c.Request.Header.Get("Request_Type")

		//call aws function to create iam access keys and send it in a mail
		code, msg := aws.MailAccessKeys(acs[0], value)

		if value == "fetch" {
			// intialize variable with type access keys
			var keys []database.AccessKeys
			// fetching data on basis of userid
			db.Where("userid= ?", claims["id"].(string)).Find(&keys)
			if len(keys) == 0 {
				// if there is no account
				c.JSON(400, gin.H{"error": true, "message": "Account not found."})
				return
			}
			c.JSON(200, gin.H{
				"key":    keys[0].Key,
				"secret": keys[0].Secret,
			})
			return
		}

		c.JSON(code, gin.H{"message": msg})
		return
	} else {
		c.JSON(200, gin.H{
			"key":    "",
			"secret": "",
		})
		return
	}

}
