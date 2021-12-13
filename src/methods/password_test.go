package methods

import "testing"

// password check
func TestCheckPassword(t *testing.T) {
	//valid pass check
	password := "RKrahul@321"
	valid_pass := CheckPassword(password)
	if valid_pass != true {
		t.Error("test case fail")
	}

	//invalid password check with no special characters
	password = "RKrahul321"
	invalid_pass := CheckPassword(password)
	if invalid_pass == true {
		t.Error("test case fail")
	}

	//invalid password check with no uppercase alphabets
	password = "rahul@321"
	invalid_pass = CheckPassword(password)
	if invalid_pass == true {
		t.Error("test case fail")
	}

	// invalid password check with no lowercase alphabets
	password = "RK@321"
	invalid_pass = CheckPassword(password)
	if invalid_pass == true {
		t.Error("test case fail")
	}

	// invalid password check with no numeric value
	password = "RKrahul@"
	invalid_pass = CheckPassword(password)
	if invalid_pass == true {
		t.Error("test case fail")
	}

}

//test to check hash password is equal to password
func TestCheckHashForPassword(t *testing.T) {
	//check hash password with valid password
	password := "xenonstack"
	hash_pass := HashForNewPassword(password)
	bol := CheckHashForPassword(hash_pass, password)
	if bol != true {
		t.Error("test case fail")
	}

	//check hash password with invalid password
	hash_pass = HashForNewPassword(password)
	bol = CheckHashForPassword(hash_pass, "xenon")
	if bol == true {
		t.Error("test case fail")
	}
}
