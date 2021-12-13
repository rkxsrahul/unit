package methods

import (
	"regexp"
	"strings"
)

// ValidateEmail is a method for validating email
func ValidateEmail(email string) bool {
	email = strings.ToLower(email)
	Re := regexp.MustCompile(`^([-\w\d]+)(\.[-\w\d]+)*@([-\w\d]+)(\.[-\w\d]{2,})$`)
	return Re.MatchString(email)
}
