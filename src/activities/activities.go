package activities

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// RecordActivity is a method use to record activity of a user in activity table
func RecordActivity(activity database.Activities) {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// recording users activities
	db.Create(&activity)
}

// GetLoginActivities is a method used to get login activities of a user
func GetLoginActivities(email string) ([]database.Activities, error) {
	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return []database.Activities{}, errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// fetching activities of user from activities table
	var activities []database.Activities
	db.Where("email=?", email).Order("timestamp desc").Limit(5).Find(&activities)
	return activities, nil
}
