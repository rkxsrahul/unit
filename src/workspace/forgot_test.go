package workspace

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

func TestForgot(t *testing.T) {
	span := opentracing.StartSpan("Forgot")
	err := Forgot(email, span)
	if err != nil {
		t.Error("test case fail")
	}
	err = Forgot("xenon@test.com", span)
	if err == nil {
		t.Error("test case fail")
	}
	// err = Forgot("te@testing.com", span)
	// if err == nil {
	// 	t.Error("test case fail")
	// }
	// err = Forgot("t@testing.com", span)
	// if err == nil {
	// 	t.Error("test case fail")
	// }

}

func TestRecoverWorkspace(t *testing.T) {
	span := opentracing.StartSpan("Recover Workspace")
	status, _ := RecoverWorkspace(token, span)
	if status != 200 {
		t.Error("test case fail")
	}
	status, _ = RecoverWorkspace("333333", span)
	if status != 404 {
		t.Error("test case fail")
	}
}
