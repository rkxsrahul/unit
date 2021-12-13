package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/admin"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
	serviceWork "git.xenonstack.com/stacklabs/stacklabs-auth/src/workspace"
)

// ListWorkspaces is an api handler for listing deployed workspaces
func ListWorkspaces(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "list workspaces")
	// fetch workspace from database
	span.LogKV("task", "call list workspaces function")
	list, err := admin.ListWorkSpaces(span)
	if err != nil {
		span.LogKV("task", "when there is some error in listing workspace")
		c.JSON(501, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "send final list")
	c.JSON(200, gin.H{
		"error":      false,
		"workspaces": list,
	})
}

//=============================================================//

// DeleteWorkspace is an api handler for deleting workspace by admin from db and helm or kube
func DeleteWorkspace(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "delete workspace")
	span.LogKV("task", "fetch token from request headers")
	token := c.Request.Header.Get("Authorization")
	// delete workspace data from db and cluster
	span.LogKV("task", "call delete workspace function")
	err := serviceWork.DeleteWorkspace(c.Param("workspace_id"), token, span)
	if err != nil {
		span.LogKV("task", "error in deleting")
		c.JSON(501, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "send final output")
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Request Accepted.",
	})
}

//============================================================//

// DeleteAccountByEmail is an api handler for deleting user account on basis of email from db
func DeleteAccountByEmail(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "delete user account")

	span.LogKV("task", "check email validation")
	if !methods.ValidateEmail(c.Param("email")) {
		span.LogKV("task", "send final output when email is wrong")
		// if there is some error passing bad status code
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Please enter valid email id."})
		return
	}

	//delete user details from database
	span.LogKV("task", "delete user on basis of email")
	err := admin.DeleteAccount(strings.ToLower(c.Param("email")), span)
	if err != nil {
		span.LogKV("task", "error in deleting")
		c.JSON(400, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err = ldap.DeleteLDAPAccount(strings.ToLower(c.Param("email")))
	if err != nil {
		span.LogKV("task", "error in deleting")
		c.JSON(400, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	span.LogKV("task", "send final output")
	c.JSON(200, gin.H{
		"error":   false,
		"message": "account delete succesfully",
	})
}
