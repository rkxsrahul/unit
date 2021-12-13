package login

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
	acc.AccountStatus = "active"
	acc.ContactNo = contactno
	acc.CreationDate = creationdate
	acc.Email = "test@testing.com"
	acc.Name = name
	acc.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc.RoleID = roleid
	acc.Userid = "1"
	acc.VerifyStatus = varifystatus

	acc2 := database.Accounts{}
	acc2.AccountStatus = "active"
	acc2.ContactNo = contactno
	acc2.CreationDate = creationdate
	acc2.Email = "testing@testing.com"
	acc2.Name = name
	acc2.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc2.RoleID = roleid
	acc2.Userid = "2"
	acc2.VerifyStatus = varifystatus

	acc3 := database.Accounts{}
	acc3.AccountStatus = "active"
	acc3.ContactNo = contactno
	acc3.CreationDate = creationdate
	acc3.Email = "t@testing.com"
	acc3.Name = name
	acc3.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc3.RoleID = roleid
	acc3.Userid = "3"
	acc3.VerifyStatus = varifystatus

	//token data
	token := database.Tokens{}
	token.Timestamp = creationdate
	token.Token = tok
	token.Userid = userid
	token.TokenTask = "recover_workspace"

	token2 := database.Tokens{}
	token2.Timestamp = creationdate
	token2.Token = "222222"
	token2.Userid = "2"
	token2.TokenTask = "recover"

	token3 := database.Tokens{}
	token3.Timestamp = creationdate
	token3.Token = "444444"
	token3.Userid = "2"
	token3.TokenTask = "recover_workspace"

	token4 := database.Tokens{}
	token4.Timestamp = creationdate
	token4.Token = "444445"
	token4.Userid = "3"
	token4.TokenTask = "recover_workspace"

	//member data
	member := database.WorkspaceMembers{}
	member.WorkspaceID = "1"
	member.Joined = 3600
	member.Role = "owner"
	member.MemberEmail = "test@testing.com"

	member2 := database.WorkspaceMembers{}
	member2.WorkspaceID = "2"
	member2.Joined = 3600
	member2.Role = "user"
	member2.MemberEmail = "testing@testing.com"

	work := database.Workspaces{}
	work.WorkspaceID = "1"
	work.UsagePolicies = "usage"
	work.Status = "new"
	work.TeamName = "team"
	work.TeamSize = "11"
	work.TeamType = "admin"
	work.Created = 1621511267

	//create table for member, account and token
	db.Create(&member)
	db.Create(&acc)
	db.Create(&acc3)
	db.Create(&token)
	db.Create(&acc2)
	db.Create(&member2)
	db.Create(&token2)
	db.Create(&token3)
	db.Create(&token4)
	db.Create(&work)
}

//test to check user
func TestCheckUser(t *testing.T) {
	acc := database.Accounts{
		AccountStatus: "active",
		ContactNo:     contactno,
		CreationDate:  creationdate,
		Email:         "test@testing.com",
		Name:          name,
		Password:      "RKrahul@321",
		RoleID:        roleid,
		Userid:        userid,
		VerifyStatus:  varifystatus,
	}

	//test user with invalid timestamp
	checkUser(acc)

	acc2 := database.Accounts{
		AccountStatus: "active",
		ContactNo:     contactno,
		CreationDate:  creationdate,
		Email:         "xenon@test.com",
		Name:          name,
		Password:      password,
		RoleID:        roleid,
		Userid:        userid,
		VerifyStatus:  varifystatus,
	}

	//test user with invalid timestamp

	checkUser(acc2)

	acc3 := database.Accounts{
		AccountStatus: "active",
		ContactNo:     contactno,
		CreationDate:  creationdate,
		Email:         "testing@testing.com",
		Name:          name,
		Password:      password,
		RoleID:        roleid,
		Userid:        userid,
		VerifyStatus:  varifystatus,
	}

	//test user with invalid timestamp

	checkUser(acc3)

}

func TestNormalLogin(t *testing.T) {
	span := opentracing.StartSpan("simple signin")
	code, _ := NormalLogin("test@testing.com", "RKrahul@321", "1", span)
	if code != 401 {
		t.Error("test case fail")
	}

	code, _ = NormalLogin(email, password, "1", span)
	if code != 401 {
		t.Error("test case fail")
	}

}
