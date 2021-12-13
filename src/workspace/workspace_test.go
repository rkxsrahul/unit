package workspace

import (
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	opentracing "github.com/opentracing/opentracing-go"
)

func TestCheckWorkspaceAvailability(t *testing.T) {

	span := opentracing.StartSpan("simple changepassword")

	data := database.Workspaces{
		WorkspaceID:   "xenonstac",
		UsagePolicies: "usage",
		Status:        "new",
		TeamName:      "team",
		TeamSize:      "11",
		TeamType:      "admin",
		Created:       1621511267,
	}
	CheckWorkspaceAvailability(data, span)
}

// func TestLogin(t *testing.T) {
// 	span := opentracing.StartSpan("simple changepassword")
// 	status, _ := Login("xenonstacks", span)
// 	log.Println("status", status)
// 	db := config.DB
// 	ws := []database.Workspaces{}
// 	db.Where("workspace_id= ?", "xenonstack").Find(&ws)
// 	log.Println("mmmmmmm", len(ws))

// }
