package verifyToken

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

// dummy database
func init() {
	os.Remove(os.Getenv("HOME") + "/account-testing.db")
	db, err := gorm.Open("sqlite3", os.Getenv("HOME")+"/account-testing.db")
	if err != nil {
		log.Println(err)
		log.Println("Exit")
		os.Exit(1)
	}
	config.DB = db

	//create table
	database.CreateDatabaseTables()
	token := database.Tokens{}
	token.Timestamp = 1621511267
	token.Token = "222222"
	token.TokenTask = "email_verification"
	token.Userid = "1"

	token2 := database.Tokens{}
	token2.Timestamp = 1621511267
	token2.Token = "222222"
	token2.TokenTask = "verification"
	token2.Userid = "2"

	db.Create(&token)
}
func TestCheckSentToken(t *testing.T) {
	token := CheckSentToken("1", "email_verification")
	if token != "222222" {
		t.Error("test case fail")
	}

	CheckSentToken("2", "verification")

}

func TestNewToken(t *testing.T) {

	//token with task --> email varification
	token := newToken("1", "email_verification")

	if len(token) >= 7 {
		t.Error("test case fail", token)
	}

	//token with task --> ""
	token = newToken("1", "")
	log.Println("token", token)
	if len(token) <= 6 {
		t.Error("test case fail", token)
	}

}

// test to check token is valid
func TestCheckToken(t *testing.T) {

	// test with ivalid timestamp
	_, err := CheckToken("222222")
	if err != nil {
		t.Error("test case fail")
	}

	_, err = CheckToken("888888")
	if err == nil {
		t.Error("test case fail")
	}

}

func TestDeleteToken(t *testing.T) {
	DeleteToken("222222")
}

func TestToggleOTP(t *testing.T) {
	ToggleOTP("882538")
}
