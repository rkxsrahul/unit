package admin

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/accounts"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

// DeleteAccount is a method to delete a account from database
func DeleteAccount(email string, parentspan opentracing.Span) error {
	// start span from parent span context
	db := config.DB
	span := opentracing.StartSpan("forgot password method", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// connecting to db
	// span.LogKV("task", "intialise db connection")
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	span.LogKV("task", "send final output after error in connecting to db")
	// 	log.Println(err)
	// 	return errors.New("Unable to connect to database")
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()

	//fetch account
	span.LogKV("task", "fetch account details from db")
	acc, err := accounts.GetAccountForEmail(email)
	if err != nil {
		span.LogKV("task", "send final output after error in fetching account from db")
		log.Println(err)
		return err
	}

	//delete from workspace Member table
	span.LogKV("task", "delete from database")
	row := db.Where("member_email=?", acc.Email).Delete(&database.WorkspaceMembers{})
	log.Println("members...", row)

	//delete account from core auth
	err = accounts.DeleteAccount(email)
	span.LogKV("task", "send final output")
	return err
}
