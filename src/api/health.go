package api

import (
	"log"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/health"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	logs "github.com/opentracing/opentracing-go/log"
)

// Healthz is an api handler to check health of service
func Healthz(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("span not found")
		c.AbortWithStatus(500)
		return
	}
	defer span.Finish()
	span.SetTag("event", "check health")
	// call health service check function
	span.LogKV("task", "start health check")
	err := health.ServiceHealth(span)
	if err != nil {
		span.LogFields(
			logs.String("task", "stop health check"),
			logs.String("error", err.Error()),
		)
		// if any error is there
		log.Println(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogFields(
		logs.String("info", "stop health check"),
		logs.String("error", "nil"),
	)
	// if no error is there
	c.JSON(200, gin.H{
		"error":       false,
		"message":     "All is okay",
		"build":       config.Conf.Service.Build,
		"environment": config.Conf.Service.Environment,
	})
}
