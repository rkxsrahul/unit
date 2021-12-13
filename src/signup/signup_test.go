package signup

import (
	"testing"

	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"github.com/opentracing/opentracing-go"
)

func TestSignup(t *testing.T) {
	span := opentracing.StartSpan("send code")
	acc := database.Accounts{
		Userid:        "3",
		Password:      "ysrAk.FOIUClxmYbzwLg6WlDVLgb3i46o=",
		Email:         "xenon@stack.com",
		Name:          "Xenon",
		VerifyStatus:  varifystatus,
		RoleID:        roleid,
		AccountStatus: accountstatus,
	}

	_, status := Signup(acc, password, span)
	if status != false {
		t.Error("test case fail")
	}
}

func TestWithEmail(t *testing.T) {
	_, err := WithEmail("xenon@stack.com")
	if err != nil {
		t.Error("test case fail")
	}
	_, err = WithEmail(email)
	if err != nil {
		t.Error("test case fail")
	}
}
