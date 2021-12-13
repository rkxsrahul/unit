package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	jwtToken "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/vm"
)

func VMRequest(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("email not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	name, ok := claims["name"].(string)
	if !ok {
		log.Println("name not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	wsID, ok := claims["workspace"].(string)
	if !ok {
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	var data database.VMRequestInfo
	if err := c.BindJSON(&data); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": true, "message": "WorkSpace URL is required."})
		return
	}
	data.UserEmail = email
	data.UserName = name
	data.Workspace = wsID
	err := vm.VMRequest(data)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	log.Println(err)
	c.JSON(200, gin.H{"error": false, "message": "Your request has been submitted successfully."})
	return

}

func ServeFile(c *gin.Context) {

	token := c.Query("token")
	//fetching only token from whole string
	token = strings.TrimPrefix(token, "Bearer ")
	// parsing token and checking its validity
	rtoken, err := jwtToken.Parse(token, func(token *jwtToken.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtToken.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Conf.JWT.PrivateKey), nil
	})
	// if any err return nil claims
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	claims := rtoken.Claims.(jwtToken.MapClaims)

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("email not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}

	//fetch the file path
	path := accounts.Server(email)
	if path == "" {
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+email+"_vpn_config.opvn")
	http.ServeFile(c.Writer, c.Request, path)
	//	return

}

func VPNAccessKeys(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("email not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again"})
		return
	}

	data := accounts.VPNAccess(email)
	if data.Email == email {
		c.JSON(200, gin.H{"error": false, "username": data.Username, "password": data.Password})
		return
	}

	c.JSON(400, gin.H{"error": true, "message": "No Data Found"})
	return
}
