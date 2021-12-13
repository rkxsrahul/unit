package methods

import (
	"testing"
)

//Checking email is valid or not
func TestValidateEmail(t *testing.T) {
	//passing valid email address
	email := "testing@xenonstack.com"
	valid_email := ValidateEmail(email)
	if valid_email != true {
		t.Error("test case fail")
	}

	//passing envalid email address
	email = "testingxenonstackcom"
	invalid_email := ValidateEmail(email)
	if invalid_email == true {
		t.Error("test case fail")
	}

}
