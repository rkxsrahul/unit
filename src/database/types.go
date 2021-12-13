package database

import (
	"time"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/methods"
)

// Flags is a strucuture use to assign userid in a systematic way to a user
type Flags struct {
	Name string
	//value is a userid of user
	Value        string
	DefaultValue string
	Usage        string
}

// Accounts is a strucutre to stores user information
type Accounts struct {
	// auto generated
	Userid string `json:"userid" gorm:"not null;unique;"`
	//If password is empty then the user registered from career forms.
	//User can register password in future using forgot password.
	Password string `json:"-"`
	//email of user
	Email string `json:"email" gorm:"not null;unique_index;"`
	//name of user
	Name string `json:"name"`
	// contact number of user
	ContactNo string `json:"contact_no"`
	//Verify_status can be verified, not_verified
	VerifyStatus string `json:"verify_status"`
	//this role id is used for auth portal management
	//value can be 'admin', 'user'
	RoleID string `json:"sys_role"`
	//account status can be active, new, deleted, blocked etc.
	AccountStatus string `json:"account_status"`
	//account creation date
	CreationDate int64
}

// Activities is a structure to record user activties
type Activities struct {
	//If username entered incorrect then these activities will also be recorded.
	Email string `json:"email" gorm:"index"`
	//can be login, failedlogin, signup
	ActivityName string `json:"activity_name"`
	ClientIP     string `json:"client_ip"`
	ClientAgent  string `json:"client_agent"`
	Timestamp    int64  `json:"timestamp"`
}

// Tokens is a structure to stores token for verifcation, invite link, forgot password
type Tokens struct {
	Userid    string `gorm:"index"`
	Token     string //if you require to store token secret then append that to it with '&' symbol.
	TokenTask string
	Timestamp int64
}

// ActiveSessions is a structure to stores active sessions
type ActiveSessions struct {
	SessionID   string
	Userid      string `gorm:"index"`
	ClientAgent string
	Start       int64
	// if value is '0' then session is remembered.
	End int64
}

// Workspaces is a structure to stores workspace inofrmation
type Workspaces struct {
	// url of workspace
	WorkspaceID   string `json:"workspace_url" gorm:"unique_index;not null" binding:"required"`
	UsagePolicies string
	// Status can be new, active, deleted, blocked
	Status string `json:"status"`

	//if not provided then intial of url is team name
	TeamName string `json:"team_name"`
	TeamSize string `json:"team_size"`
	TeamType string `json:"team_type"`
	Created  int64  `json:"created" gorm:"not null"`
}

// WorkspaceMembers is a structure to store details abont members of a particular workspace
type WorkspaceMembers struct {
	WorkspaceID string `json:"workspace_id" gorm:"not null;unique_index:edq_idx"`
	MemberEmail string `json:"member_email" gorm:"not null;unique_index:edq_idx"`
	// admin or user
	Role string `json:"workspace_role" gorm:"not null;unique_index:edq_idx"`
	// date when user joined
	Joined int64 `json:"created" gorm:"not null;default:0"`
}

type VMRequestInfo struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Flavour     string `json:"flavour"`
	Ports       string `json:"inbound_ports"`
	UserEmail   string `json:"user_email"`
	UserName    string `json:"user_name"`
	Workspace   string `json:"workspace"`
}

// InitAdminAccount is a function used to create admin account
func InitAdminAccount() Accounts {

	// fetching info from env variables
	adminEmail := config.Conf.Admin.Email
	if adminEmail == "" {
		adminEmail = "admin@xenonstack.com"
	}
	adminPass := config.Conf.Admin.Pass
	if adminPass == "" {
		adminPass = "admin"
	}
	// return struct with details of admin
	return Accounts{Userid: "0",
		Password:      methods.HashForNewPassword(adminPass),
		Email:         adminEmail,
		Name:          adminEmail,
		RoleID:        "admin",
		AccountStatus: "active",
		VerifyStatus:  "verified",
		CreationDate:  time.Now().Unix()}
}
