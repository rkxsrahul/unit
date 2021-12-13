package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// stores aws access keys data
type AccessKeys struct {
	Userid string `gorm:"primary_key"`
	Key    string //if you require to store token secret then append that to it with '&' symbol.
	Secret string
}

type OpenVPNInformation struct {
	Email string `gorm:"primary_key"`
	//	Email    string `gorm:"unique"`
	Username string `json:"username"`
	Password string `json:"password"`
	FileName string `json:"filename"`
}

// stores details about aws iam policy
type Policy struct {
	// arn of policy
	Arn string `json:"arn" binding:"required"`
	// policy type
	// values accepted all, Individual, Enterprise, Admin
	PType string `json:"type" binding:"required"`
	// if enterprise or admin then company name should be there
	Company string `json:"company"`
}

func CreateDatabaseTables() {
	// connecting db using connection string
	db := config.DB

	// creating all tables one by one but firstly checking whether table exists or not
	if !(db.HasTable(Flags{})) {
		db.CreateTable(Flags{})
	}
	if !(db.HasTable(Accounts{})) {
		db.CreateTable(Accounts{})
		// initializing flag table
		db.Exec("insert into flags (name, value, usage) values ('new_userid', '10000', 'to allocate userid to new user.');")

		//creating admin account
		adminAcc := InitAdminAccount()
		db.Create(&adminAcc)

		err := createLDAPAccount(adminAcc, config.Conf.Admin.Pass)
		if err != nil {
			log.Println("LDAP account creation error: ", err)
		}
	}

	if !(db.HasTable(Activities{})) {
		db.CreateTable(Activities{})
	}
	if !(db.HasTable(Tokens{})) {
		db.CreateTable(Tokens{})
	}
	if !(db.HasTable(ActiveSessions{})) {
		db.CreateTable(ActiveSessions{})
	}

	if !(db.HasTable(Workspaces{})) {
		db.CreateTable(Workspaces{})
	}
	if !(db.HasTable(WorkspaceMembers{})) {
		db.CreateTable(WorkspaceMembers{})
	}
	if !(db.HasTable(AccessKeys{})) {
		db.CreateTable(AccessKeys{})
	}
	if !(db.HasTable(Policy{})) {
		db.CreateTable(Policy{})
	}
	if !(db.HasTable(VMRequestInfo{})) {
		db.CreateTable(VMRequestInfo{})
	}
	if !(db.HasTable(OpenVPNInformation{})) {
		db.CreateTable(OpenVPNInformation{})
	}

	// Database migration
	db.AutoMigrate(&Flags{},
		&Accounts{},
		&Activities{},
		&Tokens{},
		&ActiveSessions{},
		&Workspaces{},
		&WorkspaceMembers{},
		&AccessKeys{},
		&Policy{},
		VMRequestInfo{},
		OpenVPNInformation{})

	// Add foreignKeys
	db.Model(&Activities{}).AddForeignKey("email", "accounts(email)", "CASCADE", "CASCADE")
	db.Model(&ActiveSessions{}).AddForeignKey("userid", "accounts(userid)", "CASCADE", "CASCADE")
	db.Model(&Tokens{}).AddForeignKey("userid", "accounts(userid)", "CASCADE", "CASCADE")
	db.Model(&WorkspaceMembers{}).AddForeignKey("workspace_id", "workspaces(workspace_id)", "CASCADE", "CASCADE")
	db.Model(&AccessKeys{}).AddForeignKey("userid", "accounts(userid)", "CASCADE", "CASCADE")
}

func CreateDatabase() {
	// connecting with cockroach database root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.User,
		config.Conf.Database.Pass,
		"postgres", config.Conf.Database.Ssl))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	// executing create database query.
	db.Exec(fmt.Sprintf("create database %s;", config.Conf.Database.Name))
}

func createLDAPAccount(acs Accounts, password string) error {

	log.Println("==============================================")
	ldapPort, err := strconv.Atoi(config.Conf.LDAP.Port)
	if err != nil {
		return err
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Conf.LDAP.Host, ldapPort))
	if err != nil {
		log.Println(err)
		return err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.Conf.LDAP.AdminDN, config.Conf.LDAP.AdminPass)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(acs.RoleID)
	addReq := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", acs.Email, config.Conf.LDAP.UserParentDN))
	addReq.Attribute("objectClass", []string{"top", "inetOrgPerson", "posixAccount"})
	addReq.Attribute("userPassword", []string{password})
	addReq.Attribute("cn", []string{acs.Name})
	addReq.Attribute("sn", []string{acs.Name})
	//	addReq.Attribute("role", []string{acs.RoleID})
	addReq.Attribute("mail", []string{acs.Email})
	addReq.Attribute("uidNumber", []string{acs.Userid})
	addReq.Attribute("gidNumber", []string{"1001"})
	addReq.Attribute("homeDirectory", []string{"/home/" + acs.Userid})
	err = l.Add(addReq)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("==============================================")
	return nil
}
