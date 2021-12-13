package login

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	email         string = "test@testing.com"
	password      string = "testing"
	accountstatus string = "active"
	contactno     string = "8825383177"
	creationdate  int64  = 3600
	name          string = "xenon"
	roleid        string = "admin"
	userid        string = "1"
	varifystatus  string = "verified"
	tok           string = "882538"
)

// check token login
func TestTokenLogin(t *testing.T) {
	span := opentracing.StartSpan("simple token")
	//check token login with invalid token
	code, _ := TokenLogin(email, tok, "1", span)

	if code != 200 {
		t.Error("test case fail")
	}
	code, _ = TokenLogin("xenon@test.com", tok, "1", span)
	if code != 400 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin(email, "111111", "1", span)
	if code != 404 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin(email, "222222", "1", span)
	if code != 404 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin("testing@testing.com", tok, "1", span)
	if code != 404 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin("testing@testing.com", "444444", "2", span)
	if code != 404 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin(email, "444444", "2", span)
	if code != 404 {
		t.Error("test case fail")
	}

	code, _ = TokenLogin("t@testing.com", "444445", "3", span)
	if code != 404 {
		t.Error("test case fail")
	}

}
