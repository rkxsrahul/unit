package accounts

import (
	"testing"
)

func TestVPNAccess(t *testing.T) {
	// test with valid email
	Vpn_data := VPNAccess(email)
	if Vpn_data.FileName != "tomal" {
		t.Error("test case fail")
	}

	//test with invalid email
	Vpn_data = VPNAccess("xenon@testing.com")
	if Vpn_data.FileName == "tomal" {
		t.Error("test case fail")
	}
}
func TestServer(t *testing.T) {
	file := Server(email)
	if file != "tomal" {
		t.Error("test case fail")
	}

}
