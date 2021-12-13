package mail

import (
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

// ReadToml is a method to read mail.toml file DialAndSend
// fetch template file path, subject of mail and image paths to be send in mail
func ReadToml(task string) (string, string, []string) {
	var fileData map[string]interface{}

	// final array of images string
	images := make([]string, 0)
	//read toml file
	_, err := toml.DecodeFile("./mail.toml", &fileData)
	if err != nil {
		log.Println(err)
		return "", "", images
	}

	// fetching data for verification mail
	value, ok := fileData[task]
	if !ok {
		log.Println(errors.New("there is no data in toml file regarding verification code"))
		return "", "", images
	}

	//type casting data in map of string  key and interface value
	data := value.(map[string]interface{})

	// fetch template file path from verify data
	tmplPath, ok := data["template"]
	if !ok {
		log.Println(errors.New("there is no template file path in toml file for verification mail"))
		return "", "", images
	}

	//fetch subject from verify data
	subject, ok := data["subject"]
	if !ok {
		log.Println(errors.New("there is no subject in toml file for verification mail"))
		return "", "", images
	}

	// check images are there
	imgInterface, ok := data["images"]
	if ok {
		// type casting array of interfaces
		imgArrayInterface := imgInterface.([]interface{})

		for i := 0; i < len(imgArrayInterface); i++ {
			images = append(images, imgArrayInterface[i].(string))
		}
		//finally return data
		return tmplPath.(string), subject.(string), images
	}
	return "", "", images
}
