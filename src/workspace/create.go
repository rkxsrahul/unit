package workspace

import (
	"log"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/aws"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	coreSchema "git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

const datConst string = "Unable to connect to database"

// CreateWorkspace is a method to create workspace
// send request to deployment service to deploy workspace
// save in database
func CreateWorkspace(email, token string, data database.Workspaces, parentspan opentracing.Span) (int, map[string]interface{}) {
	//result map
	mapd := make(map[string]interface{})
	// start span from parent span context
	span := opentracing.StartSpan("create workspace method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// checking workspace url is correct
	span.LogKV("task", "check workspace url is valid")
	wsURLParts := strings.Split(data.WorkspaceID, ".")
	if len(wsURLParts) != 0 && wsURLParts[0] == "" {
		span.LogKV("task", "workspace url is not correct")
		mapd["error"] = true
		mapd["message"] = "WorkSpace URL is required."

		return 400, mapd
	}
	if !IsWorkspaceURLValid(wsURLParts[0]) {
		span.LogKV("task", "workspace url is in-valid")
		mapd["error"] = true
		mapd["message"] = "URL should not contain special characters."
		return 400, mapd
	}

	// connect to database
	span.LogKV("task", "intialise db connection")
	db := config.DB

	// checking workspace already not present
	var count int64
	span.LogKV("task", "check workspace already exist in db")
	db.Model(&database.Workspaces{}).Where("workspace_id=?", wsURLParts[0]).Count(&count)
	if count != 0 {
		span.LogKV("task", "workspace already exist in db")
		mapd["error"] = true
		mapd["message"] = "Workspace already exists."

		return 409, mapd
	}

	loginac, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "user has no account")
		mapd["error"] = true
		mapd["message"] = "Please create account first"
		return 501, mapd
	}

	//if aws configurable
	if config.Conf.Service.ISAWS == "true" {
		//===============aws work========================================//
		//create enterprise bucket
		berr := aws.CreateBucket(wsURLParts[0] + "-sl")
		if berr == nil || strings.Contains(berr.Error(), "status code: 409") {
			//create policies related to bucket
			perr := aws.CreatePolicy(wsURLParts[0], "Admin")
			log.Println(perr)
			perr = aws.CreatePolicy(wsURLParts[0], "Enterprise")
			log.Println(perr)
		} else {
			span.LogKV("task", "unable to create bucket")
			mapd["error"] = true
			mapd["message"] = "Unable to create AWS team bucket."
			return 500, mapd
		}

		//create iam user
		_, err = aws.CreateIamUser(loginac.Userid)
		if err != nil && !strings.Contains(err.Error(), "status code: 409") {
			log.Println(err)
			span.LogKV("task", "unable to create IAM account")
			mapd["error"] = true
			mapd["message"] = "Unable to create IAM user account."
			return 500, mapd
		}
		perr := aws.AttachPolicy("sl"+loginac.Userid, "Admin", wsURLParts[0])
		if perr != nil && !strings.Contains(perr.Error(), "status code: 409") {
			log.Println(perr)
			span.LogKV("task", "unable to attach policy")
			mapd["error"] = true
			mapd["message"] = "Unable to attach policy."
			return 500, mapd
		}
		//create object in bucket
		berr = aws.CreateBucket(wsURLParts[0] + "-sl/" + "private-dataset-sl" + loginac.Userid + "/")
		if berr != nil && !strings.Contains(berr.Error(), "status code: 409") {
			log.Println(berr)
			span.LogKV("task", "unable to create bucket")
			mapd["error"] = true
			mapd["message"] = "Unable to create user folder in bucket."
			return 500, mapd
		}
	}

	// create workspace
	// span.LogKV("task", "deploy workspace")
	// isOk := createDedicatedWorkspace(data.WorkspaceID, token)
	// if !isOk {
	// 	span.LogKV("task", "unable to deploy workspace")
	// 	mapd["error"] = true
	// 	mapd["message"] = "Unable to create workspace."
	// 	return 501, mapd
	// }

	// if team name is nil
	if data.TeamName == "" {
		data.TeamName = wsURLParts[0]
	}

	//save data in database
	span.LogKV("task", "save workspace in database")

	// creating workspace
	db.Create(&database.Workspaces{
		WorkspaceID: wsURLParts[0],
		Status:      "new",
		TeamName:    data.TeamName,
		TeamSize:    data.TeamSize,
		TeamType:    data.TeamType,
		Created:     time.Now().Unix()})
	// add owner
	db.Create(&database.WorkspaceMembers{
		WorkspaceID: wsURLParts[0],
		MemberEmail: email,
		Role:        "owner",
		Joined:      time.Now().Unix(),
	})

	go sendNotificationMails(loginac, wsURLParts[0])

	// create jwt token
	span.LogKV("task", "generate jwt token")

	mapd = jwtToken.JwtTokenusingWID(loginac, wsURLParts[0], "owner")
	mapd["name"] = loginac.Name
	mapd["role_id"] = loginac.RoleID
	mapd["email"] = loginac.Email
	mapd["workspace_role"] = "owner"
	mapd["workspace_name"] = data.TeamName
	mapd["message"] = "Workspace created."

	span.LogKV("task", "send final output workspace created")
	return 200, mapd
}

func sendNotificationMails(acc coreSchema.Accounts, wsID string) {
	mapd := map[string]interface{}{
		"Workspace": wsID,
		"Name":      acc.Name,
		"Email":     acc.Email,
		"Phone":     acc.ContactNo,
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := mail.ReadToml("notification")

	// parse email template
	tmpl := mail.EmailTemplate(tmplPath, mapd)
	mails := strings.Split(config.Conf.Service.Mails, ",")
	//now sending mail
	for i := 0; i < len(mails); i++ {
		mail.SendMail(mails[i], subject, tmpl, images)
	}
}

// WorkspaceJWT is structure to be send in body for deploying a workspace
type WorkspaceJWT struct {
	WorkspaceURL string `json:"workspace_url"`
}

// createDedicatedWorkspace is a method to send request to deployment service to deploy a workspace
// func createDedicatedWorkspace(url, token string) bool {
// 	// setup body
// 	ws := WorkspaceJWT{WorkspaceURL: url}

// 	// marshal above data in JSON
// 	jsonBytes, err := json.Marshal(ws)
// 	if err != nil {
// 		log.Println("json.....", err)
// 		return false
// 	}

// 	// http request for deploying WorkSpace
// 	req, err := http.NewRequest("POST", config.Conf.Address.Deployment+"/sr/team_ingress", bytes.NewBuffer(jsonBytes))
// 	if err != nil {
// 		log.Println("http....", err)
// 		return false
// 	}
// 	req.Header.Set("Authorization", token)
// 	log.Println(req)
// 	// send request
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Println("workspace creation response error", err)
// 		return false
// 	}
// 	// check status
// 	if err == nil && resp.StatusCode == 201 {
// 		return true
// 	}
// 	log.Println(resp.StatusCode, "=-==--=", err)
// 	return false
// }
