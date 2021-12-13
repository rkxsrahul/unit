package member

import (
	"errors"
	"log"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/signup"
)

// output structure for invite member requests
type out struct {
	Role    string `json:"role"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type status struct {
	Add    []out `json:"add"`
	NotAdd []out `json:"not_add"`
}

type WorkSpaceMember struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

// Invite is a method to send invite mail to valid Users
func Invite(data interface{}, ownerEmail, ownerName, wsID string, parentspan opentracing.Span) (int, map[string]interface{}) {
	emails := data.([]WorkSpaceMember)

	mapd := make(map[string]interface{})
	// start span from parent span context
	span := opentracing.StartSpan("invite members method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	db := config.DB

	span.LogKV("task", "send invite link mail to each new user and login link mail to old users")
	add := make([]out, 0)
	notadd := make([]out, 0)
	for i := 0; i < len(emails); i++ {
		if emails[i].Email == ownerEmail {
			notadd = append(notadd, out{
				Role:    "owner",
				Email:   ownerEmail,
				Message: "This person is already in your workspace",
			})
			continue
		}
		if emails[i].Email == "" {
			continue
		}
		span.LogKV("task", "call auth-core method signup with only email")
		acc, err := signup.WithEmail(strings.ToLower(emails[i].Email))
		if err != nil {
			log.Println(err)
			continue
		}
		span.LogKV("task", "check user already there in workspace")
		var count int64
		db.Model(&database.WorkspaceMembers{}).Where("workspace_id= ? AND member_email= ?", wsID, acc.Email).Count(&count)
		if count == 0 {
			if acc.VerifyStatus == "not_verified" {
				// add in workspace member
				db.Create(&database.WorkspaceMembers{
					WorkspaceID: wsID,
					MemberEmail: acc.Email,
					Role:        emails[i].Role,
				})
				//send invite mail
				go mail.SendInviteLink(acc, wsID, ownerEmail, ownerName)
			} else {
				// add in workspace member
				db.Create(&database.WorkspaceMembers{
					WorkspaceID: wsID,
					MemberEmail: acc.Email,
					Role:        emails[i].Role,
					Joined:      time.Now().Unix(),
				})

				//================================aws=============================//
				if config.Conf.Service.ISAWS == "true" {
					err := awsWork(acc, wsID)
					if err != nil {
						notadd = append(notadd, out{
							Role:    emails[i].Role,
							Email:   acc.Email,
							Message: "AWS work not completed",
						})
						continue
					}
				}
				// send login link mail
				go mail.SendLoginLink(acc, wsID, ownerEmail, ownerName)
			}
			add = append(add, out{
				Role:    emails[i].Role,
				Email:   acc.Email,
				Message: "Successfully invited",
			})
		} else {
			if acc.VerifyStatus == "not_verified" {

				add = append(add, out{
					Role:    emails[i].Role,
					Email:   acc.Email,
					Message: "Successfully invited",
				})
				//send invite mail
				go mail.SendInviteLink(acc, wsID, ownerEmail, ownerName)
			} else {

				notadd = append(notadd, out{
					Role:    emails[i].Role,
					Email:   acc.Email,
					Message: "This person is already in your workspace",
				})
			}
		}
	}

	span.LogKV("task", "send final output")
	mapd["error"] = false
	mapd["status"] = status{
		Add:    add,
		NotAdd: notadd,
	}
	return 200, mapd
}

//===========================================================================//

// MemList is a structure to send member details in member list request
type MemList struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Joined int64  `json:"joined"`
}

// List is a method to fetch all member in a workspace from database
func List(workspace string, parentspan opentracing.Span) ([]MemList, error) {
	list := make([]MemList, 0)
	// start span from parent span context
	span := opentracing.StartSpan("list members method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	db := config.DB
	// member list from database
	span.LogKV("task", "fetch members from database on basis of workspace")
	wsMem := []database.WorkspaceMembers{}
	db.Where("workspace_id= ?", workspace).Order("joined asc").Find(&wsMem)

	span.LogKV("task", "fetch user name of each user")
	for i := 0; i < len(wsMem); i++ {
		// fetch account details on basis of email
		acc, err := accounts.GetAccountForEmail(wsMem[i].MemberEmail)
		if err != nil {
			log.Println(err)
			continue
		}
		// append final details
		list = append(list, MemList{
			Id:     acc.Userid,
			Name:   acc.Name,
			Email:  acc.Email,
			Role:   wsMem[i].Role,
			Joined: wsMem[i].Joined,
		})
	}

	span.LogKV("task", "send final list")
	return list, nil
}

//===========================================================================//

// Check is a method to check member joined the worksapce from database
func Check(workspace, email string, parentspan opentracing.Span) error {
	// start span from parent span context
	span := opentracing.StartSpan("check member method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")
	db := config.DB
	// check member in database
	span.LogKV("task", "check member in database")
	var count int64
	db.Model(&database.WorkspaceMembers{}).Where("member_email = ? AND workspace_id= ? AND joined <> 0", email, workspace).Count(&count)
	span.LogKV("task", "send final output")
	if count == 0 {
		return errors.New("member not invited or joined the workspace")
	}
	return nil
}

//===========================================================================//

// Delete is a method to delete member in a workspace from database
func Delete(workspace, email string, parentspan opentracing.Span) error {
	// start span from parent span context
	span := opentracing.StartSpan("delete member method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	span.LogKV("task", "intialise db connection")

	//fetch account on basis of email
	span.LogKV("task", "fetch account details on basis of email")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "send final output when no account is there")
		log.Println(err)
		return err
	}

	if config.Conf.Service.ISAWS == "true" {
		//delete aws things related to member
		err = awsDeleteWork(acc, workspace)
		if err != nil {
			return errors.New("Unable to delete aws bucket. Please try again later")
		}
	}
	db := config.DB

	// delete member from database
	span.LogKV("task", "delete member from database")
	count := db.Where("member_email = ? AND workspace_id= ?", email, workspace).Delete(&database.WorkspaceMembers{}).RowsAffected
	span.LogKV("task", "send final output")
	if count == 0 {
		return errors.New("this member is not invited in the workspace")
	}
	return nil
}
