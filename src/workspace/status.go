package workspace

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/jwtToken"
)

// Status is a method to check workspace deployed or not
func Status(wsURL string, parentspan opentracing.Span) (int, map[string]interface{}) {
	//result map
	mapd := make(map[string]interface{})
	// start span from parent span contex
	span := opentracing.StartSpan("login with token function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	span.LogKV("task", "check each url by making htttp request")
	// urls for checking workspace is deployed or not
	urls := []string{
		"https://" + wsURL,
		"https://" + wsURL + "/api/deployments/healthz",
	}
	log.Println(urls)

	// insecure transport enable
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// intialize http client
	client := &http.Client{Transport: tr}

	// checking each url status by sending http request
	for i := 0; i < len(urls); i++ {
		log.Println("Checking for URL: ", urls[i])
		req, err := http.NewRequest("GET", urls[i], nil)
		if err != nil {
			log.Println("Error while creating http request for "+urls[i]+": ", err)
			mapd["error"] = true
			mapd["message"] = "Please wait !! We are creating " + wsURL + " workspace."
			return 202, mapd
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error while sending request for "+urls[i]+": ", err)
			mapd["error"] = true
			mapd["message"] = "Please wait !! We are creating " + wsURL + " workspace."
			return 202, mapd
		}

		log.Println("Status for " + urls[i] + ": " + resp.Status)

		if resp.StatusCode != 200 {
			mapd["error"] = true
			mapd["message"] = "Please wait !! We are creating " + wsURL + " workspace."
			return 202, mapd
		}
	}

	// connect to database
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	mapd["error"] = true
	// 	mapd["message"] = "Unable to connect to database"
	// 	return 501, mapd
	// }
	// // close database client after every operations
	db := config.DB
	defer db.Close()

	// split workspace url in parts
	wsURLParts := strings.Split(wsURL, ".")

	// fetch member detail of workpsace
	span.LogKV("task", "fetch member details")
	var mem []database.WorkspaceMembers
	db.Where("workspace_id = ? AND role=?", wsURLParts[0], "owner").Find(&mem)
	log.Println(mem)
	if len(mem) == 0 {
		span.LogKV("task", "send final output when workspace status is hit before workspace create")
		mapd["error"] = true
		mapd["message"] = "Please create workspace first"
		return 501, mapd
	}
	// fetch member details
	span.LogKV("task", "fetch account details")
	loginac, err := accounts.GetAccountForEmail(mem[0].MemberEmail)
	if err != nil {
		span.LogKV("task", "user has no account")
		mapd["error"] = true
		mapd["message"] = "Please create account first"
		return 501, mapd
	}
	log.Println(loginac)

	// fetch workspace details
	span.LogKV("task", "fetch workspace details")
	var work []database.Workspaces
	db.Where("workspace_id = ?", wsURLParts[0]).Find(&work)
	log.Println(work)
	if len(work) == 0 {
		span.LogKV("task", "send final output when workspace status is hit before workspace create")
		mapd["error"] = true
		mapd["message"] = "Please create workspace first"
		return 501, mapd
	}

	// create jwt token
	span.LogKV("task", "generate jwt token")
	mapd = jwtToken.JwtTokenusingWID(loginac, wsURLParts[0], mem[0].Role)
	mapd["name"] = loginac.Name
	mapd["role_id"] = loginac.RoleID
	mapd["email"] = loginac.Email
	mapd["workspace_role"] = mem[0].Role
	mapd["workspace_name"] = work[0].TeamName
	mapd["message"] = "Your workspace is ready to use."

	span.LogKV("task", "send final output")
	//send final result
	return 200, mapd
}
