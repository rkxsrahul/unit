package login

import (
	"log"
	// for gorm there is need to add a blank import for dialects

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// TokenLogin is a method to login with recover workpsace token
func TokenLogin(email, token, workspace string, parentspan opentracing.Span) (int, map[string]interface{}) {
	mapd := make(map[string]interface{})
	// start span from parent span contex
	span := opentracing.StartSpan("login with token function", opentracing.ChildOf(parentspan.Context()))
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
	//fetch account
	span.LogKV("task", "fetch user detail on basis of email id")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "send final output after error in fetching user detail")
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = err.Error()
		return 400, mapd
	}
	log.Println(acc)

	//Checking token in database
	span.LogKV("task", "check token in database")
	tok, err := verifyToken.CheckToken(token)
	if err != nil {
		span.LogKV("task", "send final output when token is invalid")
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = err.Error()
		return 404, mapd
	}
	log.Println(tok)
	//check token task
	if tok.TokenTask != "recover_workspace" {
		span.LogKV("task", "send final output when token task is invalid")
		log.Println("token task")
		mapd["error"] = true
		mapd["message"] = "Invalid or expired token"
		return 404, mapd
	}

	//compare userid
	if tok.Userid != acc.Userid {
		span.LogKV("task", "send final output when token user is invalid")
		log.Println("userid compare")
		mapd["error"] = true
		mapd["message"] = "Invalid or expired token."
		return 404, mapd
	}

	//fetch member details
	span.LogKV("task", "check user belongs to that workspace")
	var member []database.WorkspaceMembers
	db.Where("workspace_id=? AND member_email=?", workspace, acc.Email).Find(&member)
	if len(member) == 0 {
		span.LogKV("task", "send final output when user not belongs to that workspace")
		log.Println("workspace find")
		mapd["error"] = true
		mapd["message"] = "Invalid or expired token."
		return 404, mapd
	}
	// fetch workspace details
	span.LogKV("task", "fetch workspace details")
	var work []database.Workspaces
	db.Where("workspace_id = ?", workspace).Find(&work)
	log.Println(work)
	if len(work) == 0 {
		span.LogKV("task", "send final output when workspace not exists")
		mapd["error"] = true
		mapd["message"] = "Invalid or expired token."
		return 404, mapd
	}

	// delete used token
	go verifyToken.DeleteToken(token)

	span.LogKV("task", "generate jwt token")
	// generate jwt token
	mapd = jwtToken.JwtTokenusingWID(acc, workspace, member[0].Role)
	mapd["name"] = acc.Name
	mapd["workspace_role"] = member[0].Role
	mapd["workspace_name"] = work[0].TeamName
	mapd["role_id"] = acc.RoleID
	mapd["email"] = acc.Email
	span.LogKV("task", "send final output")
	return 200, mapd
}
