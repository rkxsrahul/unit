package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	ot "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/api"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ginjwt"
)

// V1Routes is a method in which all the service endpoints are defined
func V1Routes(router *gin.Engine) {
	// health endpoint
	router.GET("/healthz", opengintracing.NewSpan(ot.GlobalTracer(), "checking health of service"), api.Healthz)
	// developer help endpoint
	if config.Conf.Service.Environment != "production" {
		// endpoint to read logs
		router.GET("/logs", checkToken, readLogs)
		// endpoint to read variables
		router.GET("/end", checkToken, readEnv)
	}

	// intialize v1 group
	v1 := router.Group("/v1")

	// signup routes
	// account creation endpoint
	v1.POST("/signup", opengintracing.NewSpan(ot.GlobalTracer(), "signup user"), api.SignupEndpoint)
	// verify mail id on basis of token
	v1.POST("/verifymail", opengintracing.NewSpan(ot.GlobalTracer(), "verify user mail"), api.VerifyMailEp)
	// if verify code get expired used to send code again at mail id
	v1.POST("/send_code_again", opengintracing.NewSpan(ot.GlobalTracer(), "send code for verification"), api.SendCodeAgain)

	// workspace specific routes
	// to check workspace status
	v1.POST("/workspace_status", opengintracing.NewSpan(ot.GlobalTracer(), "workspace status"), api.WorkspaceStatus)
	// login in workspace
	v1.POST("/workspace_login", opengintracing.NewSpan(ot.GlobalTracer(), "workspace login"), api.WorkSpaceLoginEp)
	// forgot worksapce
	v1.POST("/forgot_workspace", opengintracing.NewSpan(ot.GlobalTracer(), "forgot workspace"), api.ForgotWorkspaceEp)
	// get workspace list of a user
	v1.POST("/get_workspace_list", opengintracing.NewSpan(ot.GlobalTracer(), "list workspace on basis of token"), api.GetWorkspaceList)
	// signup member in a workspace
	v1.POST("/membersignup", opengintracing.NewSpan(ot.GlobalTracer(), "member signup"), api.MemberSignupEp)
	// this function sets new password according to invited link or token
	v1.POST("/memberregistration", opengintracing.NewSpan(ot.GlobalTracer(), "save invite member data"), api.MemberRegistration)

	// login routes
	// login with a token
	v1.POST("/tokenlogin", opengintracing.NewSpan(ot.GlobalTracer(), "login with token"), api.TokenLogin)
	// login with workspace
	v1.POST("/login", opengintracing.NewSpan(ot.GlobalTracer(), "login in account with worspace"), api.LoginEndpoint)

	// forgot password routes
	// used to get link when user forgot password and also for reset password
	v1.POST("/forgotpass", opengintracing.NewSpan(ot.GlobalTracer(), "forgot and reset password"), api.ForgotPassEp)

	v1.GET("/vpn/guide", checkJWT, serveFile)

	v1.GET("/vpnconnection/file.ovpn", api.ServeFile)

	//setting up middleware for protected apis
	authMiddleware := ginjwt.MwInitializer()

	//Protected resources
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		// adding custom middleware for checking token validity
		v1.Use(opengintracing.NewSpan(ot.GlobalTracer(), "check jwt token validity"), api.CheckTokenValidity)
		{
			// session apis
			v1.GET("/refresh_token", opengintracing.NewSpan(ot.GlobalTracer(), "refresh jwt token"), api.RefreshToken)
			v1.GET("/check_token", opengintracing.NewSpan(ot.GlobalTracer(), "check jwt token"), api.CheckToken)
			v1.GET("/logout", opengintracing.NewSpan(ot.GlobalTracer(), "logout session"), api.Logout)

			// user apis
			user := v1.Group("/")

			user.POST("/vm_request", api.VMRequest)

			//middleware for user
			user.Use(opengintracing.NewSpan(ot.GlobalTracer(), "check customer is user"), api.CheckUser)
			{
				// profile related routes
				// api for changing password
				user.PUT("/changepass", opengintracing.NewSpan(ot.GlobalTracer(), "change account password"), api.ChangePasswordEp)
				// api for view profile
				user.GET("/profile", opengintracing.NewSpan(ot.GlobalTracer(), "view user profile"), api.ViewProfile)
				// api for view profile
				user.PUT("/profile", opengintracing.NewSpan(ot.GlobalTracer(), "update user profile"), api.UpdateProfile)

				// user workspace routes
				// api for creating workspace
				user.POST("/workspaces", opengintracing.NewSpan(ot.GlobalTracer(), "create workspace"), api.CreateWorkSpaceEp)
				// api for checking workspace Availability
				user.POST("/workspace_availability", opengintracing.NewSpan(ot.GlobalTracer(), "check workspace availability"), api.WorkspaceAvailability)
				// send aws access key to mail
				user.GET("/accessKeys", api.MailAccessKeys)

				//
				v1.GET("/vpn/access", api.VPNAccessKeys)
				owner := user.Group("/")
				owner.Use(opengintracing.NewSpan(ot.GlobalTracer(), "check customer is user"), api.CheckOwner)
				{
					// send invite to member for joining workspace
					owner.POST("/workspaces/member", opengintracing.NewSpan(ot.GlobalTracer(), "invite member"), api.InviteWorkspaceMembers)
					// fetch members of a WorkSpace
					owner.GET("/workspaces/member", opengintracing.NewSpan(ot.GlobalTracer(), "list members"), api.WorkspaceMembers)
					owner.DELETE("/workspaces/member", opengintracing.NewSpan(ot.GlobalTracer(), "delete member"), api.DeleteWorkspaceMember)
				}
				// chech member belongs to workspace and joined the workspace
				user.GET("/workspaces/checkmember", opengintracing.NewSpan(ot.GlobalTracer(), "check member"), api.WorkspaceCheckMember)
			}
			// admin apis
			admin := v1.Group("/admin")
			//middleware for user
			admin.Use(opengintracing.NewSpan(ot.GlobalTracer(), "check admin"), api.CheckAdmin)
			{
				// admin workspace routes
				// list workspaces
				admin.GET("/workspaces", opengintracing.NewSpan(ot.GlobalTracer(), "list workspaces"), api.ListWorkspaces)
				// delete workspace
				admin.DELETE("/workspaces/:workspace_id", opengintracing.NewSpan(ot.GlobalTracer(), "delete worksapce"), api.DeleteWorkspace)
				// delete user account
				admin.DELETE("/account_emails/:email", opengintracing.NewSpan(ot.GlobalTracer(), "delete account"), api.DeleteAccountByEmail)
				// tester help endpoint
				if config.Conf.Service.Environment != "production" {
					// toggle mail service
					admin.PUT("/mail/:value", api.ChangeMail)
					// toggle otp service
					admin.PUT("/otp/:value", api.ChangeOTP)
				}
			}
		}
	}
}

// readLogs is a api handler for reading logs
func readLogs(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "info.txt")
}

// readEnv is api handler for reading configuration variables data
func readEnv(c *gin.Context) {
	if config.TomlFile == "" {
		// if configuration is done using environment variables
		env := make([]string, 0)
		for _, pair := range os.Environ() {
			env = append(env, pair)
		}
		c.JSON(200, gin.H{
			"environments": env,
		})
	} else {
		// if configuration is done using toml file
		http.ServeFile(c.Writer, c.Request, config.TomlFile)
	}
}

// checkToken is a middleware to check header is set or not for secured api
func checkToken(c *gin.Context) {
	xt := c.Request.Header.Get("STACKLABS-TOKEN")
	if xt != "slAuth1010" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

func checkJWT(c *gin.Context) {
	token := c.Query("token")
	//fetching only token from whole string
	token = strings.TrimPrefix(token, "Bearer ")
	// parsing token and checking its validity
	_, err := jwtgo.Parse(token, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Conf.JWT.PrivateKey), nil
	})
	// if any err return nil claims
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	c.Next()
}
func serveFile(c *gin.Context) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=guide.pdf")
	http.ServeFile(c.Writer, c.Request, "./guide.pdf")
}
