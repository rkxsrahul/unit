package accounts

import (
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// VerifyMail is method used to verify account
func VerifyMail(email, code string, parentspan opentracing.Span) (database.Accounts, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("verify mail function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return database.Accounts{}, false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// checking email exists or not
	span.LogKV("task", "fetch account details to check email is there or not")
	var acc []database.Accounts
	db.Where("email=?", email).Find(&acc)

	if len(acc) == 0 {
		span.LogKV("task", "send final output when no account is there")
		return database.Accounts{}, false
	}

	span.LogKV("task", "check token is valid")
	//Checking token in database on basis of token and userid
	var tok []database.Tokens
	db.Where("userid= ? AND token= ? AND token_task=?", acc[0].Userid, code, "email_verification").Find(&tok)

	//when token not found
	if len(tok) == 0 {
		span.LogKV("task", "send final when token not found")
		return database.Accounts{}, false
	}

	//check token is expired or not
	if (time.Now().Unix() - tok[0].Timestamp) > config.Conf.Service.VerifyLinkTimeout {
		span.LogKV("task", "send final output when token is invalid")
		return database.Accounts{}, false
	}

	span.LogKV("task", "update database when token is valid")
	//update account db and set verify satus to verified
	db.Model(&database.Accounts{}).Where("userid=?", tok[0].Userid).Update("verify_status", "verified")
	//if account status is 'new' then only change it to 'active'
	if acc[0].AccountStatus == "new" {
		db.Model(&database.Accounts{}).Where("userid=?", tok[0].Userid).Update("account_status", "active")
	}

	span.LogKV("task", "delete used and expired tokens")
	//deletion of expired tokens.
	db.Where("token_task=? AND timestamp<?", "email_verification", strconv.FormatInt((time.Now().Unix()-config.Conf.Service.VerifyLinkTimeout), 10)).Delete(&database.Tokens{})
	db.Where("token=?", code).Delete(&database.Tokens{})

	span.LogKV("task", "send final output")
	return acc[0], true
}
