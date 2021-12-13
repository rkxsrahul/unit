package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	ot "github.com/opentracing/opentracing-go"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/ldap"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/routes"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/scheduler"
	tracing "git.xenonstack.com/stacklabs/stacklabs-auth/src/tracing"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// setup for reading flags for deciding whether to do configuration with toml or env variables
	conf := flag.String("conf", "environment", "set configuration from toml file or environment variables")

	file := flag.String("file", "", "set path of toml file")

	flag.Parse()

	if *conf == "environment" {
		log.Println("environment")
		config.ConfigurationWithEnv()
	} else if *conf == "toml" {
		log.Println("toml")
		if *file == "" {
			log.Println("Please pass toml file path")
			os.Exit(1)
		} else {
			err := config.ConfigurationWithToml(*file)
			if err != nil {
				log.Println("Error in parsing toml file")
				log.Println(err)
				os.Exit(1)
			}
		}
	} else {
		log.Println("Please pass valid arguments, conf can be set as toml or environment")
		os.Exit(1)
	}

	// checking environment
	if config.Conf.Service.Environment != "production" {

		// removing info file if any.
		_ = os.Remove("info.txt")

		// creating and opening info.txt file for writting logs
		file, err := os.OpenFile("info.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		// changing default writer of gin to file and std output
		gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

		// setting output for logs this will writes all logs to file
		log.SetOutput(gin.DefaultWriter)
		// writing log to check all in working
		log.Print("Logging to a file in Go!")
	}

	// initialize tracer
	trace, closer := tracing.InitJaeger("stacklabs-auth")
	defer closer.Close()
	ot.SetGlobalTracer(trace)

	database.CreateDatabase()

	//	database config
	dbConfig := config.DBConfig()
	//number of ideal connections
	var ideal int
	idealStr := config.Conf.Database.Ideal
	if idealStr == "" {
		ideal = 50
	} else {
		ideal, _ = strconv.Atoi(idealStr)
	}
	// connecting db using connection string
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		log.Println(err)
		return
	}
	// close db instance whenever whole work completed
	defer db.Close()
	db.DB().SetMaxIdleConns(ideal)
	db.DB().SetConnMaxLifetime(1 * time.Hour)
	config.DB = db

	err = ldap.CreateLDAPGroup(config.GroupName, config.GroupID)
	if err != nil {
		log.Println(err)
		if !strings.Contains(err.Error(), "Already") {
			log.Println(err)
			return
		}
	}

	//create auth-team database tables
	go database.CreateDatabaseTables()

	// initialize gin router
	router := gin.Default()

	//allowing CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowMethods("DELETE")
	router.Use(cors.New(corsConfig))

	// index handler to view all registered routes
	router.GET("/", func(c *gin.Context) {
		type finalPath struct {
			Method string
			Path   string
		}

		data := router.Routes()
		finalPaths := make([]finalPath, 0)

		for i := 0; i < len(data); i++ {
			finalPaths = append(finalPaths, finalPath{
				Path:   data[i].Path,
				Method: data[i].Method,
			})
		}
		c.JSON(200, gin.H{
			"routes": finalPaths,
		})
	})

	// service specific routes
	routes.V1Routes(router)

	//start scheduler
	scheduler.Start()
	// run the service
	router.Run(":" + config.Conf.Service.Port)
}
