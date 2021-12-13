package health

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// ServiceHealth is a method to check service each components health
func ServiceHealth(parentspan opentracing.Span) error {
	span := opentracing.StartSpan("service health function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()

	// checking cockroach database and redis database
	span.LogKV("check", "cockroach and redis database")
	err := Healthz()
	if err != nil {
		log.Println(err)
		span.LogKV("return ", err.Error())
		return err
	}

	return nil
}
func Healthz() error {
	//checking health of cockroach database
	// connecting to db
	db, err := gorm.Open("postgres", config.DBConfig())
	if err != nil {
		log.Println(err)
		return errors.New("database connection not established")
	}
	// close db instance whenever whole work completed
	defer db.Close()
	//run db in debug mode
	db = db.Debug()

	//run sample query
	type Count struct {
		Count int64
	}
	var count Count
	db.Raw("select 1+1 as count").Scan(&count)
	log.Println("count...", count)
	if count.Count != 2 {
		return errors.New("Cockroach db is not working")
	}

	// checking health of redis database
	// convert redis db string to int
	redisDB, err := strconv.Atoi(config.Conf.Redis.Database)
	if err != nil {
		log.Println(err)
		return errors.New("Please pass valid redis database")
	}
	//create new redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port),
		Password: config.Conf.Redis.Pass,
		DB:       redisDB,
	})
	// check connection with server
	pong, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(pong)

	// check deployment service is working
	// http request for deploying WorkSpace
	req, err := http.NewRequest("GET", config.Conf.Address.Deployment+"/healthz", nil)
	if err != nil {
		log.Println("http....", err)
		return errors.New("Unable to connect to deployment service")
	}
	resp, err := http.DefaultClient.Do(req)
	log.Println(err)
	log.Println(resp)
	if err == nil && resp.StatusCode == 200 {
		return nil
	}
	return errors.New("Unable to connect to deployment service")

}
