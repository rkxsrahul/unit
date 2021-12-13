package methods

import "testing"

//test to check HMAC-SHA1 signer
func TestSign(t *testing.T) {

	//test sign with different key
	//first sign with same password different key --> key1
	key1 := RandomString(5)
	password := "rahul"
	first := Sign(key1, password)

	//second sign with same password different key --> key2
	key2 := RandomString(5)
	password = "rahul"
	second := Sign(key2, password)

	if first == second {
		t.Error("test case fail")
	}

}
