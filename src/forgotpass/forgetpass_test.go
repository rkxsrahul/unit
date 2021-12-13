package forgotpass

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	userid               = "1"
	email         string = "test@testing.com"
	name          string = "test"
	roleid        string = "user"
	verifystatus  string = "verified"
	filename      string = "rahul"
	password      string = "Rk123rahul@"
	username      string = "xyz"
	oldpassword   string = "old_pass"
	contact       string = "8825383117"
	code          string = "xenonstack"
	wsID          string = "workshop"
	accountstatus string = "active"
	token         string = "token"
)

const (
	workshop    string = "workshop"
	memberemail string = "member_email"
	role        string = "workspace_role"
	joined      int64  = 1620650852
)

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

	//forget password database
	forget := database.WorkspaceMembers{}
	forget.WorkspaceID = wsID
	forget.Joined = joined
	forget.MemberEmail = email
	forget.Role = role
	//token database
	tok := database.Tokens{}
	tok.Timestamp = joined
	tok.Token = token
	tok.TokenTask = "forgot_pass"
	tok.Userid = userid

	tok3 := database.Tokens{}
	tok3.Timestamp = joined
	tok3.Token = "333333"
	tok3.TokenTask = "forgot_pass"
	tok3.Userid = "3"

	tok2 := database.Tokens{}
	tok2.Timestamp = joined
	tok2.Token = "222222"
	tok2.TokenTask = "forgot"
	tok2.Userid = "2"

	//account database
	acc := database.Accounts{}
	acc.AccountStatus = accountstatus
	acc.ContactNo = contact
	acc.Email = email
	acc.Name = name
	acc.Password = password
	acc.Userid = userid
	acc.VerifyStatus = verifystatus

	acc2 := database.Accounts{}
	acc2.AccountStatus = accountstatus
	acc2.ContactNo = contact
	acc2.Email = "x123@gmail.com"
	acc2.Name = name
	acc2.Password = password
	acc2.Userid = "11"
	acc2.VerifyStatus = verifystatus

	db.Create(&acc)
	db.Create(&acc2)
	db.Create(&tok)
	db.Create(&forget)
	db.Create(&tok2)
	db.Create(&tok3)

}

//test for forget password challenge
func TestForgotPassChallenge(t *testing.T) {
	span := opentracing.StartSpan("simple forgetpassword")

	//test with valid email
	_, bol := ForgotPassChallenge(email, wsID, span)
	log.Println("xxxxx", bol)
	ForgotPassChallenge("rahul@gmail.com", wsID, span)
	ForgotPassChallenge("x123@gmail.com", wsID, span)
}

//test reset forgetten password
func TestResetForgottenPass(t *testing.T) {
	span := opentracing.StartSpan("simple ResetForgottenPass")

	//test with valid token
	ResetForgottenPass(token, password, wsID, span)

	//test with invalid token
	ResetForgottenPass("xenon", password, wsID, span)
	ResetForgottenPass(token, password, "", span)
	ResetForgottenPass(token, "rk", wsID, span)
	ResetForgottenPass("222222", password, wsID, span)
	ResetForgottenPass("333333", password, wsID, span)

}

//test for password challenge
func TestPassChallenge(t *testing.T) {

	span := opentracing.StartSpan("simple PassChallenge")

	//test with valid email
	_, bol := forgotPassChallenge(email, span)
	if bol != true {
		t.Error("test case fail")
	}

	//test with invalid email
	_, bol = forgotPassChallenge("xenon@testing", span)
	if bol == true {
		t.Error("test case fail")
	}
}

//test for reset password
func TestResetPass(t *testing.T) {
	span := opentracing.StartSpan("simple TestResetPass")

	resetForgottenPass(token, password, span)
	resetForgottenPass(token, "rk", span)
	resetForgottenPass("tok", password, span)
	resetForgottenPass("222222", password, span)
	resetForgottenPass("333333", password, span)

}

//test for update database
func TestUpdateDatabase(t *testing.T) {

	//test with valid password
	_, err := updateDatabase(userid, password)
	if err != nil {
		t.Error("test case fail")
	}
	//test with invalid password
	_, err = updateDatabase(userid, "xenon")
	if err != nil {
		t.Error("test case fail")
	}

	data, _ := updateDatabase("10", "xenon")
	if data != "" {
		t.Error("test case fail")
	}
}
