package admin

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
	workspaceid   string = "1"
	usagepolicies string = "usagepolices"
	status        string = "active"
	teamname      string = "teamname"
	teamsize      string = "11"
	teamtype      string = "teamtype"
	creaated      int64  = 1620
)

const (
	email         string = "test@testing.com"
	password      string = "RKrahulkumar@321"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
	name          string = "xenon"
	roleid        string = "user"
	userid        string = "1"
	varifystatus  string = "verified"
	tok           string = "222222"
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

	//account data

	acc := database.Accounts{}
	acc.AccountStatus = accountstatus
	acc.ContactNo = contactno
	acc.CreationDate = creationdate
	acc.Email = "test@testing.com"
	acc.Name = name
	acc.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc.RoleID = roleid
	acc.Userid = "1"
	acc.VerifyStatus = varifystatus

	//token data
	token := database.Tokens{}
	token.Timestamp = creationdate
	token.Token = tok
	token.Userid = userid
	token.TokenTask = "recover_workspace"

	//create table for member, account and token

	work := database.Workspaces{}
	work.WorkspaceID = "xenonstack"
	work.UsagePolicies = "usage"
	work.Status = "new"
	work.TeamName = "team"
	work.TeamSize = "11"
	work.TeamType = "admin"
	work.Created = 1621511267

	member := database.WorkspaceMembers{}
	member.MemberEmail = email
	member.Role = "owner"
	member.WorkspaceID = "xenonstack"
	member.Joined = 1621511267

	db.Create(&acc)
	db.Create(&token)
	db.Create(&work)
	db.Create(&member)

}
func TestDeleteAccount(t *testing.T) {
	span := opentracing.StartSpan("simple changepassword")
	err := DeleteAccount("test@testing.com", span)
	if err != nil {
		t.Error("test case fail")
	}

	err = DeleteAccount("testing@testing.com", span)
	if err == nil {
		t.Error("test case fail")
	}
}
