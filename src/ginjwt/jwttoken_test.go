package ginjwt

import (
	"testing"
)

const (
	userid string = "1"
	name   string = "xenon"
	email  string = "test@testing.com"
	roleid string = "user"
)

// test to check GinJwtToken
func TestGinToken(t *testing.T) {
	claims := make(map[string]interface{})
	// populate claims map
	claims["id"] = userid
	claims["name"] = name
	claims["email"] = email
	claims["sys_role"] = roleid
	mapd, _ := GinJwtToken(claims)

	if mapd == nil {
		t.Error("test case fail")
	}
}
