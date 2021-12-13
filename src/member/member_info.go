package member

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	config1 "git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	core "git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// SaveMemberInfo is a method to save member information of invited member through invite link token
func SaveMemberInfo(token, passwrd, password, wsID, name, contact string, parentspan opentracing.Span) (int, map[string]interface{}) {
	// start span from parent span context
	span := opentracing.StartSpan("save member info", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	mapd := make(map[string]interface{})

	// connect to database
	span.LogKV("task", "intialise db connection")
	db, err := gorm.Open("postgres", config.DBConfig())
	if err != nil {
		span.LogKV("task", "send final output after error in connecting to db")
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "Unable to connect to database."
		return 501, mapd
	}
	// close database client after every operations
	defer db.Close()

	//Checking token in database
	span.LogKV("task", "verify token")
	tok, err := verifyToken.CheckToken(token)
	if err != nil {
		span.LogKV("task", "send final output when token is invalid")
		log.Println(err)
		mapd["error"] = true
		mapd["message"] = "Invalid Join Link."
		return 404, mapd
	}
	//check token task
	if tok.TokenTask != "invite_link" {
		span.LogKV("task", "send final output when token task is invalid")
		mapd["error"] = true
		mapd["message"] = "Invalid Join Link."
		return 404, mapd
	}

	//fetch account informtion on basis of id
	span.LogKV("task", "fetch account details on basis of userid")
	acc := accounts.GetAccountForUserid(tok.Userid)

	//check workspace
	span.LogKV("task", "check user is member of workspace")
	if !CheckMember(acc.Email, wsID) {
		span.LogKV("task", "send final output when user is not member")
		mapd["error"] = true
		mapd["message"] = "You have not been invited in this workspace"
		return 404, mapd
	}
	if config1.Conf.Service.ISAWS == "true" {
		err = awsWork(acc, wsID)
		if err != nil {
			span.LogKV("task", "send final output when aws work not completed")
			mapd["error"] = true
			mapd["message"] = "Please try again later."
			return 404, mapd
		}
	}
	//update details
	span.LogKV("task", "update password in database and date in workspace member table")
	db.Model(&core.Accounts{}).Where("userid=?", tok.Userid).Updates(map[string]interface{}{"name": name, "contact_no": contact, "password": password, "creation_date": time.Now().Unix(), "account_status": "active", "verify_status": "verified"})
	db.Model(&database.WorkspaceMembers{}).Where("member_email=?", acc.Email).Update("joined", time.Now().Unix())
	acc.Name = name

	//generate jwt Token
	span.LogKV("task", "fetch jwt token")

	info := accounts.GetAccountForUserid(tok.Userid)
	//create account in ldap
	err = ldap.CreateLDAPAccount(info, passwrd)
	if err != nil {
		span.LogKV("task", "send final output when aws work not completed")
		mapd["error"] = true
		mapd["message"] = "Please try again later."
		mapd["error_message"] = err.Error()
		return 404, mapd
	}

	// delete used tokens
	go verifyToken.DeleteToken(token)

	mapd = jwtToken.JwtTokenusingWID(acc, wsID, "user")
	mapd["name"] = name
	mapd["role_id"] = acc.RoleID
	mapd["email"] = acc.Email
	return 200, mapd
}
