package signup

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

// new is a constant to remove duplicacy code
const new string = "new"

// Signup is a method for creating account if account is already not preset and send mail for verification
func Signup(newAccount database.Accounts, passwrd string, parentspan opentracing.Span) (string, bool) {
	// start span from parent span contex
	span := opentracing.StartSpan("signup function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return "Unable to connect to database.", false
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB

	span.LogKV("task", "fetch account details to check account already is there or not")
	//Checking for email whether already exists or not
	var oldAccount []database.Accounts
	db.Where("email ILIKE ?", newAccount.Email).Find(&oldAccount)

	// if no account is there with that email
	if len(oldAccount) == 0 {
		// initializing new account
		span.LogKV("task", "initialize new account")
		initNewAccount(&newAccount)

		// Create ldap account
		err := ldap.CreateLDAPAccount(newAccount, passwrd)
		if err != nil {
			log.Println("LDAP account creation error: ", err)
			return "Unable to add user.", false
		}
		// creating new account
		span.LogKV("task", "save account in database")
		db.Create(&newAccount)
		// sending mail when account status is new
		if newAccount.AccountStatus == new {
			//send verification mail
			span.LogKV("task", "send mail for verification")
			go mail.SendVerifyMail(newAccount)
			return "We have sent a confirmation link to your email, please check your email.", true
		}
		span.LogKV("task", "account registered succesfully")
		return "Registered successfully.", true
	}

	userCount, err := ldap.SearchLDAPAccount(oldAccount[0].Email)
	if err != nil {
		log.Println(err)
		return "Unable to add user.", false
	}

	if userCount > 0 {
		return "Email already exists.", false
	}

	// Create ldap account
	err = ldap.CreateLDAPAccount(oldAccount[0], passwrd)
	if err != nil {
		log.Println("LDAP account creation error: ", err)
		return "Unable to add user.", false
	}
	// when there is account with that email but status is new means not verified
	if oldAccount[0].AccountStatus == "new" {
		span.LogKV("task", "delete old account details name and password")
		// delete previous account details
		db.Exec("delete from accounts where userid='" + oldAccount[0].Userid + "';")
		oldAccount[0].Name = newAccount.Name
		oldAccount[0].Password = newAccount.Password
		// creating account with new details
		span.LogKV("task", "save account in database")
		db.Create(&oldAccount[0])
		//send verification mail
		span.LogKV("task", "send mail for verification")
		go mail.SendVerifyMail(oldAccount[0])
		return "We have sent a confirmation link to your email, please check your email.", true
	}
	span.LogKV("task", "account already exists")
	return "Email Already Exists.", false
}

//==============================================================================

// WithEmail is a method used to create account with email only
func WithEmail(email string) (database.Accounts, error) {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return database.Accounts{}, errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()

	//check account already exist
	db := config.DB
	acc := []database.Accounts{}
	db.Where("email=?", email).Find(&acc)
	if len(acc) != 0 {
		return acc[0], nil
	}

	// new account structure
	newAccount := database.Accounts{
		Email:         email,
		VerifyStatus:  "not_verified",
		AccountStatus: "new",
		CreationDate:  time.Now().Unix(),
	}

	// initialize new account
	initNewAccount(&newAccount)
	// create account in database
	db.Create(&newAccount)
	return newAccount, nil
}
