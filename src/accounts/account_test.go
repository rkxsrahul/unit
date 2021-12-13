package accounts

import (
	"fmt"
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

var (
	userid string = "1"
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
	acc2.Email = "t@testing.com"
	acc2.Name = name
	acc2.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc2.RoleID = "admin"
	acc2.Userid = "2"
	acc2.VerifyStatus = varifystatus

	acc3 := database.Accounts{}
	acc3.AccountStatus = "new"
	acc3.ContactNo = contactno
	acc3.CreationDate = creationdate
	acc3.Email = "tx@testing.com"
	acc3.Name = name
	acc3.Password = "AhVuY.k3s2hoNm9yux51ufOP9xpiotozw="
	acc3.RoleID = "admin"
	acc3.Userid = "5"
	acc3.VerifyStatus = varifystatus

	vpn := database.OpenVPNInformation{}
	vpn.Email = email
	vpn.FileName = "tomal"
	vpn.Password = password
	vpn.Username = "user"

	token := database.Tokens{}
	token.Timestamp = 1621511267
	token.Token = "222222"
	token.TokenTask = "email_verification"
	token.Userid = "1"

	token2 := database.Tokens{}
	token2.Timestamp = 1621511267
	token2.Token = "22222"
	token2.TokenTask = "email_verification"
	token2.Userid = "5"

	db.Create(&acc)
	db.Create(&token)
	db.Create(&vpn)
	db.Create(&acc2)
	db.Create(&acc3)

}

func TestChangePassword(t *testing.T) {

	span := opentracing.StartSpan("simple changepassword")
	_, status, _ := ChangePassword("1", "RKrahul@321", "RKrahulkumar@321", span)
	if status == true {
		t.Error("test case fail")
	}

	_, status, _ = ChangePassword("1", "RKrahul@321", "RKra", span)
	if status == true {
		t.Error("test case fail")
	}

	_, status, _ = ChangePassword("1", "RKra", "RKrahulkumar@321", span)
	if status == true {
		t.Error("test case fail")
	}

	_, status, _ = ChangePassword("1", "RKrahul@321", "RKrahul@321", span)
	if status == true {
		t.Error("test case fail")
	}

}
func TestUpdateProfile(t *testing.T) {
	span := opentracing.StartSpan("simple changepassword")
	err := UpdateProfile(email, name, contact, span)
	if err == nil {
		t.Error("test case fail")
	}

	err = UpdateProfile("rahul@gmail.com", name, contact, span)
	if err == nil {
		t.Error("test case fail")
	}

}

func TestVerifyMail(t *testing.T) {

	span := opentracing.StartSpan("simple changepassword")
	_, mail_status := VerifyMail(email, tok, span)
	if mail_status != true {
		t.Error("test case fail")
	}

	_, mail_status = VerifyMail("Xenon@xenonstack.com", tok, span)
	if mail_status == true {
		t.Error("test case fail")
	}

	_, mail_status = VerifyMail(email, "333333", span)
	if mail_status == true {
		t.Error("test case fail")
	}

	_, mail_status = VerifyMail("tx@testing.com", tok, span)
	if mail_status == true {
		t.Error("test case fail")
	}

}

//test case -> to get account from id
func TestGetAccountForUserid(t *testing.T) {

	//test case 1 -> right userid
	data := GetAccountForUserid(userid)
	fmt.Println("data", data)
	if data.Name != "test" {
		t.Error("test case fail", data)
	}

	//test case 2 -> wrong userid
	data = GetAccountForUserid("3")
	fmt.Println("data", data)
	if data.Name == "test" {
		t.Error("test case fail", data)

	}

}

//test case -> to get account from email
func TestGetAccountForEmail(t *testing.T) {

	//test case 1-> right email
	data, err := GetAccountForEmail(email)
	fmt.Println("err", err)
	fmt.Println("data", data)
	if data.Name != "test" {
		t.Error("test case fail", data)
	}

	//test case 2 -> wrong email
	data, err = GetAccountForEmail("testing@xenonstack.com")
	fmt.Println("err", err)
	if data.Name == "test" {
		t.Error("test case fail", data)
	}

}

//test case -> to get all account
func TestGetAllAccounts(t *testing.T) {

	data, _ := GetAllAccounts()

	if len(data) == 0 {
		t.Error("test case fail", data)
	}
}

// test case -> to delete account using email
func TestDeleteAccount(t *testing.T) {
	data := DeleteAccount(email)

	if data != nil {
		t.Error("test case fail", data)
	}

	data = DeleteAccount("testing@xenonstack.com")
	if data == nil {
		t.Error("test case fail", data)
	}

	data = DeleteAccount("t@testing.com")
	if data == nil {
		t.Error("test case fail", data)
	}

}

//=====================================================================================================
