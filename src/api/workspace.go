package api

import (
	"log"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/workspace"
	serviceWork "git.xenonstack.com/stacklabs/stacklabs-auth/src/workspace"
)

// CreateWorkSpaceEp is an api handler for creating workspace
func CreateWorkSpaceEp(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "create worksapce")
	// extract jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)

	//bind body request data
	span.LogKV("task", "binding body data")
	var data database.Workspaces
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "WorkSpace URL is required."})
		return
	}

	//fetch id from claims
	email, ok := claims["email"]
	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	// create workspace
	span.LogKV("task", "call function for sending request to deployment service to deploy workspace")
	code, mapd := serviceWork.CreateWorkspace(email.(string), c.GetHeader("Authorization"), data, span)
	span.LogKV("task", "workspace creation done")
	c.JSON(code, mapd)
}

//=========================================================================//

// WorkspaceAvailability is an api handler for checking workspace Availability
func WorkspaceAvailability(c *gin.Context) { // fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check workspace availability")
	//bind body request data
	span.LogKV("task", "binding body data")
	var data database.Workspaces
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "WorkSpace URL is required."})
		return
	}
	span.LogKV("task", "call function checking workspace exist in database")
	code, message := workspace.CheckWorkspaceAvailability(data, span)
	span.LogKV("task", "send final output")
	c.JSON(code, gin.H{
		"error":   code != 200,
		"message": message,
	})
}

//=========================================================================//

// WorkSpaceLoginEp is an api handler for login in a workspace
// mainly checks workspace exist or not
func WorkSpaceLoginEp(c *gin.Context) {
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "login in workspace")
	//bind body request data
	span.LogKV("task", "binding body data")
	var data database.Workspaces
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "WorkSpace URL is required."})
		return
	}

	span.LogKV("task", "call function checking workspace exist in database")
	code, mapd := workspace.Login(data.WorkspaceID, span)
	span.LogKV("task", "send final output")
	c.JSON(code, mapd)
}
