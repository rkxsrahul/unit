package signin

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// SimpleSignin is a method for checking email is there
// then password is matching with password saved in db
// correspondance to that email and then check account status
func SimpleSignin(email, password string, parentspan opentracing.Span) (bool, bool, string, database.Accounts) {
	// start span from parent span contex
	span := opentracing.StartSpan("signin function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return false, false, "Unable to connect to database.", database.Accounts{}
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB

	span.LogKV("task", "fetch account details to check account is there or not on basis of email")
	//Checking whether registered or not
	var account []database.Accounts
	db.Where("email=?", email).Find(&account)

	// when no account found
	if len(account) == 0 {
		span.LogKV("task", "send final output when no account is there")
		return false, false, "Invalid email or password.", database.Accounts{}
	}

	// checking previous failed logins
	span.LogKV("task", "check previous failed logins")
	// msg, isAccLocked, count := checkPreviousFailedLogins(account[0])
	// if isAccLocked {
	// 	// when account is locked
	// 	span.LogKV("task", "send final output when account is locked")
	// 	return true, false, msg, account[0]
	// }

	// checking password with saved password
	span.LogKV("task", "check password")
	if methods.CheckHashForPassword(account[0].Password, password) {
		// when password matched
		// checking account status
		span.LogKV("task", "check account status")
		switch account[0].AccountStatus {
		case "active":
			// when user is active all well
			span.LogKV("task", "send final output when status is active")
			return false, true, "", account[0]
		case "blocked":
			// when user is blocked
			span.LogKV("task", "send final output when account is block")
			return false, false, "Your account has been blocked.", database.Accounts{}
		case "new":
			// when user is new not verified
			span.LogKV("task", "send final output when account is new")
			return false, false, "Please verify your email.", database.Accounts{}
		}
	}
	// when password not matched
	span.LogKV("task", "send final output")
	return false, false, "Invalid email or password. You have  login attempts left", database.Accounts{}
}

//==============================================================================

// checkPreviousFailedLogins is a method for checking previous failed login of a user
func checkPreviousFailedLogins(account database.Accounts) (string, bool, int) {
	// declaring variables
	var lockFor int64 = 3600
	var failedloginCount int
	var msg string
	var isLocked bool

	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return "Unable to connect to database.", true, 0
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB

	var LastFailedAttempt int64
	// extracting activities on bsis of userid
	var activities []database.Activities
	db.Raw("select * from activities where email= '" + account.Email + "' order by timestamp desc limit 5;").Scan(&activities)
	for i := 0; i < len(activities); i++ {
		// if activity name is failed login and checking time interval is less then lockfor
		if activities[i].ActivityName == "failedlogin" && (time.Now().Unix()-activities[i].Timestamp) < lockFor {
			if i == 0 {
				// setting last failed attemp
				LastFailedAttempt = activities[i].Timestamp
			}
			// incrementing failed login count
			failedloginCount++
		} else {
			break
		}
	}

	// is count is more then equal to 5
	if failedloginCount >= 5 {
		msg = "Your account has been locked due to three invalid attempts. Either reset your password by clicking Forgot Password or try after " + time.Duration(1e9*(lockFor-time.Now().Unix()+LastFailedAttempt)).String() + "."
		isLocked = true
		return msg, isLocked, 0
	}

	return "", false, 5 - failedloginCount
}
