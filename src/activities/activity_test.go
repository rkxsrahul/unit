package activities

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
	activityname string = "activity_name"
	email        string = "test@testing.com"
	clintip      string = "clint_ip"
	clintagent   string = "clintagaent"
	timestamp    int64  = 1620650852
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

	//create account database
	acc := database.Activities{}
	acc.ActivityName = activityname
	acc.Email = email
	acc.ClientIP = clintip
	acc.ClientAgent = clintagent
	acc.Timestamp = timestamp

	db.Create(&acc)

}

// test to get login activities
func TestGetLoginActivities(t *testing.T) {
	var data []database.Activities
	var err error

	//test with valid email id
	data, err = GetLoginActivities(email)
	log.Println("data", data)
	if len(data) == 0 {
		t.Error("test case fail", err)
	}

	//test with invalid email
	data, _ = GetLoginActivities("test@xenonstack")
	if len(data) != 0 {
		t.Error("test case fail")
	}

}

func TestRecordActivity(t *testing.T) {

	data := database.Activities{
		Email:        email,
		ActivityName: activityname,
		ClientIP:     clintip,
		ClientAgent:  clintagent,
		Timestamp:    timestamp,
	}
	RecordActivity(data)
}
