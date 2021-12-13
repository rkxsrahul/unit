package workspace

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
	email         string = "test@testing.com"
	password      string = "RKrahulkumar@321"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
	name          string = "xenon"
	roleid        string = "user"
	userid        string = "1"
	varifystatus  string = "verified"
	token         string = "222222"
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
	acc.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc.RoleID = roleid
	acc.Userid = "1"
	acc.VerifyStatus = varifystatus

	acc2 := database.Accounts{}
	acc2.AccountStatus = accountstatus
	acc2.ContactNo = contactno
	acc2.CreationDate = creationdate
	acc2.Email = "te@testing.com"
	acc2.Name = name
	acc2.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc2.RoleID = roleid
	acc2.Userid = "1"
	acc2.VerifyStatus = "not_verified"

	acc3 := database.Accounts{}
	acc3.AccountStatus = accountstatus
	acc3.ContactNo = contactno
	acc3.CreationDate = creationdate
	acc3.Email = "t@testing.com"
	acc3.Name = name
	acc3.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc3.RoleID = roleid
	acc3.Userid = "1"
	acc3.VerifyStatus = varifystatus

	work := database.Workspaces{}
	work.WorkspaceID = "xenonstack"
	work.UsagePolicies = "usage"
	work.Status = "new"
	work.TeamName = "team"
	work.TeamSize = "11"
	work.TeamType = "admin"
	work.Created = 1621511267

	work2 := database.Workspaces{}
	work2.WorkspaceID = "xenonstacks"
	work2.UsagePolicies = "usage"
	work2.Status = "new"
	work2.TeamName = "team"
	work2.TeamSize = "11"
	work2.TeamType = "admin"
	work2.Created = 1621511267

	token2 := database.Tokens{}
	token2.Timestamp = 1621511267
	token2.Token = "222222"
	token2.TokenTask = "recover_workspace"
	token2.Userid = "1"

	token3 := database.Tokens{}
	token3.Timestamp = 1621511267
	token3.Token = "333333"
	token3.TokenTask = "email_verification"
	token3.Userid = "1"

	token4 := database.Tokens{}
	token4.Timestamp = 1621511267
	token4.Token = "444444"
	token4.TokenTask = "email_verification"
	token4.Userid = "1"

	db.Create(&work)
	db.Create(&acc)
	db.Create(&work2)
	db.Create(&acc2)
	db.Create(&acc3)
	db.Create(&token2)

}

func TestCheckWorkSpaceAvailability(t *testing.T) {
	span := opentracing.StartSpan("Check Work Space Availability")
	data := database.Workspaces{
		WorkspaceID: "",
	}
	status, _ := CheckWorkspaceAvailability(data, span)
	if status != 200 {
		t.Error("test case fail")
	}
}

func TestLogin(t *testing.T) {
	span := opentracing.StartSpan("simple changepassword")
	status, _ := Login("xenonstacks", span)
	if status != 200 {
		t.Error("test case fail")
	}
	status, _ = Login("stack", span)
	if status == 200 {
		t.Error("test case fail")
	}

}

func TestCreateWorkspace(t *testing.T) {
	span := opentracing.StartSpan("simple changepassword")
	data := database.Workspaces{
		WorkspaceID:   "xenonstac",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ := CreateWorkspace("testing@xenonstack", token, data, span)
	if status != 501 {
		t.Error("test case fail", status)
	}

	data = database.Workspaces{
		WorkspaceID:   "xenonstac",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ = CreateWorkspace(email, token, data, span)
	if status != 200 {
		t.Error("test case fail", status)
	}

	data = database.Workspaces{
		WorkspaceID:   "xenonstacks",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "team",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ = CreateWorkspace(email, token, data, span)
	if status != 409 {
		t.Error("test case fail", status)
	}

	data = database.Workspaces{
		WorkspaceID:   "",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "team",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ = CreateWorkspace(email, token, data, span)
	if status != 400 {
		t.Error("test case fail", status)
	}

	data2 := database.Workspaces{
		WorkspaceID:   "xenonstac",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "team",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ = CreateWorkspace("Xenon@test.com", token, data2, span)
	if status != 409 {
		t.Error("test case fail", status)
	}

	data = database.Workspaces{
		WorkspaceID:   "",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	status, _ = CreateWorkspace(email, token, data, span)
	if status != 400 {
		t.Error("test case fail", status)
	}

}
