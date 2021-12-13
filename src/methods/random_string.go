package methods

import (
	"log"
	"math/rand"
	"regexp"
	"time"
)

// RandomString is a method to generate random string
// on basis of length passed in parameter
func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt())
	}
	return string(bytes)
}

// generating random integer
func randInt() int {
	rdInt := rand.Intn(51)
	if rdInt <= 25 {
		rdInt = rdInt + 65
	} else {
		rdInt = rdInt + 71
	}
	return rdInt
}

//==============================================

// RandomStringIntegerOnly is a method to generate random string
// contains integer only on basis of length passed in parameter
func RandomStringIntegerOnly(l int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(rand.Intn(9) + 48)
	}
	return string(bytes)
}

//==============================================

// SlugifyEmail is method to remove special characters and only 5 character string returned
func SlugifyEmail(email string) string {
	var finalString string
	// remove all special character in email string
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	finalString = reg.ReplaceAllString(email, "")
	// return first five characters only
	return finalString[0:5]
}
