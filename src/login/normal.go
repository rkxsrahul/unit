package login

import (
	"strings"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	core "git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/signin"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// NormalLogin is a method to login with workpsace
func NormalLogin(email, password, wsID string, parentspan opentracing.Span) (int, map[string]interface{}) {
	mapd := make(map[string]interface{})
	// start span from parent span contex
	span := opentracing.StartSpan("signup function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	mapd["error"] = true
	// 	mapd["message"] = "Unable to connect to database"
	// 	return 501, mapd
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// check email and password
	span.LogKV("task", "call simplesignin function of auth-core to match email and password and check status of account")
	isAccLocked, login, msg, acc := signin.SimpleSignin(strings.ToLower(email), password, span)

	// when account is locked or email or password is wrong
	if isAccLocked || !login {
		span.LogKV("task", "send final output when login not succesfully")
		mapd["error"] = true
		mapd["message"] = msg
		return 401, mapd
	}

	status, msg, loggedinUID := ldap.Authentication(strings.ToLower(email), password)
	if status && loggedinUID != "" {

		// if workspace is passed
		if wsID != "" {
			//check user is member of that workspace
			span.LogKV("task", "check member belongs to that workspace")
			member := []database.WorkspaceMembers{}
			db.Where("member_email= ? AND workspace_id= ?", acc.Email, wsID).Find(&member)
			if len(member) == 0 {
				//when not a member
				span.LogKV("task", "when invalide workspace is passed")
				mapd["error"] = true
				mapd["message"] = "invalid email and password"
				return 401, mapd
			}

			var work []database.Workspaces
			db.Where("workspace_id = ?", wsID).Find(&work)
			if len(work) == 0 {
				//when not a member
				span.LogKV("task", "when invalid workspace is passed")
				mapd["error"] = true
				mapd["message"] = "invalid email and password"
				return 501, mapd
			}

			//generate jwt token
			span.LogKV("task", "generate jwt token")
			mapd = jwtToken.JwtTokenusingWID(acc, wsID, member[0].Role)
			mapd["name"] = acc.Name
			mapd["workspace_role"] = member[0].Role
			mapd["workspace_name"] = work[0].TeamName
			mapd["role_id"] = acc.RoleID
			mapd["email"] = acc.Email

		} else {
			// when workspace is empty
			// check user belongs to any WorkSpace
			span.LogKV("task", "check user belongs to any workspace")
			isWs, token := checkUser(acc)
			//generate jwt token
			span.LogKV("task", "generate jwt token")
			mapd = jwtToken.JwtTokenusingWID(acc, "", "")
			mapd["name"] = acc.Name
			mapd["role_id"] = acc.RoleID
			mapd["email"] = acc.Email
			mapd["isWorkspace"] = isWs
			mapd["workspace_token"] = token
		}
	} else {
		if strings.Contains(msg, "LDAP Result Code 32") {
			msg = "Internal Server Error"
			mapd["error"] = true
			mapd["message"] = msg
			return 500, mapd
		}
		mapd["error"] = true
		mapd["message"] = msg
		return 401, mapd
	}

	span.LogKV("task", "send final output with status 200")
	return 200, mapd
}

// checkUser is a method to check user belongs to any WorkSpace
// if yes follow recover workspace process
func checkUser(acc core.Accounts) (bool, string) {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return false, ""
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// cheking userid in workspace member table
	var wm []database.WorkspaceMembers
	db.Where("member_email= ?", acc.Email).Find(&wm)

	var is bool
	// if there is no user in this table then return false
	if len(wm) == 0 {
		return false, ""
	}
	is = false
	// when there is user in worksapce member check role of user for that workspace
	for i := 0; i < len(wm); i++ {
		if wm[i].Role == "owner" {
			is = true
		}
	}
	if !is {
		return is, ""
	}
	// getting token to recover workspace used to directly login in workspace.
	token := verifyToken.CheckSentToken(acc.Userid, "recover_workspace")

	return is, token
}
