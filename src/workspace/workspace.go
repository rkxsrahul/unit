package workspace

import (
	"log"

	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

//=========================================================================//

// CheckWorkspaceAvailability is a method to check workspace is there in database or not
func CheckWorkspaceAvailability(data database.Workspaces, parentspan opentracing.Span) (int, string) {
	log.Println(data)

	// start span from parent span context
	span := opentracing.StartSpan("workspace login method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// check workspace url is valid
	span.LogKV("task", "check workspace url is valid")
	if !IsWorkspaceURLValid(data.WorkspaceID) {
		span.LogKV("task", "workspace url is invalid")
		return 400, "URL should not contain special characters."
	}

	// connect to database
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return 501, datConst
	// }
	// // close database client after every operations
	// defer db.Close()
	db := config.DB

	// checking workspace already not present
	span.LogKV("task", "check workspace exist in db")
	var count int64
	db.Model(&database.Workspaces{}).Where("workspace_id=?", data.WorkspaceID).Count(&count)
	log.Println(count)
	if count != 0 {
		span.LogKV("task", "Workspace Found")
		return 409, "Workspace already exists."
	}
	span.LogKV("task", "Workspace Not Found")
	return 200, "Available."
}

func IsWorkspaceURLValid(name string) bool {
	for i := 0; i < len(name); i++ {
		if !((name[i] > 47 && name[i] < 58) || (name[i] > 64 && name[i] < 91) || (name[i] > 96 && name[i] < 123) || name[i] == 45 || name[i] == 95) {
			return false
		}
	}
	return true
}

//=========================================================================//

// Login is a method to login in workspace
func Login(wsURL string, parentspan opentracing.Span) (int, map[string]interface{}) {
	//result map
	mapd := make(map[string]interface{})
	// start span from parent span context
	span := opentracing.StartSpan("workspace login method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// connect to database
	span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	mapd["error"] = true
	// 	mapd["message"] = datConst
	// 	return 501, mapd
	// }
	// // close database client after every operations

	db := config.DB

	// details of workspace
	span.LogKV("task", "check workspace exist in db")
	ws := []database.Workspaces{}
	db.Where("workspace_id= ?", wsURL).Find(&ws)
	if len(ws) == 0 {
		span.LogKV("task", "Workspace Not Found")
		mapd["error"] = true
		mapd["message"] = "Workspace Not Found"
		return 404, mapd
	}
	span.LogKV("task", "Workspace Found")
	mapd["error"] = false
	mapd["worksapce"] = ws[0]
	return 200, mapd
}
