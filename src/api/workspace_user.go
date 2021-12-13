package api

import (
	"log"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/member"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
)

// WorkSpaceMembers is a array of structure for binding data from body during invite member request
type WorkSpaceMembers []struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

// InviteWorkspaceMembers is an api handler to send mail to members who are invited
func InviteWorkspaceMembers(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "invite members in a workspace")

	// extract jwt claims
	span.LogKV("task", " extract jwt claims")
	claims := jwt.ExtractClaims(c)

	// fetching data from request body
	span.LogKV("task", "binding body data")
	var members WorkSpaceMembers
	if err := c.BindJSON(&members); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "Bad Request."})
		return
	}

	// save emails in string array
	span.LogKV("task", "fetch workspace, email and name from payload")
	wsID, ok := claims["workspace"]
	if !ok {
		span.LogKV("task", "send final output when workspace is not set in payload")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	ownerEmail, ok := claims["email"]
	if !ok {
		span.LogKV("task", "send final output when owner email is not set in payload")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	ownerName, ok := claims["name"]
	if !ok {
		span.LogKV("task", "send final output when owner name is not set in payload")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}

	// save emails in string array
	span.LogKV("task", "save email in string array")
	emails := make([]member.WorkSpaceMember, 0)
	mapd := make(map[string]interface{}, 0)
	for i := 0; i < len(members); i++ {
		_, ok = mapd[members[i].Email]
		if !ok {
			emails = append(emails, member.WorkSpaceMember{
				Email: members[i].Email,
				Role:  members[i].Role,
			})
			mapd[members[i].Email] = true
		}
	}
	// empty this map
	for k := range mapd {
		delete(mapd, k)
	}

	//send invite or login link to users
	span.LogKV("task", "call invite function to send mail to all valid users")
	code, mapd := member.Invite(emails, ownerEmail.(string), ownerName.(string), wsID.(string), span)
	span.LogKV("task", "send final output")
	c.JSON(code, mapd)
}

//===============================================================================

// WorkspaceMembers is api handler to list members of a workspace
func WorkspaceMembers(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "invite members in a workspace")

	// fetch workspace id from jwt claims
	span.LogKV("task", " extract jwt claims")
	claims := jwt.ExtractClaims(c)

	workspace, ok := claims["workspace"]
	if !ok {
		span.LogKV("task", "send final output when workspace is not set in payload")
		c.JSON(501, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}
	// fetch workspace members from database
	span.LogKV("task", " call memberList function to fetch all members from database")
	members, err := member.List(workspace.(string), span)
	if err != nil {
		span.LogKV("task", "error in fetching list")
		c.JSON(501, gin.H{"error": true, "message": "Unable to read data from db"})
		return
	}
	span.LogKV("task", " send list")
	c.JSON(200, gin.H{
		"error":   false,
		"members": members,
	})

}

//===============================================================================

// WorkspaceCheckMember is api handler to check member joined the workspace
func WorkspaceCheckMember(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check member in a workspace")

	// fetch workspace id from jwt claims
	span.LogKV("task", " extract jwt claims")
	claims := jwt.ExtractClaims(c)
	workspace, ok := claims["workspace"]
	if !ok {
		span.LogKV("task", "send final output when workspace is not set in payload")
		c.JSON(501, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}
	// check member in workspace from database
	span.LogKV("task", " call check function to check member joined the workspace")
	err := member.Check(workspace.(string), c.Query("email"), span)
	span.LogKV("task", " send final output")
	if err != nil {
		c.JSON(501, gin.H{"error": true})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
	})
}

//======================================================================================

// DeleteWorkspaceMember is an API handler for deleting member from a workspace
func DeleteWorkspaceMember(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "delete member in a workspace")

	// fetch workspace id from jwt claims
	span.LogKV("task", " extract jwt claims")
	claims := jwt.ExtractClaims(c)
	workspace, ok := claims["workspace"]
	if !ok {
		span.LogKV("task", "send final output when workspace is not set in payload")
		c.JSON(501, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	// check member in workspace from database
	span.LogKV("task", " call delete function to delete member from the workspace")
	err := member.Delete(workspace.(string), c.Query("email"), span)
	span.LogKV("task", " send final output")
	if err != nil {
		c.JSON(501, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Member deleted Successfully.",
	})
}
