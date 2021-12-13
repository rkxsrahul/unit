package jwtToken

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
	userid        string = "1"
	password      string = "RKrahu@321"
	email         string = "test@testing.com"
	name          string = "xenon"
	contactno     string = "8825383227"
	varifystatus  string = "verified"
	roleid        string = "user"
	accountstatus string = "active"
	creationdate  int64  = 12345
	workpsace     string = "workshop"
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

	act := database.ActiveSessions{}
	act.ClientAgent = "client"
	act.End = 1800
	act.SessionID = "222222"
	act.Start = 1700

	db.Create(&act)

}

//test to check fetching of jwt token
func TestJwtToken(t *testing.T) {
	span := opentracing.StartSpan("Jwt token")

	acc := database.Accounts{
		Userid: userid,
		Email:  email,
		Name:   name,
		RoleID: roleid,
	}
	JwtToken(acc, span)

}

func TestJwtRefreshToken(t *testing.T) {
	span := opentracing.StartSpan("Jwt token")
	claims := make(map[string]interface{})
	// populate claims map
	claims["id"] = userid
	claims["name"] = name
	claims["email"] = email
	claims["sys_role"] = roleid
	JwtRefreshToken(claims, span)
}

func TestJwtTokenusingWID(t *testing.T) {
	acc := database.Accounts{
		Userid: userid,
		Email:  email,
		Name:   name,
		RoleID: roleid,
	}
	JwtTokenusingWID(acc, "", "user")

	acc = database.Accounts{
		Userid: userid,
		Email:  email,
		Name:   name,
		RoleID: roleid,
	}
	JwtTokenusingWID(acc, "xenonstack", "user")
}

func TestDeleteTokenFromDb(t *testing.T) {
	DeleteTokenFromDb("222222")
}
