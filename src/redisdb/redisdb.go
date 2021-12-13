package redisdb

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// client is a redis common client
var client *redis.Client

// initialise is a method to initialise a redis client
func initialise() error {
	//again reset the config if any changes in toml file or environment variables
	config.SetConfig()

	// convert redis db string to int
	db, _ := strconv.Atoi(config.Conf.Redis.Database)
	// redis db client creations
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Conf.Redis.Host, config.Conf.Redis.Port),
		Password: config.Conf.Redis.Pass,
		DB:       db,
	})

	// check connection with server
	pong, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(pong)
	return nil
}

// CheckToken is a method to check token exists in redis database
func CheckToken(token string) error {
	// initialise redis client
	err := initialise()
	if err != nil {
		return err
	}
	// check token
	_, err = client.Get(token).Result()
	if err != nil {
		// when token not exist
		log.Println(err)
		return err
	}
	// log.Println(val)
	return nil
}

// SaveToken is a method saving token in redis
func SaveToken(key, value string, expire time.Duration) {
	// initialise redis client
	err := initialise()
	if err != nil {
		return
	}
	// save token with expiry
	err = client.Set(key, value, expire).Err()
	if err != nil {
		log.Println(err)
		return
	}
}

// DeleteToken is a method for deleting token from redis
func DeleteToken(key string) error {
	// initialise redis client
	err := initialise()
	if err != nil {
		return err
	}
	// delete token from redis
	val, err := client.Del(key).Result()
	if err != nil {
		// if any error
		log.Println(err)
		return err
	}
	log.Println(val)
	return nil
}
