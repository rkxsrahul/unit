package workspace

import (
	"errors"
	"log"
	"net/http"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/aws"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// DeleteWorkspace is a method to delete a workspace from database and kube
func DeleteWorkspace(wsID, token string, parentspan opentracing.Span) error {
	// start span from parent span context
	span := opentracing.StartSpan("delete workspace method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connect to database
	span.LogKV("task", "intialise db connection")
	db := config.DB

	//split workspace url to fetch workspace id
	wsURLs := strings.Split(wsID, ".")

	var count int64
	db.Model(&database.Workspaces{}).Where("workspace_id=?", wsURLs[0]).Count(&count)

	span.LogKV("task", "delete s3 and iam from aws")
	if config.Conf.Service.ISAWS == "true" {
		err := deleteAWSWorkspace(wsURLs[0])
		if err != nil {
			log.Println(err)
			return err
		}
	}
	span.LogKV("task", "delete workspace from cluster")
	// err := deleteDedicatedWorkspace(wsID, token)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	span.LogKV("task", "delete workspace from database")
	row := db.Where("workspace_id=?", wsURLs[0]).Delete(&database.WorkspaceMembers{}).RowsAffected
	row = db.Where("workspace_id=?", wsURLs[0]).Delete(&database.Workspaces{}).RowsAffected
	_ = row
	span.LogKV("task", "send final output")
	return nil
}

//=============================================================================//

func deleteDedicatedWorkspace(wsID, token string) error {
	// url to be hit for deleting workspace
	url := config.Conf.Address.Deployment + "/sr/team_ingress/" + wsID

	//http Request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("http request ....", err)
		return err
	}

	// set authorization
	req.Header.Set("Authorization", token)
	// send above http request to server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("http response error ....", err)
		return err
	}
	if resp.StatusCode != 200 {
		log.Println("http response ....", resp.StatusCode)
		return errors.New("workspace is not deleted succesfully")
	}
	return nil
}

// ====================================================================================================
func deleteAWSWorkspace(wsID string) error {
	// delete policies
	err := aws.DeletePolicies(wsID)
	if err != nil {
		return err
	}
	// delete bucket
	err = aws.DeleteBucket(wsID + "-sl")
	if err != nil && !strings.Contains(err.Error(), "status code: 404") {
		return err
	}
	return nil
}
