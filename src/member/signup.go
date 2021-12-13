package member

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

// MemberSignup is a method to send invite link again if user is already invited by owner
func MemberSignup(email, wsID string, parentspan opentracing.Span) (int, string) {
	// start span from parent span context
	span := opentracing.StartSpan("member signup method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	span.LogKV("task", "check user belongs to workspace")
	if !CheckMember(email, wsID) {
		span.LogKV("task", "send final output when user not belongs to workspace")
		return 404, "You have not been invited in this workspace."
	}

	// connecting to db
	span.LogKV("task", "intialise db connection")
	db := config.DB

	owner := []database.WorkspaceMembers{}
	db.Where("role=? AND workspace_id=?", "owner", wsID).Find(&owner)

	//fetch account on basis of email
	span.LogKV("task", "fetch account details on basis of email")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "send final output when no account is there")
		log.Println(err)
		return 404, err.Error()
	}

	if acc.AccountStatus == "active" {
		span.LogKV("task", "send final output when user already registered")
		return 209, "You have already registered."
	}
	ownerAcc, err := accounts.GetAccountForEmail(owner[0].MemberEmail)
	if err != nil {
		span.LogKV("task", "send final output when no account is there")
		log.Println(err)
		return 404, err.Error()
	}

	// send mail
	span.LogKV("task", "send mail")
	go mail.SendInviteLink(acc, wsID, ownerAcc.Email, ownerAcc.Name)
	return 200, "We have sent a link to your email, Please check your email."
}

// CheckMember is a method to check user belongs to that workspace
func CheckMember(email, wsID string) bool {
	// connecting to db
	db := config.DB

	var count int64
	db.Model(&database.WorkspaceMembers{}).Where("member_email= ? AND workspace_id= ?", email, wsID).Count(&count)
	return count != 0
}
