package signup

import (
	"log"
	"os"
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/opentracing/opentracing-go"
)

const (
	email         string = "test@testing.com"
	name          string = "test"
	roleid        string = "user"
	verifystatus  string = "verified"
	filename      string = "rahul"
	password      string = "RKrahulkumar@321"
	username      string = "xyz"
	oldpassword   string = "bGvaR.ByOQKrzjb85wiUanwHUFiJ+74KI="
	contact       string = "8825383117"
	code          string = "xenonstack"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
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

	acc := database.Accounts{}
	acc.AccountStatus = accountstatus
	acc.ContactNo = contactno
	acc.CreationDate = creationdate
	acc.Email = "test@testing.com"
	acc.Name = name
	acc.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc.RoleID = roleid
	acc.Userid = "1"
	acc.VerifyStatus = varifystatus

	acc2 := database.Accounts{}
	acc2.AccountStatus = accountstatus
	acc2.ContactNo = contactno
	acc2.CreationDate = creationdate
	acc2.Email = "test@test.com"
	acc2.Name = name
	acc2.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc2.RoleID = roleid
	acc2.Userid = "2"
	acc2.VerifyStatus = "notverified"

	db.Create(&acc)
	db.Create(&acc2)

}

func TestSendCodeAgain(t *testing.T) {
	span := opentracing.StartSpan("send code")
	_, status := SendCodeAgain(email, span)
	if status != false {
		t.Error("test case fail")
	}
	_, status = SendCodeAgain("test@test.com", span)
	if status != true {
		t.Error("test case fail")
	}
	_, status = SendCodeAgain("test@xenon.com", span)
	if status != false {
		t.Error("test case fail")
	}
}

func TestInitNewAccount(t *testing.T) {
	acc := database.Accounts{
		Userid:        "1",
		VerifyStatus:  varifystatus,
		AccountStatus: accountstatus,
		RoleID:        roleid,
	}
	initNewAccount(&acc)
}
