package forgotpass

import (
	"errors"
	"log"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/member"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

//==============================================================================

// ForgotPassChallenge is a method for sending reset-password link in mail with workspace
func ForgotPassChallenge(email, wsID string, parentspan opentracing.Span) (string, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("forgot password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	// span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return "Unable to connect to database", false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()

	// // when workspace id is not send
	// if wsID == "" {
	// 	msg, ok := forgotPassChallenge(email, span)
	// 	return msg, ok
	// }

	//checking whether email exist or not
	span.LogKV("task", "fetch account details to check email is there or not")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "send final output when no account is there")
		log.Println("when no account is there")
		return "Email doesn't exists.", false
	}

	//split workspace name
	wsParts := strings.Split(wsID, ".")

	if acc.VerifyStatus == "verified" && acc.AccountStatus == "active" {
		// check member exist in workspace
		span.LogKV("task", "check user is member of that workspace")
		if member.CheckMember(acc.Email, wsParts[0]) {
			//send forgot password mail
			go mail.SendForgotPassMail(acc, wsID)
			span.LogKV("task", "send final output ")
			return "We have sent a password reset link on your mail.", true
		}
		log.Println("when workspace not match")
	}
	span.LogKV("task", "send final output when no account is there")
	return "Email doesn't exists.", false
}

// ==============================================================================

// ResetForgottenPass is method to reset password in database
// first check request is with workspace or without
// but before updating token and password is checked
func ResetForgottenPass(token, password, wsID string, parentspan opentracing.Span) (string, string, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("reset password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// if wsID is not present
	if wsID == "" {
		span.LogKV("task", "reset password without workspace")
		email, msg, ok := resetForgottenPass(token, password, span)
		span.LogKV("task", "send final output after resetting password")
		return email, msg, ok
	}

	//validation check on password
	span.LogKV("task", "check password validation")
	if !methods.CheckPassword(password) {
		span.LogKV("task", "send final output when password is wrong")
		return "", "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character.", false
	}

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return "", "Unable to connect to database.", false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()

	//Checking token in database
	span.LogKV("task", "verify token")
	tok, err := verifyToken.CheckToken(token)
	if err != nil {
		span.LogKV("task", "send final output when token is invalid")
		return "", err.Error(), false
	}

	//split workspace name
	wsParts := strings.Split(wsID, ".")

	if tok.TokenTask == "forgot_pass" {
		//fetch account informtion on basis of id
		span.LogKV("task", "fetch account details on basis of userid")
		acc := accounts.GetAccountForUserid(tok.Userid)
		// check member exist in workspace
		span.LogKV("task", "check user is member of that workspace")
		if member.CheckMember(acc.Email, wsParts[0]) {
			// update in database
			span.LogKV("task", "update password in database")
			email, _ := updateDatabase(tok.Userid, password)

			// delete used and expired tokens
			go verifyToken.DeleteToken(token)

			//password reset done.
			span.LogKV("task", "Password reset successfully")
			return email, "Password reset successfully.", true
		}
	}
	span.LogKV("task", "send final output when token is invalid")
	return "", "Invalid or expired token.", false
}

//==============================================================================

// ForgotPassChallenge is a method for sending reset-password link in mail
func forgotPassChallenge(email string, parentspan opentracing.Span) (string, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("forgot password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return "Unable to connect to database", false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	//checking whether email exist or not
	span.LogKV("task", "fetch account details to check email is there or not")
	var acc []database.Accounts
	db.Where("email=?", email).Find(&acc)
	if len(acc) != 0 {
		if acc[0].VerifyStatus == "verified" && acc[0].AccountStatus == "active" {
			span.LogKV("task", "send reset password link in mail")
			//send verification mail with task as forgot_pass
			go mail.SendForgotPassMailAccount(acc[0])
			span.LogKV("task", "send final output ")
			return "We have sent a password reset link on your mail.", true
		}
	}
	span.LogKV("task", "send final output when no account is there")
	return "Email doesn't exists.", false
}

// ==============================================================================

// ResetForgottenPass is method to reset password in database
// but before updating token and password is checked
func resetForgottenPass(token, password string, parentspan opentracing.Span) (string, string, bool) {
	// start span from parent span context
	span := opentracing.StartSpan("reset password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	//validation check on password
	span.LogKV("task", "check password validation")
	if !methods.CheckPassword(password) {
		span.LogKV("task", "send final output when password is wrong")
		return "", "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character.", false
	}

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return "", "Unable to connect to database.", false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()

	//Checking token in database
	span.LogKV("task", "verify token")
	tok, err := verifyToken.CheckToken(token)
	if err != nil {
		span.LogKV("task", "send final output when token is invalid")
		return "", err.Error(), false
	}

	if tok.TokenTask == "forgot_pass" {
		// update in database
		span.LogKV("task", "update database")
		email, err := updateDatabase(tok.Userid, password)
		if err != nil {
			span.LogKV("task", "send final output after error in updating database")
			return "", err.Error(), false
		}
		// delete used and expired tokens
		go verifyToken.DeleteToken(token)

		//password reset done.
		span.LogKV("task", "Password reset successfully")
		return email, "Password reset successfully.", true
	}
	span.LogKV("task", "send final output when token task is invalid")
	return "", "Invalid or expired token.", false
}

//==============================================================================

// UpdateDatabase is method to update password in database on basis of userid
func updateDatabase(userid, password string) (string, error) {

	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return "", errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// check account exist
	var acs []database.Accounts
	db.Where("userid= ?", userid).Find(&acs)
	if len(acs) == 0 {
		return "", errors.New("Account not found")
	}
	// hash the simple password and then update password in database
	db.Model(&database.Accounts{}).Where("userid=?", userid).Update("password", methods.HashForNewPassword(password))

	return acs[0].Email, nil
}
