package verifyToken

import (
	"errors"
	"log"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// ToggleOTP is a function to change the value of OTP in config
func ToggleOTP(value string) {
	config.OTP = value
}

// CheckSentToken is a method to generate verifications token that are send in mail
// Before generating new token it checks is there any valid old token previouly
func CheckSentToken(userid, task string) string {

	// connecting to db
	db := config.DB

	// cheking previous token
	var tok []database.Tokens
	if task == "email_verification" {
		db.Where("userid= ? AND token_task= ? ", userid, task).Find(&tok)
	} else {
		db.Where("userid= ? AND token_task= ?", userid, task).Find(&tok)
	}
	token := ""
	if len(tok) != 0 {
		// if there is token
		token = tok[0].Token
	} else {
		// else creating new token
		token = newToken(userid, task)
	}
	return token
}

//====================================================================

// newToken is a method to generate new verification token and add to database
func newToken(userid, task string) string {

	// connecting to db
	db := config.DB
	token := database.Tokens{}
	// setting token user id
	token.Userid = userid
	// generating token on basis of task
	switch task {
	case "email_verification":
		// generating random 6 digit numeric string
		token.Token = methods.RandomStringIntegerOnly(6)
	default:
		// for other tasks creating 35 length random string
		token.Token = methods.RandomString(35)
	}

	if config.OTP == "true" {
		token.Token = "111111"
	}

	token.TokenTask = task
	token.Timestamp = time.Now().Unix()

	// save data in db
	db.Create(&token)
	return token.Token
}

//========================================================================

// CheckToken is a method to check token is valid or not
func CheckToken(token string) (database.Tokens, error) {
	// connecting to db
	db := config.DB

	// fetch token details
	tok := []database.Tokens{}
	db.Where("token=?", token).Find(&tok)
	//token not found
	if len(tok) == 0 {
		log.Println("token not found")
		return database.Tokens{}, errors.New("Invalid or expired token")
	}

	//expired token
	// if (time.Now().Unix() - tok[0].Timestamp) > config.Conf.Service.InviteLinkTimeout {
	// 	log.Println("expired token")
	// 	return database.Tokens{}, errors.New("Invalid or expired token")
	// }
	return tok[0], nil
}

//========================================================================

// DeleteToken is method to delete used and expired tokens
func DeleteToken(token string) {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// delete used token
	row := db.Where("token=?", token).Delete(&database.Tokens{}).RowsAffected
	log.Println(row)
	// delete expired tokens
	row = db.Where("timestamp < ?", strconv.FormatInt((time.Now().Unix()-config.Conf.Service.InviteLinkTimeout), 10)).RowsAffected
	log.Println(row)
}
