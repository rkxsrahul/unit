package signup

import (
	"log"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

// initNewAccount is a method for intializing account
func initNewAccount(account *database.Accounts) {
	// saving new userid
	id := GetNewUserid()
	if id == "" {
		// when id is not generated
		return
	}
	account.Userid = id
	// if verify status is not_verified then account_status = new
	if account.VerifyStatus == "not_verified" {
		account.AccountStatus = "new"
	} else {
		account.AccountStatus = "active"
	}
	//default role will be user
	if account.RoleID == "" {
		account.RoleID = "user"
	}
}

//==============================================================================

// GetNewUserid is method for generating new id
func GetNewUserid() string {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return ""
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	var flags []database.Flags
	db.Where("name= ?", "new_userid").Find(&flags)
	log.Println(flags)
	if len(flags) != 0 {
		// converting value to int
		newUseridInt, _ := strconv.Atoi(flags[0].Value)
		// save value as string but after incrementing
		db.Exec("update flags set value= '" + strconv.Itoa(newUseridInt+1) + "' where name= 'new_userid';")
		return flags[0].Value
	}
	return ""
}

//==============================================================================

// SendCodeAgain is a method for sending code again to email for verification
func SendCodeAgain(email string, parentspan opentracing.Span) (string, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("verify mail function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return err.Error(), false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	db = db.Debug()

	//Checking for email whether exists or not
	span.LogKV("task", "fetch account details to check email is there or not")
	var account []database.Accounts
	db.Where("email=?", email).Find(&account)

	// if there is account and account status is new
	if len(account) != 0 {
		if account[0].VerifyStatus != "verified" {
			span.LogKV("task", "send mail for verification")
			// send mail again
			go mail.SendVerifyMail(account[0])
			return "Verification code sent.", true
		}
		span.LogKV("task", "Send final output email doesn't exists")
		return "Your account is already verified", false
	}
	span.LogKV("task", "Send final output email doesn't exists")
	return "Email doesn't exists.", false
}
