package ldap

import (
	"os"
	"testing"
)

func TestLdap_auth(t *testing.T) {

	if len(os.Getenv("LDAP_CN")) == 0 {
		t.Errorf("Environment variable LDAP_CN doesn't exist")
	}

	if len(os.Getenv("LDAP_PASS")) == 0 {
		t.Errorf("Environment variable LDAP_PASS doesn't exist")
	}

	if len(os.Getenv("LDAP_PORT")) == 0 {
		t.Errorf("Environment variable LDAP_PORT doesn't exist")
	}

	if len(os.Getenv("LDAP_HOST")) == 0 {
		t.Errorf("Environment variable LDAP_HOST doesn't exist")
	}

	if len(os.Getenv("LDAP_DC")) == 0 {
		t.Errorf("Environment variable LDAP_DC doesn't exist")
	}

	_, _, msg := Authentication("abc", "abc")
	if msg == "" {
		t.Errorf("LDAP_auth doesn't returns correct message")
	}

}
