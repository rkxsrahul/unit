package workspace

import (
	"errors"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// Forgot is a method to send recover workspace mail to a valid account
func Forgot(email string, parentspan opentracing.Span) error {
	// start span from parent span context
	span := opentracing.StartSpan("forgot password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB

	// fetch account details on basis of email
	span.LogKV("task", "fetch account details on basis of email")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "any error in fetching details")
		log.Println(err)
		return err
	}
	log.Println(acc)

	if acc.VerifyStatus == "not_verified" {
		span.LogKV("task", "when account is not verfied")
		return errors.New("Please verify your account first")
	}

	//fetch workspace ids correspondance to that user
	span.LogKV("task", "fetch workspaces correspondance to user")
	members := []database.WorkspaceMembers{}
	db.Where("member_email=?", email).Find(&members)
	if len(members) < 1 {
		return errors.New("No workspace is assigned to you")
	}

	span.LogKV("task", "send mail")
	go mail.SendRecoveryMail(acc)
	return nil
}

// WorkSpaceEmail is a structure for sending list of workspaces recovered by token
type WorkSpaceEmail struct {
	database.Workspaces
	Email string
}

// RecoverWorkspace is a method to list all the workspaces related to that user on basis of recover workspace token
func RecoverWorkspace(token string, parentspan opentracing.Span) (int, map[string]interface{}) {
	mapd := make(map[string]interface{})
	// start span from parent span context
	span := opentracing.StartSpan("recover workspace method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	db := config.DB

	// close db instance whenever whole work completed
	defer db.Close()

	//Checking token in database
	span.LogKV("task", "verify token")
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
		mapd["error"] = true
		mapd["message"] = "Invalid or expired token."
		return 404, mapd
	}

	span.LogKV("task", "fetch workspace list on basis of userid")
	mapd["error"] = false
	mapd["workspaces"] = workspaceList(tok.Userid, span)
	span.LogKV("task", "send final output")
	return 200, mapd
}

// workspaceList is a method to fetch workspace correspondance to user on basis of userid
func workspaceList(userid string, parentspan opentracing.Span) []WorkSpaceEmail {
	// start span from parent span context
	span := opentracing.StartSpan("workspace list method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// array with workspace list
	workspaces := make([]WorkSpaceEmail, 0)
	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return workspaces
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// fetch account on basis of userid
	span.LogKV("task", "fetch account details on basis of userid")
	acc := accounts.GetAccountForUserid(userid)
	log.Println(acc)

	//fetch workspace ids correspondance to that user
	span.LogKV("task", "fetch workspaces correspondance to user")
	members := []database.WorkspaceMembers{}
	db.Where("member_email=?", acc.Email).Find(&members)

	// fetch workspace detail correspondance to each workspace id
	span.LogKV("task", "fetch each workspace detail")
	for i := 0; i < len(members); i++ {
		ws := []database.Workspaces{}
		db.Where("workspace_id = ?", members[i].WorkspaceID).Find(&ws)
		if len(ws) != 0 {
			workspaces = append(workspaces, WorkSpaceEmail{ws[0], acc.Email})
		}
	}
	span.LogKV("task", "send final list")
	return workspaces
}

//1621511267
