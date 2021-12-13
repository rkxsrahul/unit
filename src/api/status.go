package api

import (
	"log"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/workspace"
)

// WorkspaceStatus is an api handler to find status of workspace that worksapce deployed successfully
func WorkspaceStatus(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check status of workspace")
	//bind body request data
	span.LogKV("task", "binding body data")
	var data database.Workspaces
	if err := c.BindJSON(&data); err != nil {
		span.LogKV("task", "send final output when request data is wrong")
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "WorkSpace URL is required."})
		return
	}

	span.LogKV("task", "call function to find status of workspace")
	code, mapd := workspace.Status(data.WorkspaceID, span)
	span.LogKV("task", "send final output")
	c.JSON(code, mapd)
}
