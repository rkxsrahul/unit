package api

import (
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/redisdb"
)

// CheckToken is API for checking token is valid
func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{})
}

// RefreshToken is api handler to generate new jwt token and expire old token
func RefreshToken(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "refresh token")
	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)

	// fetch account on basis of userid
	span.LogKV("task", "fetch account details on basis of userid")
	acc := accounts.GetAccountForUserid(claims["id"].(string))
	// generating new token
	span.LogKV("task", "generate new jwt token")
	mapd := jwtToken.JwtRefreshToken(claims, span)
	mapd["name"] = acc.Name
	mapd["email"] = acc.Email
	mapd["workspace_role"] = claims["role"].(string)
	mapd["sys_role"] = acc.RoleID

	// if any error in genrating new token
	if mapd["token"] == "" {
		span.LogKV("task", "final output when token is nil")
		c.JSON(501, gin.H{
			"error":   true,
			"message": "Error in generating new token",
		})
		return
	}
	if config.Conf.Service.IsLogoutOthers != "true" {
		// when succesfully generated new token
		// delete old token from redis
		// fetch token from header
		span.LogKV("task", "fetch old token from headers")
		token := c.Request.Header.Get("Authorization")
		// trim bearer from token
		token = strings.TrimPrefix(token, "Bearer ")
		// call delete token go function
		span.LogKV("task", "delete token from redis")
		err := redisdb.DeleteToken(token)
		if err != nil {
			span.LogKV("task", "final output when error in deleting token")
			c.JSON(501, gin.H{
				"error":   err,
				"message": "Error in deleting old token",
			})
			return
		}
	}
	span.LogKV("task", "send final output")
	c.JSON(200, mapd)
}

// CheckTokenValidity is a middleware for checking token validity using redis database
func CheckTokenValidity(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check token validitiy")

	span.LogKV("task", "fetch token from headers")
	// fetch token from header
	token := c.Request.Header.Get("Authorization")
	// trim bearer from token
	token = strings.TrimPrefix(token, "Bearer ")

	// check token exist or not
	span.LogKV("task", "check token exist or not")
	err := redisdb.CheckToken(token)
	if err != nil {
		// when token not exist
		span.LogKV("task", "final output when token not exist")
		c.Abort()
		c.JSON(401, gin.H{"error": true, "message": "Expired auth token"})
		return
	}
	span.LogKV("task", "final output when token exist")
	c.Next()
}

// CheckAdmin is a middleware for checking user is admin or not
func CheckAdmin(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check admin")

	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)
	// checking sys role
	if claims["sys_role"].(string) != "admin" {
		span.LogKV("task", "final output when user is not admin")
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "You are not authorized",
		})
		return
	}
	span.LogKV("task", "final output when user is admin")
	c.Next()
}

// CheckUser is a middleware for checking user is user or not
func CheckUser(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check user")

	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)
	// checking sys role
	if claims["sys_role"].(string) != "user" {
		span.LogKV("task", "final output when user is not user")
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "You are not authorized",
		})
		return
	}
	span.LogKV("task", "final output when user is user")
	c.Next()
}

// CheckOwner is a middleware for checking workspace user is owner or not
func CheckOwner(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		log.Println("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check user")

	// extracting jwt claims
	span.LogKV("task", "extract jwt claims")
	claims := jwt.ExtractClaims(c)

	// check user is owner of workpsace or user
	role, ok := claims["role"].(string)
	if !ok {
		span.LogKV("task", "final output when wrong token used")
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "You are not authorized",
		})
		return
	}
	if role != "owner" {
		span.LogKV("task", "final output when workspace user is not owner")
		c.Abort()
		c.JSON(403, gin.H{
			"error":   true,
			"message": "You are not authorized",
		})
		return
	}

	span.LogKV("task", "final output when workspace user is owner")
	c.Next()
}
