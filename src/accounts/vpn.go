package accounts

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

func VPNAccess(email string) database.OpenVPNInformation {
	db := config.DB

	data := database.OpenVPNInformation{}

	db.Table("open_vpn_informations").Where("email=?", email).Find(&data)

	return data
}

func Server(email string) string {
	db := config.DB

	data := database.OpenVPNInformation{}

	db.Table("open_vpn_informations").Where("email=?", email).Find(&data)

	if data.Email == email {
		return data.FileName
	}
	count := 1

	//fetch total count
	db.Table("open_vpn_informations").Count(&count)

	count = (count / 75)
	count++
	count++
	if count > 11 {
		rand.Seed(time.Now().UnixNano())
		count = rand.Intn(12)
	}
	if count < 2 {
		count = 3
	}
	//variable for store the list of files
	var files []string

	path := config.Conf.Service.Environment
	if path == "" {
		path = "dev"
	}
	//fetch the files list
	root := "./vpn-files/" + path + "/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Println(err)
		return ""
	}
	path = strings.TrimPrefix(files[count], "vpn-files/"+path+"/")

	file := strings.Split(path, ".")

	tomlFile := root + "/vpn.toml"

	mapData := make(map[string]interface{})
	// parse toml file and save data config structure
	_, err = toml.DecodeFile(tomlFile, &mapData)
	if err != nil {
		log.Println(err)
		return ""
	}

	//manage password
	mapdata := mapData[file[0]].((map[string]interface{}))
	//store infomation in database
	data = database.OpenVPNInformation{}
	data.Email = email
	data.FileName = files[count]
	data.Username = fmt.Sprint(mapdata["username"])
	data.Password = fmt.Sprint(mapdata["password"])

	db.Create(&data)

	return files[count]
}
