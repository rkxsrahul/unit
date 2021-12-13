package scheduler

import (
	"strconv"
	"time"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"github.com/robfig/cron"
)

// Start is a function to start cronjobs
func Start() {

	DeleteUsers()
	c := cron.New()
	c.AddFunc("0 20 * * *", DeleteUsers)
	c.Start()
}

func DeleteUsers() {
	db := config.DB
	db.Exec("delete from accounts where verify_status='not_verified' AND creation_date<" + strconv.FormatInt(time.Now().Unix()-86400, 10) + ";")
	db.Exec("delete from workspace_members where member_email not in (select email from accounts);")
}
