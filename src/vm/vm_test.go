package vm

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

const (
	email         string = "test@testing.com"
	password      string = "testing"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
	name          string = "xenon"
	roleid        string = "admin"
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

	//account database
	acc := database.Accounts{}
	acc.AccountStatus = accountstatus
	acc.ContactNo = contactno
	acc.CreationDate = creationdate
	acc.Email = email
	acc.Name = name
	acc.Password = password
	acc.RoleID = roleid
	acc.Userid = userid
	acc.VerifyStatus = varifystatus

	//token database
	token := database.Tokens{}
	token.Timestamp = creationdate
	token.Token = tok
	token.Userid = userid
	token.TokenTask = "recover_workspace"

	//create table for member, account and token

	//workspace database
	work := database.Workspaces{}
	work.WorkspaceID = "2"
	work.UsagePolicies = "usage"
	work.Status = "new"
	work.TeamName = "team"
	work.TeamSize = "11"
	work.TeamType = "admin"
	work.Created = 3600

	//VM database
	vm := database.VMRequestInfo{}
	vm.Description = "Description"
	vm.Flavour = "Flavour"
	vm.ID = 1
	vm.Name = name
	vm.Ports = "inbound_ports"
	vm.Source = "source"
	vm.UserEmail = email
	vm.UserName = "user"
	vm.Workspace = "xenonstack"

	db.Create(&acc)
	db.Create(&token)
	db.Create(&work)
	db.Create(&vm)

}

//test for VM request
func TestVMRequest(t *testing.T) {
	info := database.VMRequestInfo{
		Description: "Description",
		Flavour:     "Flavour",
		ID:          2,
		Name:        name,
		Ports:       "inbound_ports",
		Source:      "source",
		UserEmail:   email,
		UserName:    "user",
		Workspace:   "xenonstack",
	}

	//test with valid information
	status := VMRequest(info)
	if status != nil {
		t.Error("test case fail", status)
	}

}
