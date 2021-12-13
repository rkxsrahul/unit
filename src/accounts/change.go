package accounts

import (
	"errors"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// ChangePassword is method for updating password
func ChangePassword(userid, oldpassword, password string, parentspan opentracing.Span) (int, bool, string) {
	// start span from parent span context
	span := opentracing.StartSpan("change password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	//validation check on password
	span.LogKV("task", "check password validation")
	if !methods.CheckPassword(password) {
		return 400, false, "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."
	}

	// connecting to db
	// span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return 500, false, "Unable to connect to database"
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	//check the exist password

	acc := database.Accounts{}
	db.Where("userid=?", userid).Find(&acc)

	if !methods.CheckHashForPassword(acc.Password, oldpassword) {
		return 400, false, "please pass valid current password"
	}

	if oldpassword == password {
		return 400, false, "Current and New password are same."
	}

	// creating hash for new password
	span.LogKV("task", "encrypt the password")
	newPassHash := methods.HashForNewPassword(password)

	//updating password in db is there is user with that userid passed in parameter
	span.LogKV("task", "update the database")
	dbResult := db.Exec("update accounts set password= '" + newPassHash + "' where userid= '" + userid + "';")
	acc = GetAccountForUserid(userid)
	if dbResult.Error == nil && dbResult.RowsAffected != 0 {

		err := ldap.UpdateLDAPAccountPass(acc.Email, acc.Name, password)
		if err != nil {
			span.LogKV("task", "Unable to change password in LDAP.")
			log.Println("Unable to change password in LDAP.")
			return 400, false, "Unable to change password."
		}
		span.LogKV("task", "Password updated successfully.")
		return 200, true, "Password updated successfully."
	}
	span.LogKV("task", "Unable to change password.")
	return 400, false, "Unable to change password."
}

// ======================================================================================= //
// UpdateProfile is a method to update name and contact of user
func UpdateProfile(email, name, contact string, parentspan opentracing.Span) error {
	// start span from parent span context
	span := opentracing.StartSpan("update profile method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connecting to db
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	//update name of user
	span.LogKV("task", "update the database")
	var row int64
	if name != "" {
		row = db.Model(&database.Accounts{}).Where("email=?", email).Update("name", name).RowsAffected
	}
	// update contact number of user
	if contact != "" {
		row = db.Model(&database.Accounts{}).Where("email=?", email).Update("contact_no", contact).RowsAffected
	}
	span.LogKV("task", "send final output")
	if row == 0 {
		return errors.New("no account found")
	}
	err := ldap.UpdateLDAPAccountPass(email, name, "")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
