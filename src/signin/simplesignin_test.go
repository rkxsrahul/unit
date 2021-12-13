package signin

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
	email    string = "test@testing.com"
	password string = "RKrahul@321"

	name         string = "test"
	roleid       string = "user"
	verifystatus string = "verified"
	filename     string = "rahul"

	username      string = "xyz"
	oldpassword   string = "bGvaR.ByOQKrzjb85wiUanwHUFiJ+74KI="
	contact       string = "8825383117"
	code          string = "xenonstack"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
	varifystatus  string = "verified"
	tok           string = "222222"
	activityname  string = "activity_name"

	clintip    string = "clint_ip"
	clintagent string = "clintagaent"
	timestamp  int64  = 16206508524444
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
	acc2.AccountStatus = "blocked"
	acc2.ContactNo = contactno
	acc2.CreationDate = creationdate
	acc2.Email = "tes@testing.com"
	acc2.Name = name
	acc2.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc2.RoleID = roleid
	acc2.Userid = "2"
	acc2.VerifyStatus = varifystatus

	acc3 := database.Accounts{}
	acc3.AccountStatus = "new"
	acc3.ContactNo = contactno
	acc3.CreationDate = creationdate
	acc3.Email = "te@testing.com"
	acc3.Name = name
	acc3.Password = "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o="
	acc3.RoleID = roleid
	acc3.Userid = "3"
	acc3.VerifyStatus = varifystatus

	act := database.Activities{}
	act.ActivityName = activityname
	act.Email = "test@testing.com"
	act.ClientIP = clintip
	act.ClientAgent = clintagent
	act.Timestamp = timestamp

	act2 := database.Activities{}
	act2.ActivityName = "failedlogin"
	act2.Email = "testxs@testing.com"
	act2.ClientIP = clintip
	act2.ClientAgent = clintagent
	act2.Timestamp = timestamp

	db.Create(&acc)
	db.Create(&acc2)
	db.Create(&acc3)
	db.Create(&act)
	db.Create(&act2)

}
func TestSimpleSignin(t *testing.T) {
	span := opentracing.StartSpan("simple signin")
	status, _, _, _ := SimpleSignin(email, password, span)
	if status != false {
		t.Error("test case fail")
	}

	status, _, _, _ = SimpleSignin("tsx@testing.com", password, span)
	if status != false {
		t.Error("test case fail")
	}

	status, _, _, _ = SimpleSignin("tes@testing.com", password, span)
	if status != false {
		t.Error("test case fail")
	}

	status, _, _, _ = SimpleSignin("te@testing.com", password, span)
	if status != false {
		t.Error("test case fail")
	}

	status, _, _, _ = SimpleSignin(email, "RKrahulkumar@321", span)
	if status != false {
		t.Error("test case fail")
	}

}

func TestCheckPreviousFailedLogins(t *testing.T) {
	acc := database.Accounts{
		Email: "test@testing.com",
	}
	_, status, _ := checkPreviousFailedLogins(acc)
	if status != false {
		t.Error("test case fail")
	}

	acc2 := database.Accounts{
		Email: "testxs@testing.com",
	}
	_, status, _ = checkPreviousFailedLogins(acc2)
	if status != false {
		t.Error("test case fail")
	}

}
