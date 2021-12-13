package jwtToken

import (
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ginjwt"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/redisdb"
)

// JwtToken is a method for creating claims map to be added in a token
// and also save sessions in cockroach database and redis database
func JwtToken(acs database.Accounts, parentspan opentracing.Span) map[string]interface{} {
	// start span from parent span contex
	span := opentracing.StartSpan("jwt token function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// intialise claims map
	span.LogKV("task", "intialise claims map")
	claims := make(map[string]interface{})
	// populate claims map
	claims["id"] = acs.Userid
	claims["name"] = acs.Name
	claims["email"] = acs.Email
	claims["sys_role"] = acs.RoleID

	// generate jwt token, expiration time and extra info like (expire jwt time, start and end time)
	span.LogKV("task", "generate jwt token")
	mapd, info := ginjwt.GinJwtToken(claims)

	// check token is empty or not
	if mapd["token"].(string) == "" {
		span.LogKV("task", "send final output when token is empty")
		return mapd
	}
	span.LogKV("task", "save session")

	// remove all other sessions from session storage and save this session
	SaveSessions(acs.Userid, mapd["token"].(string), info)

	span.LogKV("task", "send final output")
	return mapd
}

// JwtRefreshToken is a method for save old claims in a token
// and also save sessions in cockroach database and redis database
func JwtRefreshToken(claims map[string]interface{}, parentspan opentracing.Span) map[string]interface{} {
	// start span from parent span contex
	span := opentracing.StartSpan("jwt refresh token function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	// intialise claims map

	// generate jwt token, expiration time and extra info like (expire jwt time, start and end time)
	span.LogKV("task", "generate jwt token")
	mapd, _ := ginjwt.GinJwtToken(claims)

	// check token is empty or not
	if mapd["token"].(string) == "" {
		span.LogKV("task", "send final output when token is empty")
		return mapd
	}
	span.LogKV("task", "save session")

	// remove all other sessions from session storage and save this session

	span.LogKV("task", "send final output")
	return mapd
}

// SaveSessions is a method for saving session details in redis and cockroachdb
func SaveSessions(userid, newSessToken string, info map[string]interface{}) {
	// save token in redis
	redisdb.SaveToken(newSessToken, userid, info["expire"].(time.Duration))

	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	// deleting other active sessions of that user
	if config.Conf.Service.IsLogoutOthers == "true" {
		// fetch active session from dbs
		var actses []database.ActiveSessions
		db.Where("userid=?", userid).Find(&actses)
		// delete active session from redis
		for i := 0; i < len(actses); i++ {
			if actses[i].End >= time.Now().Unix() {
				err := redisdb.DeleteToken(actses[i].SessionID)
				if err != nil {
					log.Println(err)
				}
				// log.Println(val)
			}
		}
		// delete all session from db
		db.Exec("delete from active_sessions where userid= '" + userid + "';")
	}

	// creating one active session
	db.Create(&database.ActiveSessions{
		Userid:    userid,
		SessionID: newSessToken,
		Start:     info["start"].(int64),
		End:       info["end"].(int64)})
}

// DeleteTokenFromDb is a method to delete saved jwt token from db
func DeleteTokenFromDb(token string) {

	// connecting to db
	// db, err := gorm.Open("postgres", config.DBConfig())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// // close db instance whenever whole work completed
	// defer db.Close()
	db := config.DB
	db.Exec("delete from active_sessions where session_id= '" + token + "';")
}

// JwtToken is a function to generate jwt token and pass workspace details in claims
func JwtTokenusingWID(acs database.Accounts, workspace, role string) map[string]interface{} {

	claims := make(map[string]interface{}, 0)

	claims["id"] = acs.Userid
	claims["name"] = acs.Name
	claims["email"] = acs.Email
	claims["sys_role"] = acs.RoleID

	if workspace != "" {
		claims["workspace"] = workspace
		claims["role"] = role
	}
	mapd, info := ginjwt.GinJwtToken(claims)
	if val, ok := mapd["token"].(string); !ok || val == "" {
		return mapd
	}

	// remove all other sessions from session storage
	SaveSessions(acs.Userid, mapd["token"].(string), info)
	return mapd
}
