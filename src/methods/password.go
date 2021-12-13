package methods

import (
	"log"
	"strings"
)

// CheckPassword is a method for validating password
func CheckPassword(pass string) bool {
	// checking length
	if len(pass) < 8 {
		log.Println("length")
		return false
	}
	// checking password contains any special characters
	if !strings.ContainsAny(pass, "!@#$%^&*-?") {
		log.Println("special")
		return false
	}
	// checking password contains any uppercase aphabets
	if !strings.ContainsAny(pass, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		log.Println("upper")
		return false
	}
	// checking password contains any lowercase aphabets
	if !strings.ContainsAny(pass, "abcdefghijklmnopqrstuvwxyz") {
		log.Println("lower")
		return false
	}
	// checking password contains any numeric
	if !strings.ContainsAny(pass, "1234567890") {
		log.Println("number")
		return false
	}
	return true
}

// HashForNewPassword is a method for encoding password
func HashForNewPassword(pass string) string {

	// random string is generated as a key
	key := RandomString(5)
	// passhash generated
	passwordHash := key + "." + Sign(key, pass)
	// return hashed password
	return passwordHash
}

// CheckHashForPassword is a method matching two passwords
func CheckHashForPassword(passwordHash, password string) bool {

	//spliting hash password on basis of .
	passHashParts := strings.Split(passwordHash, ".")
	// if there are less then or more then two parts the hashed password is wrong
	if len(passHashParts) != 2 {
		return false
	}

	// calculate hashpassword using new password
	reCalculatedHash := passHashParts[0] + "." + Sign(passHashParts[0], password)
	// chech old hash with new hash
	return reCalculatedHash == passwordHash
}
