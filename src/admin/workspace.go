package admin

import (
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// ListWorkSpaces is a method to list all deployed workspaces
func ListWorkSpaces(parentspan opentracing.Span) ([]database.Workspaces, error) {
	// start span from parent span context
	span := opentracing.StartSpan("list workspaces method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connect to database
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", coreConfig.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return []database.Workspaces{}, errors.New("Unable to connect to database")
	// }
	// // close database client after every operations
	// defer db.Close()
	db := config.DB
	// fetch workspaces from database
	span.LogKV("task", "fetch workspaces from database")
	var workspaces []database.Workspaces
	db.Find(&workspaces)
	span.LogKV("task", "send final list")
	return workspaces, nil
}

//============================================================================//

// DeleteWorkspace is a method to delete a workspace from database and kube
// func DeleteWorkspace(wsID, token string, parentspan opentracing.Span) error {
// 	// start span from parent span context
// 	span := opentracing.StartSpan("delete workspace method", opentracing.ChildOf(parentspan.Context()))
// 	defer span.Finish()
// 	// connect to database
// 	span.LogKV("task", "intialise db connection")
// 	db, err := gorm.Open("postgres", coreConfig.DBConfig())
// 	if err != nil {
// 		span.LogKV("task", "send final output after error in connecting to db")
// 		log.Println(err)
// 		return errors.New("Unable to connect to database")
// 	}
// 	// close database client after every operations
// 	defer db.Close()

// 	//split workspace url to fetch workspace id
// 	wsURLs := strings.Split(wsID, ".")
// 	if len(wsURLs) < 2 {
// 		return errors.New("Please pass full link of your workspace including company FQDN")
// 	}
// 	span.LogKV("task", "delete workspace from cluster")
// 	err = deleteDedicatedWorkspace(wsID, token)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	span.LogKV("task", "delete workspace from database")
// 	row := db.Where("workspace_id=?", wsURLs[0]).Delete(&database.WorkspaceMembers{}).RowsAffected
// 	log.Println("members.....", row)
// 	row = db.Where("workspace_id=?", wsURLs[0]).Delete(&database.Workspaces{}).RowsAffected
// 	log.Println("workspace.....", row)

// 	span.LogKV("task", "send final output")
// 	return nil
// }

// //=============================================================================//

// func deleteDedicatedWorkspace(wsID, token string) error {
// 	// url to be hit for deleting workspace
// 	url := config.Conf.Address.Deployment + "/sr/nexa_workspace/" + wsID + "/delete"
// 	log.Println(url)

// 	//http Request
// 	req, err := http.NewRequest("POST", url, nil)
// 	if err != nil {
// 		log.Println("http request ....", err)
// 		return err
// 	}

// 	// set authorization
// 	req.Header.Set("Authorization", methods.Sign("kfT6fgWg0fN", wsID))
// 	// send above http request to server
// 	resp, err := http.DefaultClient.Do(req)
// 	log.Println("http response error ....", err)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("http response ....", resp.StatusCode)
// 	if resp.StatusCode != 200 {
// 		return errors.New("workspace is not deleted succesfully")
// 	}
// 	return nil
// }
